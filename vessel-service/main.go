package main

import (
	"context"
	pb "github.com/Jimmy01010/protocol/vessel-service"
	"go-micro.dev/v4"
	"log"
	"os"
)

//type Repository interface {
//	FindAvailable(spec *pb.Specification) (*pb.Vessel, error)
//}
//
//type VesselRepository struct {
//	vessels []*pb.Vessel
//}
//
//// Our grpc service handler
//type vesselService struct {
//	repo Repository
//}
//
//// FindAvailable 为寻找可用的船提供rpc方法
//func (s *vesselService) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
//
//	// Find the next available vessel
//	vessel, err := s.repo.FindAvailable(req)
//	if err != nil {
//		return err
//	}
//
//	// Set the vessel as part of the response message type
//	res.Vessel = vessel
//	return nil
//}
//
//// FindAvailable 寻找可用的船(Vessel)，如果货物容量和最大重量低于船舶容量和最大重量，则返回该船舶。
//func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
//	// 选择最近一条容量、载重都符合的货轮
//	for _, v := range repo.vessels {
//		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
//			return v, nil
//		}
//	}
//	return nil, errors.New("no vessel found by that spec")
//}

func main() {
	// 初始化船(Vessel)的基本信息、规格
	//vessels := []*pb.Vessel{
	//	{
	//		Id:        "vessel001",
	//		Name:      "Boaty McBoatface",
	//		MaxWeight: 200000,
	//		Capacity:  500,
	//	},
	//}

	// repo := &VesselRepository{vessels}

	service := micro.NewService(
		micro.Name("shippy.service.vessel"),
	)
	service.Init()
	uri := os.Getenv("DB_HOST")
	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())
	vesselCollection := client.Database("shippy").Collection("vessels")
	repo := &MongoRepository{vesselCollection}
	h := &handler{repo}

	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
