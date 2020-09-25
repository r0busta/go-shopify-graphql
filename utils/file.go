package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		panic(err)
	}
}

func ReadFile(file string) (data string, err error) {
	var bytes []byte
	bytes, err = ioutil.ReadFile(file)
	data = string(bytes)
	return
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer CloseFile(out)

	_, err = io.Copy(out, resp.Body)
	return err
}
