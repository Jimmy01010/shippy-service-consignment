package main

import (
	"context"
	"os"

	vesselProto "github.com/Jimmy01010/protocol/vessel-service"
	"go-micro.dev/v4"

	pb "github.com/Jimmy01010/protocol/consignment-service"
	//"go-micro.dev/v4/cmd/protoc-gen-micro/plugin/micro"
	"log"
)

func main() {
	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
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
	vesselClient := vesselProto.NewVesselService("shippy.service.client", service.Client())
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
