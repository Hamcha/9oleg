package main

import (
	"fmt"
)

const listenaddr = "*"

func main() {
	ofs := makeFs("data", "oleg")
	defer ofs.db.Close()

	fmt.Println("Listening on " + listenaddr)
	err := ofs.vfs.Listen(listenaddr)
	if err != nil {
		panic(err.Error())
	}
}
