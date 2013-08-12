package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "net"
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


func main() {
    dict := parseDictionary("acronyms")

    // channels for using dictionary
    get := make(chan Request, 100)
    add := make(chan Request, 100)
    remove := make(chan Request, 100)
    update := make(chan Request, 100)
    quit := make(chan bool)

    defer func() { 
        // stopping dictionary manager
        quit <- true 
        <- quit
    }()

    go manageDictionary(dict, get, add, remove, update, quit)

    listener, err := net.Listen("tcp", "0.0.0.0:6666")
    if err != nil {
        panic(err)
    }
     
    for {
        conn, err := listener.Accept()
        if err != nil {
            panic(err)
        }

        go handleClient(conn, get, add, remove, update)
    }

}

func readMessage(conn net.Conn) (string, bool) {
    buffer := make([]byte, 4096)
    msg := ""

    for {
        bytesRead, err := conn.Read(buffer)
        if err != nil {
            panic(err)
            return "", false
        }

        msg += string(buffer[:bytesRead])
        if strings.Contains(msg, "\n") {
            break
        }
    }

    msg = strings.Replace(msg, "\r\n", "", 1)

    return msg, true
}

func handleClient(conn net.Conn, get, add, remove, update chan Request) {
    defer conn.Close()

    msg, ok := readMessage(conn)
    if !ok {
        return
    }

    fmt.Println(msg)

    replyChan := make(chan Reply)

    switch msg[0] {
    case '1': // get
        key := msg[3:]
        get <- Request{key, "", replyChan}
        reply := <- replyChan
        if reply.ok {
            fmt.Fprintf(conn, "5: %s\r\n", reply.value)
        } else {
            fmt.Fprintf(conn, "7: Error\r\n")
        }

    case '2': // add
    case '3': // remove
    case '4': // update
    case '8': // server get

    }
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
