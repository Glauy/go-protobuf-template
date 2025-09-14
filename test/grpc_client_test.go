package cmd_test

import (
	"context"
	"log"
	"testing"
	"time"

	userpb "go-protos/proto/userpb"

	"google.golang.org/grpc"
)

func TestGrpcClient(t *testing.T) {
	// 1. 建立 gRPC 连接
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// 2. 创建 UserService client
	client := userpb.NewUserServiceClient(conn)

	// 3. 调用 CreateUser
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Username:     "alice",
		Email:        "alice@example.com",
		PasswordHash: "hashedpassword123",
	})
	if err != nil {
		log.Fatalf("CreateUser error: %v", err)
	}

	log.Printf("Created user: %+v\n", createResp.User)

	// 4. 调用 GetUserByUsername
	getResp, err := client.GetUserByUsername(ctx, &userpb.GetUserByUsernameRequest{
		Username: "alice",
	})
	if err != nil {
		log.Fatalf("GetUserByUsername error: %v", err)
	}

	log.Printf("Fetched user: %+v\n", getResp.User)
}
