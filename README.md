# shippy-service-consignment

本项目包含下列微服务:
* consignment-service（货运服务）
* inventory-service（仓库服务）
* user-service（用户服务）
* authentication-service（认证服务）
* role-service （角色服务）
* vessel-service（货船服务）

系列一:
1. 创建了一个consignments (货运)服务和与之交互的客户端;
2. 实现简单的gRPC通信。

系列二:
1. 使用Docker容器来运行我们的服务;
   1. 构建客户端镜像 `docker build -t shippy-cli-consignment -f Dockerfile ..`
   2. 运行 `docker run shippy-cli-consignment`
   3. 构建服务端镜像 `docker build -t shippy-service-consignment .`
   4.  `docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 shippy-service-consignment`
2. 使用go-micro框架来进行服务发现。
3. Vessel 服务
   1. 构建服务端镜像 `docker build -t vessel-service .`

系列三:
1. 使用docker-compose来管理容器
   1. 停止所有当前运行的所有容器: `docker stop $(docker ps -qa)`
   2. 构建镜像时不适用缓存：`docker-compose build --no-cache`
2. 使用mongo db替代本地文件存储


上诉服务都在docker中启动:
Docker 有自己独立的 mdns，与宿主主机 Mac 的 mdns 不一致。把客户端也 Docker 化，这样服务端与客户端就在同一个网络层下，顺利使用 mdns 做服务发现
