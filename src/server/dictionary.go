package main

import (
    "os"
    "bufio"
    "strings"
)

type Dictionary map[string]string

type Reply struct {
    value string
    ok bool
}

type Request struct {
    key string
    value string
    replyChan chan Reply
}


func manageDictionary(dictionary Dictionary, get, add, remove, update chan Request, quit chan bool) {
    done := false
    for !done {
        select {
        case req := <- get:
            value, ok := dictionary[req.key]
            req.replyChan <- Reply{value, ok}
        case req := <- add:
            _, present := dictionary[req.key]
            if !present {
                dictionary[req.key] = req.value
            }

            req.replyChan <- Reply{"", !present}
        case req := <- remove:
            _, present := dictionary[req.key]
            if present {
                delete(dictionary, req.key)
            }

            req.replyChan <- Reply{"", present}
        case req := <- update:
            _, present := dictionary[req.key]
            if present {
                dictionary[req.key] = req.value
            }

            req.replyChan <- Reply{"", present}
        case <- quit:
            done = true
        }
    }

    quit <- true
}

func parseDictionary(name string) (Dictionary) {
    file, err := os.Open("src/server/" + name)
    if err != nil { panic(err) }
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()

    d := make(Dictionary)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()      
        words := strings.Split(line, "\t")
        d[words[0]] = words[1]
    }

    return d
}