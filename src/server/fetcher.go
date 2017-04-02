package main

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

func fetchFile(connection net.Conn, connectionStatus net.Conn) {
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
		panic(err)
	}
	defer newFile.Close()

	var receivedBytes int64
	ticker := time.NewTicker(time.Millisecond * 50)
	status := make(chan float64)
	stopchan := make(chan bool, 0)

	go func() {
		for range ticker.C {
			select {
			case <-stopchan:
				ticker.Stop()
				return
			default:
				b := (float64(receivedBytes) / float64(fileSize)) * 100.0
				status <- b
			}
		}
	}()

	go func() {
		for {
			select {
			case percent := <-status:
				connectionStatus.Write([]byte(strconv.FormatFloat(percent, 'f', -1, 32) + " %" + "\n"))
				glog.Infoln(percent, "%")
			default:
			}
		}
	}()

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			stopchan <- true
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
	} else {
		glog.Errorln("Error while fetching")
	}

}
