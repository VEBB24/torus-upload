package main

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"path/filepath"

	"bufio"

	"github.com/golang/glog"
)

func fetchFile(connection net.Conn, basePath string, redis *Redis) {
	glog.Infoln("A client has connected!")
	defer connection.Close()
	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)
	var user string
	baseDir, err := filepath.Abs(basePath)
	if err != nil {
		glog.Errorln(err.Error())
		return
	}

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	bufferFileHash := make([]byte, 32)

	id, _ := reader.ReadString('\n')
	id = strings.TrimRight(id, "\n")
	glog.Infoln("Id = " + id)

	user = redis.GET(id)

	glog.Infoln("User = " + user)

	if user == "" {
		glog.Errorln("Unknown id")
		writer.WriteString("error@Unknown id")
		writer.Flush()
		return
	}

	writer.WriteString("success@client connected to " + connection.LocalAddr().String() + "\n")
	writer.Flush()

	connection.Read(bufferFileHash)
	fileHash := string(bufferFileHash)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	searchPath := filepath.Join(baseDir, "/", user)

	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		os.MkdirAll(searchPath, os.ModePerm)
	}

	fileName = filepath.Join(searchPath, "/", fileName)

	newFile, err := os.Create(fileName)

	if err != nil {
		glog.Errorln(err.Error())
	}
	defer newFile.Close()

	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
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
