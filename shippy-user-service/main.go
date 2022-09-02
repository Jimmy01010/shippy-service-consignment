package main

import (
	pb "github.com/Jimmy01010/protocol/shippy-user"
	"go-micro.dev/v4"
	"log"
)

const schema = `
	create table if not exists users (
		id varchar(36) not null,
		name varchar(125) not null,
		email varchar(225) not null unique,
		password varchar(225) not null,
		company varchar(125),
		primary key (id)
	);
`

func main() {
	// Creates a database connection and handles
	// closing it again before exit.
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer db.Close()

	// Run schema query on start-up, as we're using "create if not exists"
	// this will only be ran once. In order to create updates, you'll need to
	// use a migrations library
	db.MustExec(schema)

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("shippy.service.user"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	service.Init()

	repo := NewPostgresRepository(db)
	tokenService := &TokenService{repo}
	// Register handler
	err = pb.RegisterUserServiceHandler(service.Server(), &handler{repo, tokenService})
	if err != nil {
		log.Panic(err)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
