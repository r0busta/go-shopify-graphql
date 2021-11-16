package utils

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func CloseFile(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func ReadFile(file string) (string, error) {
	var bytes []byte
	bytes, err := ioutil.ReadFile(file)
	data := string(bytes)
	return data, err
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
