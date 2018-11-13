FROM golang:1.11-alpine

EXPOSE 8080

COPY common.env Gopkg.lock Gopkg.toml GeoLite2-City-Blocks-IPv6.csv /go/src/github.com/patrickvalle/heatmap/
COPY ./cmd/apid /go/src/github.com/patrickvalle/heatmap/cmd/apid
COPY ./internal /go/src/github.com/patrickvalle/heatmap/internal

WORKDIR /go/src/github.com/patrickvalle/heatmap

RUN apk update && \
    apk add git && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure
