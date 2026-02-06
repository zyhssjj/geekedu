# 在线视频学习平台（geekedu-project）
## 项目介绍
基于Go微服务（Gin+gRPC）+ React的前后端分离项目，实现付费视频课程的上传、购买、时效播放，核心采用阿里云OSS私有Bucket保护资源。

## 前置准备
1.  安装Go 1.20+、Node.js 18+、Docker & Docker Compose
2.  开通阿里云OSS，创建私有Bucket，获取AccessKey ID、SecretKey、Endpoint、Bucket名称
3.  填写项目根目录的`.env`文件，替换阿里云OSS配置

## 快速启动
### 1. 克隆项目（本地已创建可忽略）
```bash
git clone [仓库地址]
cd geekedu-project