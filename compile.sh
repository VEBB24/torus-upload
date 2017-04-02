#!/bin/bash

go get "github.com/golang/glog"
go get "github.com/buger/goterm"

mkdir -p build

go build -o build/client src/client/*.go
go build -o build/server src/server/*.go