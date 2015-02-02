package main

import (
	"bytes"
	"strings"
	//"code.google.com/p/goauth2/oauth"
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

var (
	Error        *log.Logger
	github_token string
)

func start() {
	config := readConfig()
	github_token = config["github_token"]

	file := loggerInit(config["log_file"])
	defer file.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/github", githubHandler)

	http.Handle("/", router)
	err := http.ListenAndServe(config["server"], nil)
	if err != nil {
		Error.Println("Error on starting server: %v", err)
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

func loggerInit(filename string) (file *os.File) {
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		fmt.Print("Error on parsing abs path for logger: %v", err1)
	}

	filename = filepath.Join(dir, filename)

	f, err2 := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		fmt.Print("Error opening file: %v", err2)
	}
	err3 := f.Truncate(0)
	if err3 != nil {
		fmt.Print("Error on clearing log file: %v", err3)
	}

	Error = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	Error.Println("Logger start")

	return f
}
