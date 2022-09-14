package main

import (
	"context"
	pb "github.com/Jimmy01010/protocol/vessel-service"
	"go-micro.dev/v4"
	"log"
	"os"
)

func main() {
	// 初始化船(Vessel)的基本信息、规格
	// 初始化一些船用于测试
	vessels := []*Vessel{
		{
			ID:        "vessel001",
			Name:      "Boaty McBoatface",
			MaxWeight: 200000,
			Capacity:  500,
		},
		//{
		//	ID:        "vessel002",
		//	Name:      "Boaty Mc",
		//	MaxWeight: 400000,
		//	Capacity:  1000,
		//},
	}

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

	// MongoDB 将自动创建不存在的数据库或集合。
	vesselCollection := client.Database("shippy").Collection("vessels")
	repo := &MongoRepository{vesselCollection}
	h := &handler{repo}

	for _, v := range vessels {
		if err := repo.Create(context.Background(), v); err != nil {
			log.Panic("创建vessels失败:", err)
		}
	}

	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
