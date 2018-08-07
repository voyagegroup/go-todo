package controller

import (
	"net/http"

	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Video struct{}

const maxUploadSize = 2 * 1024 * 1024 * 1000000000000

func (t *Video) Upload(w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("maxUploadSize: %v, headerContentLength: %v", maxUploadSize, r.Header.Get("Content-Length"))

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		Error(w, fmt.Errorf("%s", err.Error()), http.StatusRequestEntityTooLarge)
		return nil
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	filetype := http.DetectContentType(fileBytes)
	if filetype != "video/mp4" {
		return fmt.Errorf("%s is not supported", filetype)
	}

	filepath := "public/upload/" + time.Now().Format("20060102_150405") + "_" + handler.Filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)

	return JSON(w, http.StatusCreated, filepath)
}
