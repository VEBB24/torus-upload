package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const BUFFERSIZE = 1024

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

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
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return
}
