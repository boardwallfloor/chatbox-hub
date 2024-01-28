package main

import (
	pbAuth "boardwallfloor/auth_module/pb/auth/v1"
	"boardwallfloor/auth_module/sqlite"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type AuthService struct {
	service pbAuth.AuthServiceServer
}

func StartGrpcSvr(service pbAuth.AuthServiceServer) *grpc.Server {
	authService := AuthService{service: service}
	grpcSvr := grpc.NewServer()
	pbAuth.RegisterAuthServiceServer(grpcSvr, authService.service)
	return grpcSvr
}

func startsql3_test() {
	log.Println("Running sql3 server for test in port 8080")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("tcp error, :%s \n", err)
	}

	log.Println("init sql3 service")
	os.Remove("./auth.db")
	if err != nil {
		log.Fatalf("error dumping test db : %s", err)
	}
	sql3, err := sqlite.NewSqliteDB("./auth.db")
	if err != nil {
		log.Fatalf("sqlite3 init error : %s", err)
	}
	defer sql3.Db.Close()

	sqlStatements := `
		CREATE TABLE IF NOT EXISTS user (
			user_id INTEGER PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS session (
			session_id TEXT PRIMARY KEY,
			user_id INTEGER NOT NULL,
	 is_long_term BOOLEAN NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES user(user_id),
			UNIQUE (user_id, session_id)
		);
	`

	// Execute the SQL statements
	_, err = sql3.Db.Exec(sqlStatements)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("starting server")
	authService := AuthService{service: &sql3}
	grpcSvr := grpc.NewServer()
	pbAuth.RegisterAuthServiceServer(grpcSvr, authService.service)
	grpcSvr.Serve(lis)
}

func main() {
	// lis, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatalf("tcp error, :%s \n", err)
	// }
	// // inmem := inmem.InMemStorage{}
	// sql3, err := sqlite.NewSqliteDB("../sqlite3/auth.db")
	// if err != nil {
	// 	log.Fatalf("sqlite3 init error : %s", err)
	// }
	// defer sql3.Db.Close()
	// grpcSvr := StartGrpcSvr(&sql3)
	// grpcSvr.Serve(lis)
	startsql3_test()
}
