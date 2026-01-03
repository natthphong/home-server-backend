package grpc_server

import (
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	gauth "github.com/natthphong/home-server-backend/grpc/auth"
	"github.com/natthphong/home-server-backend/grpc_server/auth"
	"google.golang.org/grpc"
)

func StartGRPCServer(db *pgxpool.Pool, jwtSecret string, accessTokenDuration, refreshTokenDuration time.Duration) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	gauth.RegisterAuthServiceServer(grpcServer, &auth.AuthServiceServer{
		DB:                   db,
		JWTSecret:            jwtSecret,
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	})

	fmt.Println("Starting gRPC server on port 8081")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
