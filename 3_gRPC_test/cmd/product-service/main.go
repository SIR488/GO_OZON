package main

import (
	"MyMod/cmd/product-service/service"
	"MyMod/pkg"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to Listen: %v", err)
	}
	server := grpc.NewServer()
	reflection.Register(server)
	svc, err := service.NewProductService()
	if err != nil {
		log.Fatalf("failed to create product service: %v", err)
	}
	pkg.RegisterProductServer(server, svc)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
