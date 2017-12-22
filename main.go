package main

import (
	"net/http"
	"io"
	"os"
	"crypto/md5"
	"encoding/hex"
)

func main() {
	const md5Original = "2f282b84e7e608d5852449ed940bfc51"
	fileName := "100MB.zip"
	out, _ := os.Create(fileName)
	defer out.Close()
	resp, _ := http.Get("http://90.130.70.73/100MB.zip")
	defer resp.Body.Close()
	println("Downloading file ...")
	io.Copy(out, resp.Body)

	md5, err := getMD5Hash(out)
	println(md5)
	if err == nil && md5 == md5Original {
		println("Successfully Downloaded file ")
	} else {
		println("File corrupted")
	}
	os.Remove(fileName)
}

func getMD5Hash(file *os.File) (string, error) {

	var md5String string
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return md5String, err
	}
	hashInBytes := hash.Sum(nil)
	md5String = hex.EncodeToString(hashInBytes)
	return md5String, nil

}
