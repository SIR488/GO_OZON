package main

import (
	httpserver "MyMod/cmd/cart/http-server"
	client "MyMod/cmd/cart/product-client"
	"MyMod/pkg"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	productClient := client.NewProductClient(pkg.NewProductClient(conn))
	service := httpserver.NewCartService(productClient)
	ui := httpserver.NewUI(service)

	http.HandleFunc("POST /user/{user_id}/cart/{sku_id}", ui.PostSKU)
	http.HandleFunc("GET /user/{user_id}/cart", ui.GetCart)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("failed to ListenAndServe: %v", err)
	}
}
