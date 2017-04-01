package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const BUFFERSIZE = 1024

func main() {
	host := flag.String("host", "localhost", "a file")
	port := flag.String("port", "3784", "a valid port")
	flag.Parse()
	server, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error listetning: ", err)
		os.Exit(1)
	}
	defer server.Close()
	fmt.Println("Server started! Waiting for connections...")
	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("Client connected")
		go fetchFile(connection)
	}
}

func fetchFile(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

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
	status := make(chan float32)
	stopchan := make(chan bool, 0)

	go func() {
		for range ticker.C {
			select {
			case <-stopchan:
				file, err := newFile.Stat()
				if err != nil {
					fmt.Println(err.Error())
				}
				if file.Size() == fileSize {
					fmt.Println("Received file completely!")
				} else {
					fmt.Println("Error")
				}
				return
			default:
				b := (float32(receivedBytes) / float32(fileSize)) * 100.0
				status <- b
			}
		}
	}()

	for {
		select {
		case percent := <-status:
			fmt.Printf("%f %% \n", percent)
		default:
			if (fileSize - receivedBytes) < BUFFERSIZE {
				io.CopyN(newFile, connection, (fileSize - receivedBytes))
				connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
				stopchan <- true
				break
			}
			io.CopyN(newFile, connection, BUFFERSIZE)
			receivedBytes += BUFFERSIZE
		}
	}

}
