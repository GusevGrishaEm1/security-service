package main

import (
	"context"
	"fmt"
	"log"
	"time"

	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.NewClient(
		"127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("did not sss: %v", err)
	}

	client := security_servicev1.NewAuthClient(conn)
	res, err := client.Login(ctx, &security_servicev1.LoginRequest{
		Email:    "test",
		Password: "test",
	})
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Print(res.Token)

	defer conn.Close()
}
