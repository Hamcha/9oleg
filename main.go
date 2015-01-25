package main

import (
	"./lib9p"
	"fmt"
)

const listenaddr = "*"

func main() {
	vfs := new(lib9p.Vfs)
	vfs.OnConnError = func(err error) {
		fmt.Println(err.Error())
	}
	fmt.Println("Listening on " + listenaddr)
	err := vfs.Listen(listenaddr)
	if err != nil {
		panic(err.Error())
	}
}
