package main

import (
	"./lib9p"
	"fmt"
)

const listenaddr = "*"

func main() {
	vfs := new(lib9p.Server)
	vfs.OnConnError = func(err error) {
		fmt.Println(err.Error())
	}
	vfs.OnAttach = func(req lib9p.AttachRequest) (out lib9p.AttachResponse) {
		out.Qid = lib9p.Qid{
			Type:    lib9p.QtDir,
			Version: 1,
			PathId:  0,
		}
		return
	}

	fmt.Println("Listening on " + listenaddr)
	err := vfs.Listen(listenaddr)
	if err != nil {
		panic(err.Error())
	}
}
