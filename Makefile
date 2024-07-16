#!make
include .development.env
export $(shell sed 's/=.*//' .development.env)

test_env:
	env

start:
	@echo "Starting the server"
	@go run *.go