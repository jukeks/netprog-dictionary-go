package main

import (
    "fmt"
    "net"
    "protocol"
)

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

func handleClient(conn net.Conn, get, add, remove, update chan Request) {
    defer conn.Close()

    msg, ok := protocol.ReadMessage(conn)
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
            protocol.WriteMessage(conn, "5: " + reply.value)
        } else {
            protocol.WriteMessage(conn, "7: Error")
        }

    case '2': // add
    case '3': // remove
    case '4': // update
    case '8': // server get

    }
}
