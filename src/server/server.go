package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

type Dictionary map[string]string



func main() {
    parseDictionary("acronyms")
}


func parseDictionary(name string) (Dictionary) {
    file, err := os.Open("src/server/" + name)
    if err != nil { panic(err) }
    // close fi on exit and check for its returned error
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
        fmt.Println(words[0], words[1])
        d[words[0]] = words[1]
    }

    return d
}
