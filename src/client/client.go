package main

import (
	"flag"
	"io"
	"net"
	"os"
	"strconv"

	"path/filepath"

	"bufio"

	"strings"

	tm "github.com/buger/goterm"
	"github.com/golang/glog"
)

//BUFFERSIZE size for the socket buffer
const BUFFERSIZE = 1024

func main() {

	path := flag.String("file", "file.jp2", "a file")
	host := flag.String("host", "localhost", "a valid hostname")
	port := flag.String("port", "3784", "a valid port")
	id := flag.String("id", "1234abbc45", "a valid id")
	flag.Parse()

	connection, err := net.Dial("tcp", *host+":"+*port)

	if err != nil {
		panic(err)
	}

	defer connection.Close()
	p, err := filepath.Abs(*path)
	if err != nil {
		glog.Errorln(err.Error())
		os.Exit(1)
	}

	file, err := os.Open(p)
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

	connection.Write([]byte(*id + "\n"))

	isValid, _ := bufio.NewReader(connection).ReadString('\n')

	info := strings.Split(isValid, "@")
	msgType, msg := info[0], info[1]

	if msgType == "error" {
		tm.Println(tm.Color(tm.Bold(msg), tm.RED))
		tm.Flush()
		os.Exit(1)
	} else {
		tm.Println(tm.Color(tm.Bold(msg), tm.GREEN))
		tm.Flush()
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
	response, _ := bufio.NewReader(connection).ReadString('\n')
	response = strings.TrimRight(response, "\n")
	if response == "OK" {
		tm.Println(tm.Color(tm.Bold(response), tm.GREEN))
	} else {
		tm.Println(tm.Color(tm.Bold(response), tm.RED))
	}
	tm.Flush()
	glog.Infoln("File has been sent, closing connection!")
	return
}
