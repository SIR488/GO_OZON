.PHONY: generate
generate:
	protoc \
	--go_out=./pkg \
	--go-grpc_out=./pkg \
	--go-grpc_opt=paths=source_relative \
	--go_opt=paths=source_relative *.proto

.PHONY: run-product
run-product:
	go run ./cmd/product-service

.PHONY: run-cart
run-cart:
	go run ./cmd/cart

.PHONY:  run-all
run-all:
	@$(MAKE) -j2 run-product run-cart
