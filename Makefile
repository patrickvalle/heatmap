APID_DIR := "cmd/apid"

default: ensure generate test

ensure:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/protobuf/protoc-gen-go
	dep ensure

generate:
	protoc --go_out=. internal/ipv6/model.proto

test:
	go test ./... -cover
