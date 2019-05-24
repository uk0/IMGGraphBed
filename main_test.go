package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	filePath:="/Users/zhangjianxin/Pictures/WX20190324-005914@2x.png"
	fileWriter, _ := bodyWriter.CreateFormFile("files", filePath)

	file, _ := os.Open(filePath)
	defer file.Close()

	io.Copy(fileWriter, file)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post("http://localhost:3000/upload/image", contentType, bodyBuffer)
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.Status)
	log.Println(string(resp_body))
}
