APID_DIR := "cmd/apid"
WEB_DIR := "cmd/web"
IPV6_INTERNAL_DIR := "internal/ipv6"

default: ensure up

ensure:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/protobuf/protoc-gen-go
	dep ensure
	npm install

generate:
	protoc --go_out=. ${IPV6_INTERNAL_DIR}/model.proto
	./node_modules/protobufjs/bin/pbjs -t json ${IPV6_INTERNAL_DIR}/model.proto > ${WEB_DIR}/js/proto-bundle.json

lint:
	./node_modules/.bin/eslint ${WEB_DIR}/js

test:
	go test ./... -cover

build: generate
	docker-compose build

up: stop build
	docker-compose up

stop:
	docker-compose stop
	
push:
	# Push apid to Docker
	docker build -t patrickvalle/heatmap-apid:latest -f ${APID_DIR}/Dockerfile .
	docker login
	docker push patrickvalle/heatmap-apid:latest