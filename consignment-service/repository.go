package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	pb "github.com/Jimmy01010/protocol/consignment-service"
)

// Consignment pb的Consignment映射
type Consignment struct {
	ID          string     `json:"id"`
	Weight      int32      `json:"weight"`
	Description string     `json:"description"`
	Containers  Containers `json:"containers"`
	VesselID    string     `json:"vessel_id"`
}

type Containers []*Container

type Container struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	UserID     string `json:"user_id"`
}

//  marshal 和 unmarshal
//	用于在 protobuf 定义生成的结构和我们的内部数据存储模型之间进行转换。
//	理论上，您也可以将protobuf生成的结构用作您的数据model, 但从软件设计的角度来看，这不是必要的推荐。
//	因为您现在将数据模型耦合到您的交付层(你的通信和数据存储之间就有很高的耦合性)。
//	保持软件中各种职责之间的这些界限是很好的。 这似乎是额外的开销，但这对于您的软件的可扩展性很重要。

// MarshalConsignment an input consignment type to a consignment model
func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	return &Consignment{
		ID:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  marshalContainerCollection(consignment.Containers),
		VesselID:    consignment.VesselId,
	}
}

// UnmarshalConsignmentCollection 将存储类型转换为通信传输类型
func UnmarshalConsignmentCollection(consignments []*Consignment) []*pb.Consignment {
	collection := make([]*pb.Consignment, 0)
	for _, consignment := range consignments {
		collection = append(collection, unmarshalConsignment(consignment))
	}
	return collection
}

func marshalContainerCollection(containers []*pb.Container) []*Container {
	collection := make([]*Container, 0)
	for _, container := range containers {
		collection = append(collection, marshalContainer(container))
	}
	return collection
}

func marshalContainer(container *pb.Container) *Container {
	return &Container{
		ID:         container.Id,
		CustomerID: container.CustomerId,
		UserID:     container.UserId,
	}
}

func unmarshalConsignment(consignment *Consignment) *pb.Consignment {
	return &pb.Consignment{
		Id:          consignment.ID,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  unmarshalContainerCollection(consignment.Containers),
		VesselId:    consignment.VesselID,
	}
}

func unmarshalContainerCollection(containers []*Container) []*pb.Container {
	collection := make([]*pb.Container, 0)
	for _, container := range containers {
		collection = append(collection, unmarshalContainer(container))
	}
	return collection
}

func unmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.ID,
		CustomerId: container.CustomerID,
		UserId:     container.UserID,
	}
}

type repository interface {
	Create(ctx context.Context, consignment *Consignment) error
	GetAll(ctx context.Context) ([]*Consignment, error)
}

// MongoRepository implementation
type MongoRepository struct {
	collection *mongo.Collection
}

// Create -
func (repository *MongoRepository) Create(ctx context.Context, consignment *Consignment) error {
	_, err := repository.collection.InsertOne(ctx, consignment)
	return err
}

// GetAll -
func (repository *MongoRepository) GetAll(ctx context.Context) ([]*Consignment, error) {
	cur, err := repository.collection.Find(ctx, nil, nil)
	var consignments []*Consignment
	for cur.Next(ctx) {
		var consignment *Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, consignment)
	}
	return consignments, err
}

//func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
//	updated := append(repo.consignments, consignment)
//	repo.consignments = updated
//	return consignment, nil
//}
//
//func (repo *Repository) GetAll() []*pb.Consignment {
//	return repo.consignments
//}
