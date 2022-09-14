package main

import (
	"context"
	"encoding/json"
	pb "github.com/Jimmy01010/protocol/consignment-service"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
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
	service := micro.NewService(micro.Name("shippy.cli.consignment"))
	service.Init()

	client := pb.NewShippingService("shippy.service.consignment", service.Client())

	// Set up a connection to the server.
	//conn, err := grpc.Dial(address, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("Did not connect: %v", err)
	//}
	//defer conn.Close()

	// client := pb.NewShippingServiceClient(conn)

	// Contact the server and print out its response.
	file := defaultFilename
	// 9/7修改： 在命令行中加token信息, 然后把这个token放到rpc请求的context中，准备在服务端验证
	// 在命令行中指定新的货物信息 json 件
	// Contact the server and print out its response.
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	var token string
	if token = os.Getenv("TOKEN"); token == "" {
		log.Fatal("'TOKEN' is empty You must set your 'TOKEN' environmental variable.")
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}
	// 创建带有用户 token 的 context
	// consignment-service 服务端将从中取出 token，解密取出用户身份
	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})
	// 调用 RPC
	// 将货物存储到指定用户的仓库里
	r, err := client.CreateConsignment(tokenContext, consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Printf("Consignment created: %t", r.Created)
	// 列出目前所有托运的货物
	getAll, err := client.GetConsignments(tokenContext, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
