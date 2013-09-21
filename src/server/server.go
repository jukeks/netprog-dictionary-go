package main

import (
	"bufio"
	"dictionary"
	"fmt"
	"net"
	"os"
	"protocol"
	"strings"
)

func main() {
	dict := parseDictionary("acronyms")

	listener, err := net.Listen("tcp", "0.0.0.0:6666")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleClient(conn, dict)
	}

}

func handleClient(conn net.Conn, dict dictionary.Dictionary) {
	defer conn.Close()

	msg, ok := protocol.ReadMessage(conn)
	if !ok {
		return
	}

	fmt.Println(msg)

	errorMsg := "7: Error"
	successMsg := "6: Success"
	var retMsg string

	switch msg[0] {
	case '1': // get
		key := msg[3:]
		value, ok := dict.Get(key)
		if ok {
			retMsg = "5: " + value
		} else {
			retMsg = errorMsg
		}

	case '2': // add
		split := strings.Split(msg[3:], "\t")
		key, value := split[0], split[1]
		ok := dict.Add(key, value)
		if ok {
			retMsg = successMsg
		} else {
			retMsg = errorMsg
		}

	case '3': // remove
		key := msg[3:]
		ok := dict.Remove(key)
		if ok {
			retMsg = successMsg
		} else {
			retMsg = errorMsg
		}

	case '4': // update
		split := strings.Split(msg[3:], "\t")
		key, value := split[0], split[1]
		ok := dict.Update(key, value)
		if ok {
			retMsg = successMsg
		} else {
			retMsg = errorMsg
		}

	case '8': // server get

	}

	protocol.WriteMessage(conn, retMsg)
}

func parseDictionary(name string) dictionary.Dictionary {
	file, err := os.Open("src/server/" + name)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	d := dictionary.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, "\t")
		d.Add(words[0], words[1])
	}

	return d
}
