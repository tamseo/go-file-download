package main

import (
	"io"
	"os"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cavaliercoder/grab"
	"time"
)

func main() {
	const md5Original = "cd573cfaace07e7949bc0c46028904ff"
	url := "http://90.130.70.73/1GB.zip"
	fmt.Printf("Downloading %s...\n", url)

	resp := downLoad(url);
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

func downLoad(url string) *grab.Response {
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)
	return resp
}
