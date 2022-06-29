# shippy-service-consignment
系列一:
1. 创建了一个consignments (货运)服务和与之交互的客户端;
2. 实现简单的gRPC通信。

系列二:
1. 使用Docker容器来运行我们的服务;
   1. 构建客户端镜像 `docker build -t shippy-cli-consignment -f Dockerfile ..`
   2. 运行 `docker run shippy-cli-consignment`
2. 使用go-micro框架来进行服务发现。

