package main

import (
	"io"
	"os"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cavaliercoder/grab"
)

func main() {
	const md5Original = "cd573cfaace07e7949bc0c46028904ff"
	url := "http://90.130.70.73/1GB.zip"
	fmt.Printf("Downloading %s...\n", url)
	resp, err := grab.Get(".", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", url, err)
		os.Exit(1)
	}

	fmt.Printf("File size %d \n", resp.Size)

	file, _ := os.Open(resp.Filename)
	defer file.Close()

	md5, err := getMD5Hash(file)
	println(md5)
	if err == nil && md5 == md5Original {
		println("Successfully Downloaded file ")
	} else {
		println("File corrupted")
	}

	os.Remove(resp.Filename)
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
