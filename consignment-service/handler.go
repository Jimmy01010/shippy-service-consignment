package main

import (
	"context"
	"fmt"
	pb "github.com/Jimmy01010/protocol/consignment-service"
	vesselProto "github.com/Jimmy01010/protocol/vessel-service"
	"log"
)

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
//type Repository struct {
//	consignments []*pb.Consignment
//}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
// 实现服务中定义的方法
type handler struct {
	repository
	vesselClient vesselProto.VesselService
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	// 为本次请求的货物寻找可用的船
	vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	if vesselResponse == nil {
		return fmt.Errorf("error fetching vessel, returned nil")
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	// We set the VesselId as the vessel we got back from our vessel service
	// 我们将从(vessel)货船服务中找到的船ID(VesselID)设置货运服务的船
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	err = s.repository.Create(ctx, MarshalConsignment(req))
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = req
	return nil
}

// GetConsignments -
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments, err := s.repository.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Consignments = UnmarshalConsignmentCollection(consignments)
	return nil
}
