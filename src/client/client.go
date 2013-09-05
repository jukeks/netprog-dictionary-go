package main

import (
	"flag"
	"fmt"
	"protocol"
	"net"
	"strconv"
)

type Server struct {
    ip string
    port int
}

func main() {
	ip := flag.String("i", "localhost", "server's ip address")
	port := flag.Int("p", 6666, "server's port number")
	get_arg := flag.Bool("get", false, "get value of key")
	add_arg := flag.Bool("add", false, "add key-value")
	update_arg := flag.Bool("update", false, "update key-value")
	remove_arg := flag.Bool("remove", false, "remove key")
	key_arg := flag.String("key", "", "key parameter for query")
	value_arg := flag.String("value", "", "value parameter for query")

	flag.Parse()

	server := Server{*ip, *port}

	var return_value string
	if *get_arg {
		return_value = get(server, *key_arg)
	} else if *add_arg {
		return_value = add(server, *key_arg, *value_arg)
	} else if *remove_arg {
		return_value = remove(server, *key_arg)
	} else if *update_arg {
		return_value = update(server, *key_arg, *value_arg)
	}

	fmt.Println(return_value)
}

func get(server Server, key string) (string) {
	return query(server, "1: " + key)
}

func add(server Server, key, value string) (string) {
	return query(server, "2: " + key + "\t" + value)
}

func remove(server Server, key string) (string) {
	return query(server, "3: " + key)
}

func update(server Server, key, value string) (string) {
	return query(server, "4: " + key + "\t" + value)
}

func query(server Server, message string) (string) {
	conn, err := net.Dial("tcp", server.ip + ":" + strconv.Itoa(server.port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	protocol.WriteMessage(conn, message)
	msg, ok := protocol.ReadMessage(conn)
	if !ok {
		panic("OMG")
	}

	return msg
}
