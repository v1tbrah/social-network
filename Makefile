start:
	docker-compose up -d --scale api-gateway=3

stop:
	docker-compose down

dockgen:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -d ./internal/api --parseDependency -g router.go
	swag fmt
