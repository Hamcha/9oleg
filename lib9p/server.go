/*
   Virtual File System factory

   This class exposes all 9P methods to the user.
*/

package lib9p

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	OnConnError func(error) /* "On connection error" Handler */
}

func (s *Server) Listen(address string) error {
	listen := parseAddr(address)

	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.OnConnError(err)
			continue
		}

		go readClient(s, conn)
	}
}

func readClient(s *Server, con net.Conn) {
	b := bufio.NewReader(con)
	for {
		/* Read the total message length */
		bytes, err := b.Peek(4)
		if err != nil {
			if err.Error() != "EOF" {
				s.OnConnError(err)
			}
			break
		}
		length := uint32(dle(bytes))

		/* Read the whole message */
		remaining := length
		rawmsg := make([]byte, 0)
		for remaining > 0 {
			bmsg := make([]byte, remaining)
			n, err := b.Read(bmsg)
			if err != nil {
				s.OnConnError(err)
				break
			}
			rawmsg = append(rawmsg[:], bmsg[:]...)
			remaining -= uint32(n)
		}

		go handle(s, con, rawmsg)
	}
}

func handle(s *Server, con net.Conn, rawmsg []byte) {
	msg, data := parseMsg(rawmsg)
	switch data.(type) {
	case VersionData:
		ver := data.(VersionData)
		if Debug {
			fmt.Printf("(VERSION) MaxSize %d Version %s\n", ver.MaxSize, ver.Version)
		}
		write(s, con, makeMsg(Rversion, msg.Tag, ver))
	case UnknownData:
		if Debug {
			fmt.Printf("(UNKNOWN) Type %d Tag %x Data %x\n", msg.Type, msg.Tag, data.(UnknownData).Raw)
		}
	}
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

func write(s *Server, con net.Conn, data []byte) {
	remaining := len(data)
	for remaining > 0 {
		n, err := con.Write(data)
		if err != nil {
			if err.Error() != "EOF" {
				s.OnConnError(err)
			}
			break
		}
		remaining -= n
		data = data[n:]
	}
}
