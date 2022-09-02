FROM golang:alpine as builder

# 我们不使用apt-get install来安装包，而是使用apk add --no cache package 其中--no cache选项意味着下载的包不会被缓存。
# 同样地，我们将使用apk del而不是apt-get remove来卸载包。
RUN apk update &&  apk add --no-cache git

RUN mkdir /app

WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shippy-service-user *.go


# Run container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/shippy-service-user .

CMD ["./shippy-service-user"]