package main

import (
	"fmt"
	"log"
	"net"

	"github.com/through-this-dunya/finalProject/pkg/config"
	"github.com/through-this-dunya/finalProject/pkg/database"
	"github.com/through-this-dunya/finalProject/pkg/service"
	"github.com/through-this-dunya/finalProject/pkg/utility"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := database.Init(c.DBUrl)

	jwt := utility.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "finalProject",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := service.Server{
		Handler: h,
		Jwt:     jwt,
	}

	grpcServer := grpc.NewServer()

	database.Register(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
