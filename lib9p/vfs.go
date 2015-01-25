/*
   Virtual File System factory

   This class exposes all 9P methods to the user.
*/

package lib9p

import (
	"bufio"
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
		go handle(v, conn)
	}
}

func handle(v *Vfs, con net.Conn) {
	b := bufio.NewReader(con)
	for {
		/* Read the total message length */
		bytes, err := b.Peek(4)
		if err != nil {
			v.OnConnError(err)
			break
		}
		length := dle(bytes)

		/* Read the whole message */
		remaining := length
		rawmsg := make([]byte, 0)
		for remaining > 0 {
			bmsg := make([]byte, length)
			n, err := b.Read(bmsg)
			if err != nil {
				v.OnConnError(err)
				break
			}
			rawmsg = append(rawmsg[:], bmsg[:]...)
			remaining -= uint64(n)
		}

		go handleMsg(v, rawmsg)
	}
}

func handleMsg(v *Vfs, msg []byte) {

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
