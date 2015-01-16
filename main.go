package main

import (
	"github.com/google/go-github/github"
	"fmt"
)

func main()  {

}

func test() string {
	return "hello world"
}

func connect() {
	client := github.NewClient(nil)
	orgs, _, _ := client.Organizations.List("willnorris", nil)
	fmt.Print(orgs)
}
