// shippy-service-consignment/main.go
package main

import (
	"context"
	"log"
	"sync"

	// Import the generated protobuf code
	pb "github.com/Jimmy01010/shippy-service-consignment/proto/consignment"
	"go-micro.dev/v4"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
// 虚拟存储库，这模拟了某种数据存储的使用。稍后我们将用一个真正的实现来替换它。
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment // 运输的仓库货物
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// GetAll consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
// 服务应该实现所有的方法来满足我们在protobuf定义中定义的服务。
// 你可以在生成的代码中检查这个接口，看看是否有确切的方法签名等等，这样你会有更好的想法。
type service struct {
	pb.UnimplementedShippingServiceServer
	repo repository
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
// CreateConsignment——我们在服务上只创建了一个方法，这是一个create方法，它接受一个上下文和一个请求作为参数，这些都是由gRPC服务器处理的。
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// GetConsignments -
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	service1 := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("shippy.service.consignment"))
	// Register service
	// (service1.Server() 接口怎么存储注册的数据呢？
	if err := pb.RegisterShippingServiceHandler(service1.Server(), &consignmentService{repo}); err != nil {
		log.Panic(err)
	}
	// Init will parse the command line flags.
	service1.Init()

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}

	/*
		// Set-up our gRPC server.
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		// Register our service with the gRPC server, this will tie our
		// implementation into the auto-generated interface code for our
		// protobuf definition.
		// 向gRPC服务器注册我们的服务，这将把我们的实现绑定到protobuf定义的自动生成的接口代码中。
		pb.RegisterShippingServiceServer(s, &service{repo: repo})

		// Register reflection service on gRPC server.
		reflection.Register(s)

		log.Println("Running on port:", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	*/
}
