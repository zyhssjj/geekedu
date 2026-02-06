package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/panjf2000/ants/v2"
)

// 全局变量：远程上报地址（从.env读取）
var reportURL string

// 全局Goroutine池（控制并发数，可根据服务器调整）
var taskPool *ants.Pool

// FileProcessResult：单个文件的处理结果（匹配字段定义要求）
type FileProcessResult struct {
	FileName    string   // 文件名（对应details的key）
	ResultStr   string   // 逗号分隔的计算结果（对应details的value）
	Sum         int      // 该文件结果总和（对应sum_list的元素）
	ExprCount   int      // 该文件的算式数量（用于汇总total_count）
	Expressions []string // 新增：存储该文件的原始表达式列表
}

// CalculateResponse：计算接口的响应结构体
type CalculateResponse struct {
	ReqStart   int64             `json:"req_start"`
	ReqEnd     int64             `json:"req_end"`
	Details    map[string]string `json:"details"`
	SumList    string            `json:"sum_list"`
	TotalCount int               `json:"total_count"`
	ResultFile string            `json:"result_file"`
}

// ResultListResponse：结果列表接口的响应结构体
type ResultListResponse struct {
	Files []string `json:"files"`
}

// ReportPayload：远程上报的请求体结构体（严格匹配上报格式）
type ReportPayload struct {
	Username   string `json:"username"`
	UUID       string `json:"uuid"`
	ReqStart   int64  `json:"req_start"`
	ReqEnd     int64  `json:"req_end"`
	Details    string `json:"details"` // 注意：是JSON字符串（需转义），非JSON对象
	SumList    string `json:"sum_list"`
	TotalCount int    `json:"total_count"`
	MD5        string `json:"md5"`
}

// 初始化：加载配置、创建目录、初始化Goroutine池
func init() {
	// 1. 加载.env配置
	err := godotenv.Load()
	if err != nil {
		log.Fatal("加载.env配置失败：", err)
	}
	reportURL = os.Getenv("REPORT_URL")
	if reportURL == "" {
		log.Fatal("REPORT_URL未在.env中配置")
	}

	// 2. 自动创建results目录
	err = os.MkdirAll("./results", 0755)
	if err != nil {
		log.Fatal("创建results目录失败：", err)
	}

	// 3. 初始化Goroutine池（最大并发数10，可调整）
	taskPool, err = ants.NewPool(10)
	if err != nil {
		log.Fatal("初始化Goroutine池失败：", err)
	}
}

func main() {
	// 初始化Gin引擎
	r := gin.Default()

	// 注册接口
	r.POST("/upload", handleFileUpload)               // 原上传接口（可选保留）
	r.POST("/api/calculate", handleCalculate)         // 计算接口
	r.GET("/api/result/list", handleResultList)       // 结果列表接口
	r.GET("/api/result/detail", handleResultDetail)   // 结果详情/下载接口

	// 启动服务（端口8080）
	log.Println("服务启动：http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gin服务启动失败：", err)
	}
}

// handleFileUpload：原文件上传接口（可选保留）
func handleFileUpload(c *gin.Context) {
	// 1. 获取上传的多文件（前端表单name需为"calc_files"）
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get upload files: " + err.Error(),
		})
		return
	}
	files := form.File["calc_files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files uploaded",
		})
		return
	}

	// 2. 校验文件：仅允许.txt后缀
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if ext != ".txt" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "only txt files are supported: " + file.Filename,
			})
			return
		}
	}

	// 3. 并发处理所有文件
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		currentFile := file

		err := taskPool.Submit(func() {
			defer wg.Done()
			_, err := processSingleFile(currentFile)
			if err != nil {
				log.Printf("处理文件[%s]失败：%v", currentFile.Filename, err)
			} else {
				log.Printf("文件[%s]处理完成", currentFile.Filename)
			}
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to submit task: " + err.Error(),
			})
			return
		}
	}

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{
		"message": "all files processed successfully, results saved to ./results directory",
	})
}

// handleCalculate：计算接口（POST /api/calculate）
func handleCalculate(c *gin.Context) {
	// 1. 记录请求开始时间
	reqStart := time.Now().UnixMilli()

	// 2. 校验并获取必填参数（合并检查username/uuid）
	username := c.PostForm("username")
	uuid := c.PostForm("uuid")
	if username == "" || uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username and uuid are required",
		})
		return
	}

	// 3. 获取并校验上传文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get upload files: " + err.Error(),
		})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files uploaded",
		})
		return
	}
	for _, file := range files {
		if filepath.Ext(file.Filename) != ".txt" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "only txt files are supported: " + file.Filename,
			})
			return
		}
	}

	// 4. 并发处理文件，收集结果
	var (
		wg           sync.WaitGroup
		fileResults  []FileProcessResult
		resultLock   sync.Mutex
		processErr   error // 存储具体的处理错误信息
	)
	for _, file := range files {
		wg.Add(1)
		currentFile := file

		err := taskPool.Submit(func() {
			defer wg.Done()
			res, err := processSingleFile(currentFile)
			if err != nil {
				// 锁定并记录具体文件的解析错误
				resultLock.Lock()
				processErr = fmt.Errorf("invalid expression in file: %s", currentFile.Filename)
				resultLock.Unlock()
				return
			}
			// 无错误则记录结果
			resultLock.Lock()
			fileResults = append(fileResults, res)
			resultLock.Unlock()
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to submit task: " + err.Error(),
			})
			return
		}
	}
	wg.Wait()

	// 检查是否存在文件解析错误
	if processErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": processErr.Error(),
		})
		return
	}

	// 5. 汇总结果（匹配响应+上报要求）
	details := make(map[string]string, len(fileResults))
	for _, res := range fileResults {
		details[res.FileName] = res.ResultStr
	}

	// 按文件名排序生成sum_list
	sort.Slice(fileResults, func(i, j int) bool {
		return fileResults[i].FileName < fileResults[j].FileName
	})
	
	var sumListBuilder strings.Builder
	for i, res := range fileResults {
		if i > 0 {
			sumListBuilder.WriteString(",")
		}
		sumListBuilder.WriteString(strconv.Itoa(res.Sum))
	}
	sumList := sumListBuilder.String()

	// 计算total_count
	totalCount := 0
	for _, res := range fileResults {
		totalCount += res.ExprCount
	}

	// 6. 构造远程上报数据
	// 6.1 将details序列化为JSON字符串（上报要求的格式）
	detailsBytes, err := json.Marshal(details)
	var detailsStr string
	if err != nil {
		log.Printf("序列化details为JSON字符串失败：%v", err)
		detailsStr = "" // 失败时填充空字符串，不阻断上报
	} else {
		detailsStr = string(detailsBytes)
	}

	// 6.2 计算detailsStr的MD5哈希（32位小写）
	md5Hash := md5.New()
	md5Hash.Write([]byte(detailsStr))
	md5Str := hex.EncodeToString(md5Hash.Sum(nil))

	// 6.3 构造上报Payload
	reqEnd := time.Now().UnixMilli()
	reportPayload := ReportPayload{
		Username:   username,
		UUID:       uuid,
		ReqStart:   reqStart,
		ReqEnd:     reqEnd,
		Details:    detailsStr,
		SumList:    sumList,
		TotalCount: totalCount,
		MD5:        md5Str,
	}

	// 6.4 异步远程上报（不阻塞客户端响应）
	go func() {
		if err := sendReport(reportPayload); err != nil {
			log.Printf("远程上报失败：%v", err)
		} else {
			log.Printf("远程上报成功")
		}
	}()

	// 7. 生成result_file（响应给客户端）
	resultFileName := fmt.Sprintf("result_%d.txt", reqStart)
	resultFilePath := filepath.Join("./results", resultFileName)
	
	// 关键修改：构造“原始表达式=计算结果”的内容
	var resultContent strings.Builder
	// 按文件名升序排列（已在前面做过排序，直接遍历fileResults即可）
	for _, res := range fileResults {
		// 拆分该文件的结果字符串为结果列表
		results := strings.Split(res.ResultStr, ",")
		// 遍历原始表达式+对应结果，生成“表达式=结果”行
		for idx, expr := range res.Expressions {
			// 避免索引越界（理论上表达式和结果数量一致）
			if idx >= len(results) {
				break
			}
			// 拼接“原始表达式=计算结果”
			line := fmt.Sprintf("%s=%s\n", expr, results[idx])
			resultContent.WriteString(line)
		}
	}
	// 写入结果文件
	if err := os.WriteFile(resultFilePath, []byte(resultContent.String()), 0644); err != nil {
		log.Printf("写入结果文件[%s]失败：%v", resultFileName, err)
	}

	// 8. 返回响应给客户端
	c.JSON(http.StatusOK, CalculateResponse{
		ReqStart:   reqStart,
		ReqEnd:     reqEnd,
		Details:    details,
		SumList:    sumList,
		TotalCount: totalCount,
		ResultFile: resultFileName,
	})
}

// handleResultList：结果列表接口（GET /api/result/list）
func handleResultList(c *gin.Context) {
	// 读取results目录下的result_*.txt文件
	dir := "./results"
	entries, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read result directory: " + err.Error(),
		})
		return
	}

	// 过滤出result_开头的txt文件
	var resultFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "result_") && filepath.Ext(entry.Name()) == ".txt" {
			resultFiles = append(resultFiles, entry.Name())
		}
	}

	// 按文件名（时间戳）排序（从旧到新）
	sort.Strings(resultFiles)

	// 返回响应
	c.JSON(http.StatusOK, ResultListResponse{
		Files: resultFiles,
	})
}

// handleResultDetail：结果详情/下载接口（GET /api/result/detail）
func handleResultDetail(c *gin.Context) {
	// 1. 获取参数
	fileid := c.Query("fileid")
	if fileid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "fileid is required",
		})
		return
	}
	download := c.Query("download") == "1"

	// 2. 拼接文件路径
	filePath := filepath.Join("./results", fileid)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
		})
		return
	}

	// 3. 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read file: " + err.Error(),
		})
		return
	}

	// 4. 处理下载/查看
	if download {
		// 触发文件下载
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileid))
		c.Header("Content-Type", "application/octet-stream")
	} else {
		// 展示文本内容
		c.Header("Content-Type", "text/plain; charset=utf-8")
	}
	c.Writer.Write(content)
}

// processSingleFile：处理单个上传文件，返回该文件的计算结果
func processSingleFile(file *multipart.FileHeader) (FileProcessResult, error) {
	fileRes := FileProcessResult{
		FileName: file.Filename,
	}

	// 1. 打开上传文件
	srcFile, err := file.Open()
	if err != nil {
		return fileRes, fmt.Errorf("打开文件失败：%w", err)
	}
	defer srcFile.Close()

	// 2. 读取表达式（保留原始表达式，方便后续还原）
	var allExpressions []string // 存储所有非空行的表达式（包括错误的）
	scanner := bufio.NewScanner(srcFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			allExpressions = append(allExpressions, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fileRes, fmt.Errorf("读取文件内容失败：%w", err)
	}
	if len(allExpressions) == 0 {
		return fileRes, fmt.Errorf("文件无有效算式")
	}

	// 3. 计算表达式（跳过错误行，只处理合法表达式）
	var validExpressions []string // 存储合法的原始表达式
	var exprResults []int         // 存储合法表达式的计算结果
	for _, expr := range allExpressions {
		calcRes, err := calculateBinaryExpr(expr)
		if err != nil {
			// 记录错误日志，跳过当前错误行，继续处理下一行
			log.Printf("文件[%s]中表达式[%s]无效，已跳过：%v", file.Filename, expr, err)
			continue
		}
		// 仅将合法的表达式和结果加入列表
		validExpressions = append(validExpressions, expr)
		exprResults = append(exprResults, calcRes)
	}

	// 检查：如果所有表达式都无效，返回文件级错误
	if len(validExpressions) == 0 {
		return fileRes, fmt.Errorf("文件中无有效算式（所有行均为无效表达式）")
	}

	// 4. 填充有效结果到返回结构体
	fileRes.Expressions = validExpressions // 仅保存合法的表达式
	fileRes.ExprCount = len(validExpressions) // 仅统计合法表达式数量

	// 5. 生成结果字符串（仅包含合法表达式的计算结果）
	var resultStrBuilder strings.Builder
	for i, res := range exprResults {
		if i > 0 {
			resultStrBuilder.WriteString(",")
		}
		resultStrBuilder.WriteString(strconv.Itoa(res))
	}
	fileRes.ResultStr = resultStrBuilder.String()
// 新增日志：打印读取到的所有非空行
log.Printf("【文件读取完成】文件名：%s，读取到非空行数量：%d，内容：%v", file.Filename, len(allExpressions), allExpressions)
	
if err := scanner.Err(); err != nil {
	log.Printf("【文件读取错误】文件名：%s，错误：%v", file.Filename, err)
	return fileRes, fmt.Errorf("读取文件内容失败：%w", err)
}
if len(allExpressions) == 0 {
	log.Printf("【文件无内容】文件名：%s", file.Filename)
	return fileRes, fmt.Errorf("文件无有效算式")
}
	// 6. 计算合法表达式的总和
	sum := 0
	for _, res := range exprResults {
		sum += res
	}
	fileRes.Sum = sum

	// 7. 写入结果文件（原文件对应的.result文件）
	resultFilePath := fmt.Sprintf("./results/%s.result", file.Filename)
	if err := os.WriteFile(resultFilePath, []byte(fileRes.ResultStr), 0644); err != nil {
		return fileRes, fmt.Errorf("写入结果文件失败：%w", err)
	}

	return fileRes, nil
}

// calculateBinaryExpr：计算单行二元表达式
func calculateBinaryExpr(expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("表达式格式错误：%w", err)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("计算失败：%w", err)
	}
	resInt, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("结果非整数")
	}
	return int(resInt), nil
}

// sendReport：向远程地址上报数据（严格匹配上报格式）
func sendReport(payload ReportPayload) error {
	// 1. 序列化上报数据为JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("上报数据序列化失败：%w", err)
	}

	// 2. 构造POST请求
	req, err := http.NewRequest("POST", reportURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("构造上报请求失败：%w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// 3. 发送请求（设置超时，避免阻塞）
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送上报请求失败：%w", err)
	}
	defer resp.Body.Close()

	// 4. 检查响应状态（非2xx视为失败）
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("上报响应异常，状态码：%d", resp.StatusCode)
	}

	return nil
}