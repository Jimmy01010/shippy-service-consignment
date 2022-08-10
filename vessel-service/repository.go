package main

import (
	"context"
	"fmt"

	pb "github.com/Jimmy01010/protocol/vessel-service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository interface {
	FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error)
	Create(ctx context.Context, vessel *Vessel) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

type Specification struct {
	Capacity  int32
	MaxWeight int32
}

func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func UnmarshalSpecification(spec *Specification) *pb.Specification {
	return &pb.Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

func MarshalVessel(vessel *pb.Vessel) *Vessel {
	return &Vessel{
		ID:        vessel.Id,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerID:   vessel.OwnerId,
	}
}

func UnmarshalVessel(vessel *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id:        vessel.ID,
		Capacity:  vessel.Capacity,
		MaxWeight: vessel.MaxWeight,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerId:   vessel.OwnerID,
	}
}

type Vessel struct {
	ID        string `bson:"ID" json:"ID"`
	Capacity  int32  `bson:"capacity" json:"capacity"`
	Name      string `bson:"name" json:"name"`
	Available bool   `bson:"available" json:"available"`
	OwnerID   string `bson:"ownerID" json:"ownerID"`
	MaxWeight int32  `bson:"maxWeight" json:"maxWeight"`
}

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
// FindAvailable - 寻找可用的船(Vessel)
// 如果货物的容量和最大重量低于船舶的容量和最大重量，则返回该船舶。
func (repository *MongoRepository) FindAvailable(ctx context.Context, spec *Specification) (*Vessel, error) {
	// 设置过滤条件
	// 感觉这个过滤有些问题？为啥一个capacity key,设置两个过滤条件？
	filter := bson.D{{
		"capacity",
		bson.D{{
			"$lte",
			spec.Capacity,
		}, {
			"$lte",
			spec.MaxWeight,
		}},
	}}
	vessel := &Vessel{}
	err := repository.collection.FindOne(ctx, filter).Decode(vessel)
	if err == mongo.ErrNoDocuments {
		// Do something when no record was found
		return vessel, fmt.Errorf("record does not exist on map of vessels")
	} else if err != nil {
		return vessel, err
	}
	return vessel, nil
}

// Create a new vessel
func (repository *MongoRepository) Create(ctx context.Context, vessel *Vessel) error {
	//_, err := repository.collection.InsertOne(ctx, vessel)
	upsert := true

	filter := bson.D{{
		"ID",
		bson.M{"ID": "vessel001"},
	}}
	_, err := repository.collection.UpdateOne(ctx, filter, vessel, &options.UpdateOptions{Upsert: &upsert})
	return err
}

//func (repository *MongoRepository) Upsert() {
//	_, err := repository.collection.UpdateOne()
//	return err
//}
