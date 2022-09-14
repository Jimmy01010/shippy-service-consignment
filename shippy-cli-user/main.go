package main

import (
	"context"
	"fmt"
	userProto "github.com/Jimmy01010/protocol/shippy-user"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"log"
	"os"
)

// createUser 创建一个新用户
func createUser(ctx context.Context, service micro.Service, user *userProto.User) error {
	client := userProto.NewUserService("shippy.service.user", service.Client())
	rsp, err := client.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create a new user: %s", err.Error())
	}

	// print the response
	fmt.Println("a new user is created: ", rsp.User)

	authResponse, err := client.Auth(context.TODO(), &userProto.User{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", user.Email, err)
	}
	log.Printf("Your access token is: %s \n", authResponse.Token)

	return nil
}

func main() {
	// create and initialise a new service
	service := micro.NewService(
		micro.Flags(
			&cli.StringFlag{
				Name:  "name",
				Usage: "Your Name",
			},
			&cli.StringFlag{
				Name:  "email",
				Usage: "E-Mail",
			},
			&cli.StringFlag{
				Name:  "company",
				Usage: "Company Name",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Password",
			}))
	service.Init(
		micro.Action(func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			company := c.String("company")
			password := c.String("password")

			log.Println("test:", name, email, company, password)

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
			log.Printf("Created: %s", user.Name)

			return nil
		}),
	)
	// let's just exit because
	os.Exit(0)
}
