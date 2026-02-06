# geekedu-project 在线视频学习平台
## 项目介绍
基于Go微服务（Gin+gRPC）+ React的前后端分离付费视频学习平台，实现课程发布、购买、安全播放等核心能力，采用阿里云OSS私有Bucket保护资源，支持Docker Compose一键部署。

## 项目目录结构
```
geekedu-project/
├── system-design.md       # 系统设计文档（含架构图、接口文档、数据库设计、OSS配置）
├── README.md              # 项目部署&使用说明
├── web-server/            # Gin Web服务（提供RESTful API、JWT鉴权）
├── logic-server/          # gRPC核心服务（用户/课程/订单/播放业务逻辑）
├── frontend/              # React前端项目（B站风格交互界面）
├── deploy/                # 部署配置文件（Docker Compose、环境变量模板）
└── .env                   # 环境配置文件（含OSS Key等敏感信息，已加入.gitignore）
```


## 前置准备
### 1. 环境安装
需提前安装以下工具：
- Git（克隆项目）
- Docker & Docker Compose（容器化部署）
- 阿里云账号（开通OSS服务）


### 2. 配置阿里云OSS Key（必选）
**步骤1：开通阿里云OSS并获取配置信息**
1. 登录阿里云控制台，进入「对象存储OSS」服务；
2. 创建**私有Bucket**（读写权限选择「私有」，防止资源公开访问），记录Bucket名称；
3. 复制Bucket所属地域的**Endpoint**（如`oss-cn-beijing.aliyuncs.com`）；
4. 进入「AccessKey管理」页面，创建/获取**AccessKey ID**和**AccessKey Secret**（建议使用子账号并配置最小权限）；
5. 为Bucket配置跨域规则：允许前端/后端域名的`GET/POST`请求（跨域规则配置路径：OSS控制台 → 对应Bucket → 权限管理 → 跨域设置）。

**步骤2：填写项目配置文件**
在项目根目录`geekedu-project/`下创建`.env`文件，将以下内容复制后替换为你的OSS配置：
```env
# 阿里云OSS配置（必填）
OSS_ENDPOINT=你的OSS地域节点（如oss-cn-beijing.aliyuncs.com）
OSS_ACCESS_KEY_ID=你的AccessKey ID
OSS_ACCESS_KEY_SECRET=你的AccessKey Secret
OSS_BUCKET_NAME=你的私有Bucket名称
OSS_PRESIGNED_EXPIRE=3600  # 预签名URL有效期（单位：秒）

# MySQL配置（Docker启动无需修改）
MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=123456
MYSQL_DB=geekedu

# JWT&服务配置
JWT_SECRET=GeekEdu@2026#Secret
JWT_EXPIRE=3600
PORT=8080
GRPC_PORT=50051
CORS_ALLOW_ORIGIN=http://localhost
```


## 快速启动项目（对应演示视频「启动演示」）
采用Docker Compose一键启动所有服务（包含Web服务、gRPC服务、MySQL、前端）：
1. 克隆项目并进入目录：
   ```bash
   git clone [你的Git仓库地址]
   cd geekedu-project
   ```
2. 确保已完成「前置准备」中的OSS配置（`.env`文件已正确填写）；
3. 一键启动所有服务：
   ```bash
   docker-compose up -d
   ```
4. 验证服务启动状态：
   ```bash
   docker-compose ps
   ```
   若所有服务状态为`Up`，则启动成功。


## 核心功能演示（对应演示视频「功能演示」）
启动成功后，访问`http://localhost`（前端地址），演示流程如下：
1. **注册/登录**：
   - 点击「登录/注册」，选择「注册」填写用户名/密码（仅学生角色）；
   - 注册完成后切换到「登录」，使用账号密码登录。
2. **管理员发布课程**：
   - 使用管理员账号登录（需提前在数据库配置`role=1`）；
   - 点击「发布课程」，填写课程标题/价格/简介，上传封面（自动同步至OSS）后提交。
3. **上传课程视频**：
   - 管理员在「上传课程视频」模块选择已发布的课程，上传视频文件（自动压缩后同步至OSS）。
4. **学生购买课程**：
   - 学生账号登录后，在课程列表选择课程，点击「购买课程」完成订单创建。
5. **播放课程视频**：
   - 购买完成后点击「播放课程」，弹出视频列表弹窗，选择视频即可加载OSS预签名URL播放。


## 代码设计亮点（对应演示视频「亮点展示」）
1. **微服务拆分**：
   - 拆分`web-server`（API网关）和`logic-server`（gRPC核心服务），解耦请求转发与业务逻辑，便于横向扩展。
2. **OSS资源安全**：
   - 采用私有Bucket存储视频/封面，仅通过后端生成**时效预签名URL**提供访问，避免资源泄露。
3. **细粒度权限控制**：
   - 基于JWT+角色（学生/管理员）实现接口权限校验，管理员专属发布/上传接口，学生仅能访问已购课程资源。
4. **容器化部署**：
   - 通过Docker Compose统一管理所有服务，一键启动/停止，环境一致性强，降低部署成本。


## 提交物说明
1. **系统设计文档**：项目根目录下`system-design.md`，包含架构图、接口文档、数据库设计、OSS配置截图；
2. **完整源码**：包含`web-server`、`logic-server`、`frontend`、`deploy`等目录；
3. **演示视频**：项目提交视频中`project.mp4`，包含启动演示、功能演示、代码亮点讲解。