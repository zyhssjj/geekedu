geekedu-project 在线视频学习平台
项目介绍
基于 Go 微服务（Gin+gRPC）+ React 的前后端分离付费视频学习平台，实现课程发布、购买、安全播放等核心能力，采用阿里云 OSS 私有 Bucket 保护资源，支持 Docker Compose 一键部署。
项目目录结构
plaintext
geekedu-project/
├── system-design.md       # 系统设计文档（含架构图、接口文档、数据库设计、OSS配置）
├── README.md              # 项目部署&使用说明
├── web-server/            # Gin Web服务（提供RESTful API、JWT鉴权）
├── logic-server/          # gRPC核心服务（用户/课程/订单/播放业务逻辑）
├── frontend/              # React前端项目（B站风格交互界面）
├── deploy/                # 部署配置文件（Docker Compose、环境变量模板）
└── .env                   # 环境配置文件（含OSS Key等敏感信息，已加入.gitignore）
前置准备
1. 环境安装
需提前安装以下工具：
Git（克隆项目）
Docker & Docker Compose（容器化部署）
阿里云账号（开通 OSS 服务）
2. 配置阿里云 OSS Key（必选）
步骤 1：开通阿里云 OSS 并获取配置信息
登录阿里云控制台，进入「对象存储 OSS」服务；
创建私有 Bucket（读写权限选择「私有」，防止资源公开访问），记录 Bucket 名称；
复制 Bucket 所属地域的Endpoint（如oss-cn-beijing.aliyuncs.com）；
进入「AccessKey 管理」页面，创建 / 获取AccessKey ID和AccessKey Secret（建议使用子账号并配置最小权限）；
为 Bucket 配置跨域规则：允许前端 / 后端域名的GET/POST请求（跨域规则配置路径：OSS 控制台 → 对应 Bucket → 权限管理 → 跨域设置）。
步骤 2：填写项目配置文件
在项目根目录geekedu-project/下创建.env文件，将以下内容复制后替换为你的 OSS 配置：
env
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
快速启动项目（对应演示视频「启动演示」）
采用 Docker Compose 一键启动所有服务（包含 Web 服务、gRPC 服务、MySQL、前端）：
克隆项目并进入目录：
bash
运行
git clone [你的Git仓库地址]
cd geekedu-project
确保已完成「前置准备」中的 OSS 配置（.env文件已正确填写）；
一键启动所有服务：
bash
运行
docker-compose up -d
验证服务启动状态：
bash
运行
docker-compose ps
若所有服务状态为Up，则启动成功。
核心功能演示（对应演示视频「功能演示」）
启动成功后，访问http://localhost（前端地址），演示流程如下：
注册 / 登录：
点击「登录 / 注册」，选择「注册」填写用户名 / 密码（仅学生角色）；
注册完成后切换到「登录」，使用账号密码登录。
管理员发布课程：
使用管理员账号登录（需提前在数据库配置role=1）；
点击「发布课程」，填写课程标题 / 价格 / 简介，上传封面（自动同步至 OSS）后提交。
上传课程视频：
管理员在「上传课程视频」模块选择已发布的课程，上传视频文件（自动压缩后同步至 OSS）。
学生购买课程：
学生账号登录后，在课程列表选择课程，点击「购买课程」完成订单创建。
播放课程视频：
购买完成后点击「播放课程」，弹出视频列表弹窗，选择视频即可加载 OSS 预签名 URL 播放。
代码设计亮点（对应演示视频「亮点展示」）
微服务拆分：
拆分web-server（API 网关）和logic-server（gRPC 核心服务），解耦请求转发与业务逻辑，便于横向扩展。
OSS 资源安全：
采用私有 Bucket 存储视频 / 封面，仅通过后端生成时效预签名 URL提供访问，避免资源泄露。
细粒度权限控制：
基于 JWT + 角色（学生 / 管理员）实现接口权限校验，管理员专属发布 / 上传接口，学生仅能访问已购课程资源。
容器化部署：
通过 Docker Compose 统一管理所有服务，一键启动 / 停止，环境一致性强，降低部署成本。
提交物说明
系统设计文档：项目根目录下system-design.md，包含架构图、接口文档、数据库设计、OSS 配置截图；
完整源码：包含web-server、logic-server、frontend、deploy等目录；
演示视频：项目根目录下project.mp4，包含启动演示、功能演示、代码亮点讲解。