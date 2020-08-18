//article: https://bionic.fullstory.com/why-you-should-be-using-errgroup-withcontext-in-golang-server-handlers/

package main

import (
    "context"
    "fmt"

    serial "server/serial"
    parallel "server/parallel"
    errgroup "server/errgroup"
)

const (
    USER = 1
)

func main() {
    friends, err := serial.GetFriends(context.Background(), USER)
    fmt.Println(err, friends)

    friends, err = parallel.GetFriends(context.Background(), USER)
    fmt.Println(err, friends)

    friends, err = errgroup.GetFriends(context.Background(), USER)
    fmt.Println(err, friends)
}
