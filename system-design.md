在线视频学习平台 - 系统设计文档
1. 项目概述
本项目是一款在线视频学习平台，支持管理员发布课程、上传视频，学生注册登录、购买课程、观看视频等核心功能。技术栈采用前后端分离架构，前端基于React+AntD构建，后端以Go为开发语言，通过gRPC实现服务通信，MySQL存储业务数据，阿里云OSS存储课程封面及视频文件，整体架构高可用、易扩展。

2. 系统架构设计

2.1 架构图

采用分层架构设计，自上而下分为前端层、网关层、服务层、数据存储层，各层职责清晰、解耦性强。


2.2 核心组件职责

- 前端层：负责用户交互界面渲染，通过Axios调用后端RESTful API，实现登录、课程浏览、视频播放等功能，基于React Player封装视频播放组件。

- Web服务层：基于Gin框架提供RESTful API，接收前端请求，通过JWT中间件做身份认证，再通过gRPC客户端调用核心服务，处理跨域、请求转发等逻辑。

- gRPC核心服务层：拆分用户、课程、订单、播放四大服务，实现业务逻辑封装，通过gRPC协议与Web服务通信，操作数据库和OSS完成数据读写。

- 数据存储层：MySQL存储结构化业务数据，OSS存储非结构化文件（封面、视频），通过预签名URL机制实现视频安全播放。

3. 数据库设计

3.1 数据库配置

# 数据库连接配置（Go后端）
db:
  driver: mysql
  dsn: root:123456@tcp(127.0.0.1:3306)/geekedu?charset=utf8mb4&parseTime=True&loc=Local
  max_open_conns: 20
  max_idle_conns: 10
  conn_max_lifetime: 3600s

3.2 核心表结构

3.2.1 用户表（user）

字段名

类型

是否主键

备注

id

int unsigned

是

用户ID，自增

username

varchar(50)

否

用户名，唯一

password

varchar(100)

否

BCrypt加密后的密码

role

tinyint

否

角色（0-学生，1-管理员）

created_at

datetime

否

创建时间

updated_at

datetime

否

更新时间

3.2.2 课程表（course）

字段名

类型

是否主键

备注

id

int unsigned

是

课程ID，自增

title

varchar(100)

否

课程标题

price

decimal(10,2)

否

课程价格

intro

text

否

课程简介

cover_oss_key

varchar(255)

否

课程封面在OSS中的存储键

create_user_id

int unsigned

否

创建者ID（关联user表id）

created_at

datetime

否

创建时间

updated_at

datetime

否

更新时间

3.2.3 视频表（video）

字段名

类型

是否主键

备注

id

int unsigned

是

视频ID，自增

course_id

int unsigned

否

所属课程ID（关联course表id）

title

varchar(100)

否

视频标题

oss_key

varchar(255)

否

视频在OSS中的存储键

created_at

datetime

否

创建时间

updated_at

datetime

否

更新时间

3.2.4 订单表（order）

字段名

类型

是否主键

备注

id

int unsigned

是

订单ID，自增

user_id

int unsigned

否

用户ID（关联user表id）

course_id

int unsigned

否

课程ID（关联course表id）

order_status

tinyint

否

订单状态（1-已支付，0-未支付）

created_at

datetime

否

创建时间

updated_at

datetime

否

更新时间

4. OSS配置（阿里云）

4.1 配置信息

# OSS配置（Go后端，封装在oss工具类中）
type OSSConfig struct {
  Endpoint        string // OSS地域节点，如：oss-cn-beijing.aliyuncs.com
  AccessKeyID     string // 阿里云AccessKeyID
  AccessKeySecret string // 阿里云AccessKeySecret
  BucketName      string // OSS桶名
  PresignedURLExpire int64 // 预签名URL有效期（秒），默认3600
}

// 全局配置实例
var ossConfig = OSSConfig{
  Endpoint:        "oss-cn-beijing.aliyuncs.com",
  AccessKeyID:     "your-access-key-id",
  AccessKeySecret: "your-access-key-secret",
  BucketName:      "geekedu-video-platform",
  PresignedURLExpire: 3600,
}

4.2 核心功能封装

- 文件上传：
       

  - 封面上传：采用简单上传模式，上传图片文件至OSS，返回存储键（cover_oss_key）。

  - 视频上传：采用分片上传模式，支持大文件（≤1GB），上传完成后返回存储键（oss_key）。

- 预签名URL生成：通过`GeneratePresignedURL(ossKey string)`方法生成临时播放URL，仅允许已购买课程的用户获取，有效期1小时，避免文件泄露。

- 权限控制：OSS桶设置为私有读写，仅通过后端生成的预签名URL可访问文件，前端无法直接操作OSS。

5. 接口文档

5.1 接口概述

接口分为RESTful API（前端调用）和gRPC接口（Web服务与核心服务通信），所有REST接口前缀为`/api/v1`，统一返回JSON格式数据，响应码遵循HTTP标准。

统一响应格式：

{
  "code": 200, // 状态码（200-成功，400-参数错误，500-服务异常）
  "msg": "操作成功", // 提示信息
  "data": {} // 业务数据
}

5.2 RESTful API接口

5.2.1 用户相关接口

接口路径

请求方法

请求参数

响应数据

说明

/auth/register

POST

{"username":"xxx","password":"xxx"}

{"code":200,"msg":"注册成功","data":{}}

学生注册

/auth/login

POST

{"username":"xxx","password":"xxx"}

{"code":200,"msg":"登录成功","data":{"uid":1,"token":"xxx","role":0}}

用户登录，返回JWT令牌

5.2.2 课程相关接口

接口路径

请求方法

请求参数

响应数据

说明

/courses

GET

无

{"code":200,"msg":"成功","data":{"courses":[{"id":1,"title":"xxx","price":99.00,"coverSignedUrl":"xxx",...}]}}

获取所有课程列表

/courses

POST

{"title":"xxx","price":99.00,"intro":"xxx","cover_oss_key":"xxx"}

{"code":200,"msg":"发布成功","data":{"course_id":1}}

管理员发布课程（需JWT认证，角色为1）

/courses/upload/cover

POST

FormData（key=cover，value=图片文件）

{"code":200,"msg":"上传成功","data":{"cover_oss_key":"xxx"}}

上传课程封面（管理员）

/courses/:course_id/videos

GET

路径参数course_id

{"code":200,"msg":"成功","data":{"videos":[{"video_id":1,"title":"xxx",...}]}}

获取课程下所有视频（需购买课程或管理员）

5.2.3 视频相关接口

接口路径

请求方法

请求参数

响应数据

说明

/courses/upload/video

POST

FormData（course_id=1，title=xxx，video=视频文件）

{"code":200,"msg":"上传成功","data":{"video_id":1}}

上传课程视频（管理员）

/player/:video_id

GET

路径参数video_id

{"code":200,"msg":"成功","data":{"signed_url":"xxx"}}

获取视频预签名URL（需购买课程或管理员）

5.2.4 订单相关接口

接口路径

请求方法

请求参数

响应数据

说明

/orders

POST

{"course_id":1}

{"code":200,"msg":"购买成功","data":{"order_id":1}}

学生购买课程（需JWT认证）

5.3 gRPC接口（核心服务）

gRPC接口定义在`geekedu.proto`文件中，核心方法如下：

// 用户服务
rpc Register(RegisterRequest) returns (RegisterResponse);
rpc UserLogin(UserLoginRequest) returns (UserLoginResponse);

// 课程服务
rpc GetCourseList(CourseListRequest) returns (CourseListResponse);
rpc CreateCourse(CreateCourseRequest) returns (CreateCourseResponse);
rpc UploadVideo(UploadVideoRequest) returns (UploadVideoResponse);

// 订单服务
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);

// 播放服务
rpc GetCourseVideos(GetCourseVideosRequest) returns (GetCourseVideosResponse);
rpc GetVideoPlayUrl(GetVideoPlayUrlRequest) returns (GetVideoPlayUrlResponse);

6. 部署架构

采用Docker容器化部署，拆分前端、Web服务、gRPC服务、MySQL四个容器，通过Docker Compose编排，OSS采用阿里云托管服务，无需自建。

# docker-compose.yml核心配置
version: '3'
services:
  frontend:
    image: geekedu-frontend:v1
    ports:
      - "80:80"
    depends_on:
      - web-server

  web-server:
    image: geekedu-web-server:v1
    ports:
      - "8080:8080"
    environment:
      - GRPC_ADDR=grpc-server:50051
    depends_on:
      - grpc-server

  grpc-server:
    image: geekedu-grpc-server:v1
    environment:
      - DB_DSN=root:123456@tcp(mysql:3306)/geekedu
      - OSS_ENDPOINT=oss-cn-beijing.aliyuncs.com
      - OSS_ACCESS_KEY=xxx
      - OSS_SECRET_KEY=xxx
      - OSS_BUCKET=xxx
    depends_on:
      - mysql

  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=123456
      - MYSQL_DATABASE=geekedu
    volumes:
      - mysql-data:/var/lib/mysql

volumes:
  mysql-data:

7. 核心安全策略

- 身份认证：基于JWT实现用户登录认证，令牌有效期1小时，接口请求需携带Authorization头（Bearer Token）。

- 权限控制：管理员接口（发布课程、上传视频）校验角色为1，学生仅可购买、观看已购课程。

- 数据加密：用户密码通过BCrypt加密存储，传输过程采用HTTPS（生产环境）。

- 文件安全：OSS文件私有访问，仅通过后端生成的预签名URL临时访问，避免文件泄露。