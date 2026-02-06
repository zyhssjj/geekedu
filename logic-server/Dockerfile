# 构建阶段
FROM golang:1.25.3-alpine AS builder

# 1. 创建上级目录（模拟本地项目根目录结构）
WORKDIR /app/parent

# 2. 复制 common 目录（到 /app/parent/common，对应 ../common）
COPY ./common ./common

# 3. 复制 proto 目录（到 /app/parent/proto，对应 ../proto，解决新报错）
COPY ./proto ./proto

# 4. 复制 logic-server 目录（到 /app/parent/logic-server，和本地结构一致）
COPY ./logic-server ./logic-server

# 5. 进入 logic-server 目录（容器内路径：/app/parent/logic-server）
WORKDIR /app/parent/logic-server
# 关键新增：配置国内 GOPROXY，解决第三方依赖下载失败（核心！）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# 6. 下载依赖（此时 ../common 和 ../proto 都存在，和本地一致，无报错）
RUN go mod download

# 7. 编译项目
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logic-server .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/parent/logic-server/logic-server .

# 暴露端口
EXPOSE 50051

# 启动命令
CMD ["./logic-server"]