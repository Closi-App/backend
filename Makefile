.PHONY:
.SILENT:
.DEFAULT_GOAL := run

run:
	go mod tidy && go run cmd/server/main.go

swag:
	swag init -g cmd/server/main.go && swag fmt