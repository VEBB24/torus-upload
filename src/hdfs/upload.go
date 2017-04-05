package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"

	"os"

	"path/filepath"

	"os/exec"

	"strings"

	tm "github.com/buger/goterm"
	"github.com/colinmarc/hdfs"
	"github.com/golang/glog"
)

func main() {
	host := flag.String("host", "127.0.0.1", "a valid host")
	path := flag.String("file", "", "a file")
	name := flag.String("name", "test.txt", "a valid name")
	token := flag.String("token", "1", "a valid id")
	flag.Parse()
	if *path == "" {
		exec.Command("./hdfs -h").Run()
		os.Exit(1)
	}

	redis := RedisFactory(*host)

	redis.SET("1", "paul")
	redis.SET("2", "thomas")
	redis.SET("3", "john")

	user := redis.GET(*token)

	if user == "" {
		os.Exit(1)
	}

	pushFile(user, *path, *name, *host)

}

func ask(r *bufio.Reader, name string) bool {
	var response string
	for {
		fmt.Println("The file " + name + " already exist, do you want to replace it ? (yes/no)")
		response, _ = r.ReadString('\n')
		response = strings.TrimRight(response, "\n")
		if response == "yes" || response == "no" {
			break
		}
	}
	return response == "yes"
}

func pushFile(user string, path string, name string, host string) {
	dirPath := filepath.Join("/user/admin/", user)
	filePath := filepath.Join(dirPath, "/", name)
	client, err := hdfs.New(host + ":8020")

	if err != nil {
		fmt.Println(err.Error())
	}

	client.MkdirAll(dirPath, 0755)
	_, e := client.Stat(filePath)

	if e == nil {
		reader := bufio.NewReader(os.Stdin)
		if !ask(reader, name) {
			os.Exit(0)
		}
		err := client.Remove(filePath)
		if err != nil {
			glog.Errorln(err.Error())
			os.Exit(1)
		}
	}

	w, err := client.CreateFile(filePath, 1, 1048576, 0755)
	if err != nil {
		glog.Errorln(err.Error())
	}

	file, _ := os.Open(path)

	_, errFile := io.Copy(w, file)

	if errFile != nil {
		glog.Errorln(e.Error())
	}

	w.Close()

	tm.Println(tm.Color(tm.Bold("File has been pushed to hdfs"), tm.GREEN))
	tm.Flush()

}
