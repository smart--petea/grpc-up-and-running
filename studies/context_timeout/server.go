package main

import (
    "context"
    "fmt"

    serial "server/serial"
)

const (
    USER = 1
)

func main() {
    friends, err := serial.GetFriends(context.Background(), USER)
    fmt.Println(friends)
    fmt.Println(err)
}
