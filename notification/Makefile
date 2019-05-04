
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:${PWD}/..:. --micro_out=. --go_out=. proto/notification/notification.proto

.PHONY: build
build: proto

	go build -o notification-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t notification-srv:latest
