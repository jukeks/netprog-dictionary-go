package main

import (
    "flag"
    "fmt"
    "protocol"
    "net"
    "strconv"
)

func main() {
	ip := flag.String("i", "localhost", "server's ip address")
	port := flag.Int("p", 6666, "server's port number")
	get_key := flag.String("get", "LOL", "get value of key")
    flag.Parse()

    fmt.Println(get(*ip, *port, *get_key))
}

func get(ip string, port int, key string) (string) {
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
    if err != nil {
    	panic(err)
    }
    defer conn.Close()

    protocol.WriteMessage(conn, "1: " + key)
    msg, ok := protocol.ReadMessage(conn)
    if !ok {
    	panic("OMG")
    }

    return msg
}
