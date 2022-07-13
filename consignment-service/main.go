package main

import (
	// Import the generated protobuf code
	pb "github.com/Jimmy01010/shippy-service-consignment/consignment-service/proto/consignment"
	vesselProto "github.com/Jimmy01010/shippy-service-consignment/vessel-service/proto/vessel"
	"go-micro.dev/v4"
	"golang.org/x/net/context"
	"log"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         IRepository
	vesselClient vesselProto.VesselService
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
	)
	// initialise flags
	srv.Init()

	vesselClient := vesselProto.NewVesselService("shippy.service.client", srv.Client())

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
