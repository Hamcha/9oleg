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
	OnAuth      func(AuthRequest) AuthResponse
	OnAttach    func(AttachRequest) AttachResponse
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
			fmt.Printf("(VERSION) MaxSize %d Version \"%s\"\n", ver.MaxSize, ver.Version)
		}
		err := write(con, makeMsg(Rversion, msg.Tag, ver))
		if err != nil {
			s.OnConnError(err)
			break
		}
	case AuthRequest:
		auth := data.(AuthRequest)
		if Debug {
			fmt.Printf("(AUTH) Afid %0#8x Uname \"%s\" Aname \"%s\"\n", auth.Afid, auth.Uname, auth.Aname)
		}
		if s.OnAuth != nil {
			resp := s.OnAuth(auth)
			err := write(con, makeMsg(Rauth, msg.Tag, resp))
			if err != nil {
				s.OnConnError(err)
				break
			}
		} else {
			err := sendErr(con, msg.Tag, "auth not required")
			if err != nil {
				s.OnConnError(err)
				break
			}
		}

	case AttachRequest:
		att := data.(AttachRequest)
		if Debug {
			fmt.Printf("(ATTACH) Fid %0#8x Afid %0#8x Uname \"%s\" Aname \"%s\"\n", att.Fid, att.Afid, att.Uname, att.Aname)
		}
		if s.OnAttach != nil {
			resp := s.OnAttach(att)
			err := write(con, makeMsg(Rattach, msg.Tag, resp))
			if err != nil {
				s.OnConnError(err)
				break
			}
		} else {
			// This will be fatal
			err := sendErr(con, msg.Tag, "not implemented")
			if err != nil {
				s.OnConnError(err)
				break
			}
		}
	case FlushRequest:
		flu := data.(FlushRequest)
		if Debug {
			fmt.Printf("(FLUSH) Tag %0#8x OldTag %0#8x\n", msg.Tag, flu.OldTag)
		}
		//TODO abort operation with specified tag
		err := write(con, makeMsg(Rflush, msg.Tag, nil))
		if err != nil {
			s.OnConnError(err)
			break
		}
	case UnknownData:
		if Debug {
			fmt.Printf("(UNKNOWN) Type %d Tag %0#8x Data %x\n", msg.Type, msg.Tag, data.(UnknownData).Raw)
		}
		err := sendErr(con, msg.Tag, "unknown command")
		if err != nil {
			s.OnConnError(err)
			break
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

func write(con net.Conn, data []byte) error {
	remaining := len(data)
	for remaining > 0 {
		n, err := con.Write(data)
		if err != nil {
			return err
		}
		remaining -= n
		data = data[n:]
	}
	return nil
}

func sendErr(conn net.Conn, msgTag uint16, msg string) error {
	return write(conn, makeMsg(Rerror, msgTag, ErrorData{msg}))
}
