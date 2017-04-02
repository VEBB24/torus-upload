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
	if err != nil {
		glog.Errorln(err)
		os.Exit(1)
	}

	defer server.Close()

	glog.Infoln("Server started, waiting for connections...")

	for {
		connection, err := server.Accept()
		if err != nil {
			glog.Errorln(err)
			os.Exit(1)
		}
		glog.Infoln("Client connected")

		go fetchFile(connection)
	}

}
