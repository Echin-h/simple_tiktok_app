# 使用的基础镜像
FROM golang:1.16

# 设置工作目录
WORKDIR /app

# 将应用程序依赖拷贝到工作目录
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 拷贝其他源代码  
COPY . .

# 构建应用程序
RUN go build -o main .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]