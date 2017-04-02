package main

import (
	"flag"
	"net"
	"os"

	"github.com/golang/glog"
)

const BUFFERSIZE = 1024

func main() {
	host := flag.String("host", "localhost", "a file")
	port := flag.String("port", "3784", "a valid port")
	flag.Parse()

	server, err := net.Listen("tcp", *host+":"+*port)
	serverStatus, errStatus := net.Listen("tcp", *host+":5050")
	if err != nil {
		glog.Errorln(err)
		os.Exit(1)
	}
	if errStatus != nil {
		panic(errStatus)
	}
	defer server.Close()
	defer serverStatus.Close()

	glog.Infoln("Server started, waiting for connections...")

	for {
		connection, err := server.Accept()
		connectionStatus, errStatus := serverStatus.Accept()
		if err != nil {
			glog.Errorln(err)
			os.Exit(1)
		}
		if errStatus != nil {
			glog.Errorln(errStatus)
			os.Exit(1)
		}
		glog.Infoln("Client connected")

		go fetchFile(connection, connectionStatus)
	}

}
