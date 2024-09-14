#!make
include .development.env
export $(shell sed 's/=.*//' .development.env)

test_env:
	env

start:
	@echo "Starting the server"
	@go run *.go

build:
	@echo "Building the server"
	@go mod tidy && go mod download && go build -v -o engine && chmod +x engine

run:
	@echo "Injecting environment variables"
	@echo "Running the server"
	@./engine