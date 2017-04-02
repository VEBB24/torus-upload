package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

func fetchFile(connection net.Conn) {
	glog.Infoln("A client has connected!")
	defer connection.Close()

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	bufferFileHash := make([]byte, 32)

	connection.Read(bufferFileHash)
	fileHash := string(bufferFileHash)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		glog.Errorln(err.Error())
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			fmt.Println("here")
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			break
		}

		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}

	newFile.Close()

	f, _ := os.Open(fileName)

	newHash, _ := generateHash(f)
	glog.Infoln(newHash)
	glog.Infoln(fileHash)
	f.Close()
	if newHash == fileHash {
		glog.Infoln("Received file completely!")
		connection.Write([]byte("OK\n"))
	} else {
		glog.Errorln("Error while fetching")
		connection.Write([]byte("ERROR\n"))
	}

}
