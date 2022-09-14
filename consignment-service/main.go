package main

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"

	//"github.com/micro/go-micro/v2/client"
	//"github.com/micro/go-micro/v2/metadata"
	//"github.com/micro/go-micro/v2/server"

	"os"

	pb "github.com/Jimmy01010/protocol/consignment-service"
	userPb "github.com/Jimmy01010/protocol/shippy-user"
	vesselProto "github.com/Jimmy01010/protocol/vessel-service"
	//"go-micro.dev/v4/cmd/protoc-gen-micro/plugin/micro"
	"log"
)

func main() {
	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
		micro.WrapHandler(AuthWrapper),
	)
	// initialise flags
	service.Init()

	// 创建mongo客户端
	var uri string
	if uri = os.Getenv("DB_HOST"); uri == "" {
		log.Fatal("'DB_HOST' is empty You must set your 'DB_HOST' environmental variable.")
	}
	client, err := CreateMongoClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// 连接mongo服务端
	consignmentCollection := client.Database("shippy").Collection("consignment")
	repository := &MongoRepository{consignmentCollection}

	// 创建一个货船服务？ 为啥这里这个名字不是vessel服务名
	// vesselClient := vesselProto.NewVesselService("shippy.service.client", service.Client())
	vesselClient := vesselProto.NewVesselService("shippy.service.vessel", service.Client())
	h := &handler{repository, vesselClient}

	// Register service
	if err := pb.RegisterShippingServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	// start the service
	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}

	/*	// Set-up our gRPC server.
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		// Register our service with the gRPC server, this will tie our
		// implementation into the auto-generated interface code for our
		// protobuf definition.
		pb.RegisterShippingServiceServer(s, &service{repo})

		// Register reflection service on gRPC server.
		// 如果启动了gprc反射服务，那么就可以通过reflection包提供的反射服务查询gRPC服务或调用gRPC方法。
		// 从而在通过gRPC CLI 工具，我们可以在没有客户端代码的环境下测试gRPC服务。
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}*/
}

// AuthWrapper 是一个高阶函数，入参是 ”下一步“ 函数，出参是认证函数
// 在返回的函数内部处理完认证逻辑后，再手动调用 fn() 进行下一步处理
// token 是从 consignment-ci 上下文中取出的，再调用 user-service 将其做验证
// 认证通过则 fn() 继续执行，否则报错
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		// consignment-service 独立测试时不进行认证
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]

		// Auth here
		authClient := userPb.NewUserService("shippy.service.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(context.Background(), &userPb.Token{
			Token: token,
		})
		log.Println("ValidateToken Resp:", authResp)
		if err != nil {
			return fmt.Errorf("validateToken failed: %s", err.Error())
		}
		err = fn(ctx, req, resp)
		return err
	}
}
