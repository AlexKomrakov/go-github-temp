package main

import (
	"github.com/alexkomrakov/gohub/server"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/graceful"

	"time"
)

func main() {
	config, _ := server.GetServerConfig()
	router := server.Router()

	n := negroni.Classic() 
	n.UseHandler(router)
	graceful.Run(config.Adress, 30*time.Second, n) 
}
