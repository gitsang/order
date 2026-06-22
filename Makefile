.PHONY: all init proto dev dev-server dev-web build build-server build-web docker-build docker-build-server docker-build-web docker-build-allinone docker-up docker-down lint test migrate-up migrate-down migrate-create

all: build

init:
	go mod tidy
	cd web && pnpm install
	podman-compose -f compose.yml up -d postgres

proto:
	protoc --go_out=. --go-grpc_out=. api/order/v1/*.proto
	protoc --grpc-gateway_out=. api/order/v1/*.proto

dev:
	make -j2 dev-server dev-server

dev-server:
	go run cmd/server/main.go

dev-web:
	cd web && pnpm dev

build: build-server build-web

build-server:
	go build -o bin/server cmd/server/main.go

build-web:
	cd web && pnpm build

docker-build: docker-build-server docker-build-web docker-build-allinone

docker-build-server:
	podman build -t order-server -f Containerfile.order .

docker-build-web:
	podman build -t order-web -f Containerfile.order-web .

docker-build-allinone:
	podman build -t order -f Containerfile .

docker-up:
	podman-compose -f compose.yml --profile allinone up -d

docker-down:
	podman-compose -f compose.yml down

lint:
	golangci-lint run
	cd web && pnpm lint

test:
	go test ./...
	cd web && pnpm test

migrate-up:
	migrate -path scripts/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path scripts/migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir scripts/migrations $(NAME)

seed:
	go run cmd/seed/main.go -admin-password $(ADMIN_PASSWORD)
