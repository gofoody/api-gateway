PROJECT_ROOT?=$(shell pwd)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o .out/api-gateway

.PHONY: build-ci
build-ci: build
	docker build -f Dockerfile -t api-gateway:latest .

.PHONY: clean
clean:
	-rm -rf .out/
