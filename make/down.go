package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// MakeDown download static resources
func MakeDown() {
	get := func(url string) ([]byte, error) {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("response status is " + strconv.Itoa(res.StatusCode))
		}
		return ioutil.ReadAll(res.Body)
	}

	for _, url := range staticFiles {
		name := last(strings.Split(url, "/"))
		filename := staticDir + "/" + name
		exist := fileExists(filename)
		if exist {
			println(filename, "already exist")
			continue
		}
		data, err := get(url)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(filename, data, 0664)
		if err != nil {
			panic(err)
		}
		println(filename, "download successful")
	}
}
