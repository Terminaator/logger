package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	redirect string
)

func logToFile(req []byte, res []byte) {
	var buf []byte

	f, err := os.OpenFile("response.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	buf = append(buf, []byte("---------------START---------------\n")...)
	buf = append(buf, []byte("---------------REQUEST---------------\n")...)
	buf = append(buf, req...)
	buf = append(buf, []byte("\n---------------RESPONSE---------------\n")...)
	buf = append(buf, res...)
	buf = append(buf, []byte("----------------END-----------------\n")...)

	if _, err = f.Write(buf); err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func handle(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	bodyRequest, _ := ioutil.ReadAll(req.Body)

	var response *http.Response
	if method == "GET" {
		response, _ = http.Get(redirect + path)
		log.Println("eee")
	} else if method == "POST" {
		response, _ = http.Post(redirect+path, req.Header.Get("Content-type"), bytes.NewReader(bodyRequest))
	}

	bodyResponse, _ := ioutil.ReadAll(response.Body)

	res.Write(bodyResponse)

	go logToFile(bodyRequest, bodyResponse)

	log.Println("method:", method, "path:", path)
}

func server() {
	http.HandleFunc("/", handle)

	if err := http.ListenAndServe("127.0.0.1:7000", nil); err != nil {
		panic(err)
	}
}

func main() {
	redirect = getEnv("REDIRECT", "")

	server()
}
