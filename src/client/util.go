package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func generateHash(file *os.File) (string, error) {
	var MD5String string
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return MD5String, err
	}
	hashBytes := hash.Sum(nil)[:16]
	MD5String = hex.EncodeToString(hashBytes)
	return MD5String, nil
}

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
