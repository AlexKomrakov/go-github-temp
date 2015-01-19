package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	Error *log.Logger
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	filename := filepath.Join(dir, "/gohub.log")
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	f.Truncate(0)
	if err != nil {
		fmt.Print("error opening file: %v", err)
	}
	defer f.Close()

	loggerInit(f)
	Error.Println("Hello")

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)

	http.Handle("/", router)
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		Error.Println("Error on starting server: %v", err2)
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello world")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		Error.Println("error reading body")
	}
	Error.Println(string(body))
}

func loggerInit(f *os.File) {
	Error = log.New(f, "", log.Lshortfile)
	Error.Println("Logger start")
}

func test() string {
	return "hello world"
}

func connect() {
	client := github.NewClient(nil)
	orgs, _, _ := client.Organizations.List("willnorris", nil)
	fmt.Print(orgs)
}
