package main

import (
	"encoding/json"
	pb "github.com/Jimmy01010/shippy/shippy-service-consignment/proto/consignment"
	"go-micro.dev/v4"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	service := micro.NewService(micro.Name("shippy.consignment.cli"))
	service.Init()

	client := pb.NewShippingService("shippy.consignment.service", service.Client())

	// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("Did not connect: %v", err)
	//}
	//defer conn.Close()

	// client := pb.NewShippingServiceClient(conn)

	// Contact the server and print out its response.
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
