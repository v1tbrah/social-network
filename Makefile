start:
	docker-compose up -d --scale api-gateway=3

stop:
	docker-compose down

mockgen:
	go install github.com/vektra/mockery/v2@v2.24.0
	go generate ./...

dockgen:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -d ./internal/api --parseDependency -g router.go
	swag fmt

test:
	go test ./...

push_to_dockerhub:
	docker build -t v1tbrah/api-gateway:v1 .
	docker tag v1tbrah/api-gateway:v1 v1tbrah/api-gateway:v1-release
	docker push v1tbrah/api-gateway:v1-release
