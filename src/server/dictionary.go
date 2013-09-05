package main

import (
    "os"
    "bufio"
    "strings"
    "sync"
)

type Dictionary struct {
    l sync.Mutex
    m map[string]string
}

func NewDictionary() (Dictionary) {
    d := Dictionary{}
    d.m = make(map[string]string)

    return d
}


func (d *Dictionary) Add(key, value string) (bool) {
    d.l.Lock()
    _, ok := d.m[key]
    if !ok {
        d.m[key] = value
    }

    d.l.Unlock()

    return !ok
}

func (d *Dictionary) Update(key, value string) (bool) {
    d.l.Lock()
    _, ok := d.m[key]
    if ok {
        d.m[key] = value
    }

    d.l.Unlock()

    return ok
}

func (d *Dictionary) Get(key string) (string, bool) {
    d.l.Lock()
    value, ok := d.m[key]
    d.l.Unlock()

    return value, ok
}

func (d *Dictionary) Remove(key string) (bool) {
    d.l.Lock()
    _, ok := d.m[key]
    if ok {
        delete(d.m, key)
    }

    d.l.Unlock()

    return ok
}

func parseDictionary(name string) (Dictionary) {
    file, err := os.Open("src/server/" + name)
    if err != nil { panic(err) }
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()

    d := NewDictionary()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()      
        words := strings.Split(line, "\t")
        d.m[words[0]] = words[1]
    }

    return d
}
