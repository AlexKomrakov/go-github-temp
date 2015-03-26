package service
import (
    "log"
    "os"
    "github.com/codegangsta/negroni"
)

func GetFileLogger (filename string) *log.Logger {
    f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
    if err != nil {
        panic(err)
    }

    return log.New(f, "[gohub] ", log.Ldate | log.Ltime)
}

func GetRecoveryLogger (filename string) *negroni.Recovery {
    logger := GetFileLogger(filename)
    return &negroni.Recovery{
        Logger:     logger,
        PrintStack: true,
        StackAll:   false,
        StackSize:  1024 * 8,
    }
}
