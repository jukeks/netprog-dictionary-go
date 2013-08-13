package protocol

import (
    "fmt"
    "strings"
    "net"
)

func ReadMessage(conn net.Conn) (string, bool) {
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

func WriteMessage(conn net.Conn, msg string) (bool) {
	fmt.Fprintf(conn, "%s\r\n", msg)
	return true
}