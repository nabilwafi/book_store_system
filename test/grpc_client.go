package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create user service client
	userClient := proto.NewUserServiceClient(conn)

	fmt.Println("=== Testing gRPC User Service ===")
	fmt.Println()

	// Test Register
	fmt.Println("1. Testing User Registration")
	registerReq := &proto.RegisterRequest{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	registerResp, err := userClient.Register(ctx, registerReq)
	if err != nil {
		fmt.Printf("❌ Register failed: %v\n", err)
	} else {
		fmt.Printf("✅ Register successful!\n")
		fmt.Printf("   Message: %s\n", registerResp.Message)
		if registerResp.User != nil {
			fmt.Printf("   User ID: %d\n", registerResp.User.Id)
			fmt.Printf("   Name: %s\n", registerResp.User.Name)
			fmt.Printf("   Email: %s\n", registerResp.User.Email)
			fmt.Printf("   Role: %s\n", registerResp.User.Role)
		}
	}
	fmt.Println()

	// Test Login
	fmt.Println("2. Testing User Login")
	loginReq := &proto.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	loginResp, err := userClient.Login(ctx2, loginReq)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
	} else {
		fmt.Printf("✅ Login successful!\n")
		fmt.Printf("   Message: %s\n", loginResp.Message)
		fmt.Printf("   Token: %s\n", loginResp.Token)
		if loginResp.User != nil {
			fmt.Printf("   User ID: %d\n", loginResp.User.Id)
			fmt.Printf("   Name: %s\n", loginResp.User.Name)
			fmt.Printf("   Email: %s\n", loginResp.User.Email)
			fmt.Printf("   Role: %s\n", loginResp.User.Role)
		}
	}
	fmt.Println()

	// Test Login with wrong password
	fmt.Println("3. Testing Login with Wrong Password")
	wrongLoginReq := &proto.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "wrongpassword",
	}

	ctx3, cancel3 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel3()

	wrongLoginResp, err := userClient.Login(ctx3, wrongLoginReq)
	if err != nil {
		fmt.Printf("❌ Login failed (expected): %v\n", err)
	} else {
		fmt.Printf("⚠️  Login should have failed but succeeded: %s\n", wrongLoginResp.Message)
	}
	fmt.Println()

	// Test Register with existing email
	fmt.Println("4. Testing Register with Existing Email")
	dupRegisterReq := &proto.RegisterRequest{
		Name:     "Jane Doe",
		Email:    "john.doe@example.com", // Same email
		Password: "password456",
	}

	ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel4()

	dupRegisterResp, err := userClient.Register(ctx4, dupRegisterReq)
	if err != nil {
		fmt.Printf("❌ Register failed (expected): %v\n", err)
	} else {
		fmt.Printf("⚠️  Register should have failed but succeeded: %s\n", dupRegisterResp.Message)
	}
	fmt.Println()

	fmt.Println("=== gRPC Testing Complete ===")
}