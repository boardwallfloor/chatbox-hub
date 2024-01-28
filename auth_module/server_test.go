package main

import (
	pbauth "boardwallfloor/auth_module/pb/auth/v1"
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSqlServer(t *testing.T) {
	go startsql3_test()

	timer := time.After(1 * time.Minute)
	go func() {
		<-timer
		log.Fatal("timeout error, time exceed 1m")
	}()

	log.Println("starting grpc client")
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc dial error: %s\n", err)
	}

	client := pbauth.NewAuthServiceClient(conn)

	t.Run("sign up", func(t *testing.T) {
		req := pbauth.SignUpRequest{Username: "user1", Password: "pass1", Email: "test@mail.com"}
		ctx := context.Background()
		res, err := client.SignUp(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		log.Println(res)
	})
	var sessionId string
	t.Run("sign in", func(t *testing.T) {
		req := pbauth.SignInRequest{Username: "user1", Password: "pass1"}
		ctx := context.Background()
		res, err := client.SignIn(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		log.Println(res)
		sessionId = res.SessionId
	})
	t.Run("is logged in", func(t *testing.T) {
		req := pbauth.CheckLoggedInRequest{SessionId: sessionId}
		ctx := context.Background()
		res, err := client.CheckLoggedIn(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		log.Println(res)
	})

	t.Run("sign up duplicate", func(t *testing.T) {
		req := pbauth.SignUpRequest{Username: "user1", Password: "pass1", Email: "test@mail.com"}
		ctx := context.Background()
		res, err := client.SignUp(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		log.Println(res)
	})

	t.Run("log out", func(t *testing.T) {
		req := pbauth.LogoutRequest{Session: sessionId}
		ctx := context.Background()
		res, err := client.Logout(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		log.Println(res)
	})

	t.Run("log out non existent account", func(t *testing.T) {
		req := pbauth.LogoutRequest{Session: sessionId}
		ctx := context.Background()
		res, err := client.Logout(ctx, &req)
		if err != nil {
			t.Fatalf("request error  : %s\n", err)
		}
		if res.Success {
			t.Errorf("failed log out response, expected %t but received %t", false, res.Success)
		}
	})
}
