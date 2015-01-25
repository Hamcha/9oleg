/*
   Virtual File System factory

   This class exposes all 9P methods to the user.
*/

package lib9p

import (
	"net"
	"strconv"
	"strings"
)

type Vfs struct {
	OnConnError func(error) /* "On connection error" Handler */
}

func (v *Vfs) Listen(address string) error {
	listen := parseAddr(address)
	ln, err := net.Listen("udp", listen)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			v.OnConnError(err)
		}
		go handle(conn)
	}
}

func handle(net.Conn) {
	//TODO
}

func parseAddr(addr string) string {
	if addr[0] == '*' {
		addr = addr[1:]
	}
	if strings.Index(addr, ":") < 0 {
		addr = addr + ":" + strconv.Itoa(DefaultPort)
	}
	return addr
}
