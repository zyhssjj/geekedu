# 阶段1：构建前端项目
FROM node:18-alpine AS builder
WORKDIR /app

# 核心修改：基于 geekedu-project 上下文，直接访问 frontend 目录的 package.json
COPY ./frontend/package.json .

# 安装依赖（使用国内镜像，加快构建）
RUN npm install --registry=https://registry.npmmirror.com

# 核心修改：基于 geekedu-project 上下文，直接复制 frontend 目录的所有源码
COPY ./frontend/ .

# 构建静态文件
RUN npm run build

# 阶段2：使用Nginx托管
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 核心修改：基于 geekedu-project 上下文，直接访问 deploy/nginx/nginx.conf
COPY ./deploy/nginx/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露80端口
EXPOSE 80

# 启动Nginx
CMD ["nginx", "-g", "daemon off;"]