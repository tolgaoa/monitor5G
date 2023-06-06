# Makefile

SRC_DIR := ./cmd

# Set the name of the Go module
MODULE_NAME := proxy

build:
	go build -o ./bin/proxyapp $(SRC_DIR)/main.go

image:
	docker build . -t tolgaomeratalay/revprox:v1

run:
	go run $(SRC_DIR)/main.go

deps:
	go mod download

clean:
	rm -f ./bin/proxyapp

.PHONY: build run deps clean
