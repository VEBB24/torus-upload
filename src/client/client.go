package main

import (
	"flag"
	"io"
	"net"
	"os"
	"strconv"

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

	if err != nil {
		panic(err)
	}

	defer connection.Close()

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

	sendBuffer := make([]byte, BUFFERSIZE)
	glog.Infoln("Start sending file")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	buffer := make([]byte, 10)
	size, _ := connection.Read(buffer)
	response := string(buffer[:size])
	if response == "OK\n" {
		tm.Println(tm.Color(tm.Bold(response), tm.GREEN))
	} else {
		tm.Println(tm.Color(tm.Bold(response), tm.RED))
	}
	tm.Flush()
	glog.Infoln("File has been sent, closing connection!")
	return
}
