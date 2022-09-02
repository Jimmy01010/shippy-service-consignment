package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"log"

	userProto "github.com/Jimmy01010/protocol/shippy-user"
)

func createUser(ctx context.Context, service micro.Service, user *userProto.User) error {
	client := userProto.NewUserService("shippy.service.user", service.Client())
	rsp, err := client.Create(ctx, user)
	if err != nil {
		return err
	}

	// print the response
	fmt.Println("Response: ", rsp.User)

	return nil
}

func main() {
	// create and initialise a new service
	service := micro.NewService()
	service.Init(
		micro.Action(func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			company := c.String("company")
			password := c.String("password")

			ctx := context.Background()
			user := &userProto.User{
				Name:     name,
				Email:    email,
				Company:  company,
				Password: password,
			}

			if err := createUser(ctx, service, user); err != nil {
				log.Println("error creating user: ", err.Error())
				return err
			}

			return nil
		}),
	)
}
