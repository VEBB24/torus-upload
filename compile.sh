#!/bin/bash

go get "github.com/golang/glog"
go get "github.com/buger/goterm"
go get "github.com/mediocregopher/radix.v2"
go get "gopkg.in/resty.v0"
go get "github.com/colinmarc/hdfs"

mkdir -p build

go build -o build/client src/client/*.go
go build -o build/server src/server/*.go
go build -o build/hdfs src/hdfs/*.go