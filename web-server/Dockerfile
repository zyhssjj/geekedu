# 构建阶段
FROM golang:1.25.3-alpine AS builder

# 1. 创建上级目录（模拟本地项目根目录结构）
WORKDIR /app/parent

# 2. 复制 common 目录
COPY ./common ./common

# 3. 复制 proto 目录
COPY ./proto ./proto

# 4. 复制 web-server 目录
COPY ./web-server ./web-server

# 5. 进入 web-server 目录
WORKDIR /app/parent/web-server

# 关键：先配置 GOPROXY（移到前面，必须在 go mod download 之前）
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# 6. 再下载依赖（此时代理已生效，顺利下载第三方依赖）
RUN go mod download

# 7. 编译项目
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-server .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制编译好的二进制文件
COPY --from=builder /app/parent/web-server/web-server .

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./web-server"]