# 第一阶段：构建阶段
FROM golang:1.23.5-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制源代码
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o davlin ./cmd/main.go

# 第二阶段：运行阶段
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装基本工具和SSL证书
RUN apk --no-cache add ca-certificates tzdata curl

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/davlin .

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露应用端口
EXPOSE 8080

# 运行应用
CMD ["./davlin"]