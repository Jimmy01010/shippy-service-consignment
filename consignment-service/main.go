package main

import (
	// Import the generated protobuf code
	pb "github.com/Jimmy01010/shippy-service-consignment/consignment-service/proto/consignment"
	vesselProto "github.com/Jimmy01010/shippy-service-consignment/vessel-service/proto/vessel"
	"go-micro.dev/v4"
	"log"
)

func main() {

	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
	)
	// initialise flags
	srv.Init()

	vesselClient := vesselProto.NewVesselService("go.micro.srv.vessel", srv.Client())

	// Register service
	if err := pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient}); err != nil {
		log.Panic(err)
	}

	// start the service
	// Run the server
	if err := srv.Run(); err != nil {
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
