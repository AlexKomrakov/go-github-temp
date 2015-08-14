package main

import (
    "github.com/codegangsta/negroni"
    "github.com/stretchr/graceful"
    "time"
    "log"
    "net/http"
    "github.com/alexkomrakov/gohub/service"
    "github.com/alexkomrakov/gohub/server"
    "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
)

func main() {
    config := service.GetServerConfig()
    router := server.Router()

    n := negroni.New(service.GetRecoveryLogger(config.Logs.Error), negroni.NewLogger(), negroni.NewStatic(http.Dir("public")))

    store := cookiestore.New([]byte(config.SessionSecretKey))
    n.Use(sessions.Sessions("SESSION", store))

    n.UseHandler(router)

    log.Println("Starting server on address " + config.Adress)
    graceful.Run(config.Adress, 30*time.Second, n)
}
