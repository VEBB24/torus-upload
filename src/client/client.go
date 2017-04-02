package main

import (
	"flag"
	"io"
	"net"
	"os"
	"strconv"

	"bufio"

	tm "github.com/buger/goterm"
	"github.com/golang/glog"
)

const BUFFERSIZE = 1024

func main() {

	path := flag.String("file", "file.jp2", "a file")
	host := flag.String("host", "localhost", "a valid hostname")
	port := flag.String("port", "3784", "a valid port")
	flag.Parse()

	connection, err := net.Dial("tcp", *host+":"+*port)
	connectionStatus, errStatus := net.Dial("tcp", *host+":5050")
	if err != nil {
		panic(err)
	}

	if errStatus != nil {
		glog.Errorln(errStatus)
		os.Exit(1)
	}

	defer connection.Close()
	defer connectionStatus.Close()

	file, err := os.Open(*path)
	if err != nil {
		glog.Errorln(err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		glog.Errorln(err)
		return
	}

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	glog.Infoln("Sending filename and filesize to server")

	hash, _ := generateHash(file)

	file, err = os.Open(*path)
	if err != nil {
		glog.Errorln(err)
		return
	}

	connection.Write([]byte(hash))
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))

	stopchan := make(chan bool, 0)

	go func() {
		for {
			a, _ := bufio.NewReader(connectionStatus).ReadString('\n')
			tm.Clear()
			tm.MoveCursor(1, 1)
			tm.Println(tm.Color(tm.Bold(a), tm.GREEN))
			tm.Flush()
			if a == "100 %\n" {
				stopchan <- true
				break
			}
		}
	}()

	sendBuffer := make([]byte, BUFFERSIZE)
	glog.Infoln("Start sending file")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	for {
		if <-stopchan {
			break
		}
	}
	defer connectionStatus.Close()
	tm.Println(tm.Color(tm.Bold("File has been sent"), tm.GREEN))
	tm.Flush()
	glog.Infoln("File has been sent, closing connection!")
	return
}
