FROM golang:1.11-alpine

EXPOSE 8080

COPY ./ /go/src/github.com/patrickvalle/heatmap
WORKDIR /go/src/github.com/patrickvalle/heatmap

RUN apk update && \
    apk add git && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure
