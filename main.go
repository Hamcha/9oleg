package main

import (
	"fmt"
)

const listenaddr = "*"

func main() {
	vfs := makeFs("data", "oleg")

	fmt.Println("Listening on " + listenaddr)
	err := vfs.Listen(listenaddr)
	if err != nil {
		panic(err.Error())
	}
}
