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
	OnConnError func(net.Conn, error) /* "On connection error" Handler */
	OnAuth      func(net.Conn, AuthRequest) (AuthResponse, error)
	OnAttach    func(net.Conn, AttachRequest) (AttachResponse, error)
	OnWalk      func(net.Conn, WalkRequest) (WalkResponse, error)
	OnOpen      func(net.Conn, OpenRequest) (OpenResponse, error)
	OnRead      func(net.Conn, ReadRequest) ([]byte, error)
	OnStat      func(net.Conn, StatRequest) (StatResponse, error)
	OnClunk     func(net.Conn, ClunkRequest) error
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
			s.OnConnError(conn, err)
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
				s.OnConnError(con, err)
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
				s.OnConnError(con, err)
				break
			}
			rawmsg = append(rawmsg[:], bmsg[:]...)
			remaining -= uint32(n)
		}

		go handle(s, con, rawmsg)
	}
}

func handle(s *Server, con net.Conn, rawmsg []byte) {
	if DebugBytes {
		fmt.Printf(col(CBytes, "\nRECV > %0#x\n"), rawmsg)
	}
	msg, data := parseMsg(rawmsg)
	switch data.(type) {
	case VersionData:
		ver := data.(VersionData)
		if DebugReq {
			fmt.Printf(col(CRecv, "(VERSION) MaxSize %d Version \"%s\"\n"), ver.MaxSize, ver.Version)
		}
		err := write(con, makeMsg(Rversion, msg.Tag, ver))
		if err != nil {
			s.OnConnError(con, err)
			break
		}

	case AuthRequest:
		auth := data.(AuthRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(AUTH) Afid %0#8x Uname \"%s\" Aname \"%s\"\n"), auth.Afid, auth.Uname, auth.Aname)
		}
		if s.OnAuth != nil {
			resp, err := s.OnAuth(con, auth)
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rauth, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
		} else {
			err := sendErr(con, msg.Tag, "auth not required")
			if err != nil {
				s.OnConnError(con, err)
			}
		}

	case AttachRequest:
		att := data.(AttachRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(ATTACH) Fid %0#8x Afid %0#8x Uname \"%s\" Aname \"%s\"\n"), att.Fid, att.Afid, att.Uname, att.Aname)
		}
		if s.OnAttach != nil {
			resp, err := s.OnAttach(con, att)
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rattach, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case WalkRequest:
		walk := data.(WalkRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(WALK) Fid %0#8x NewFid %0#8x Paths %v\n"), walk.Fid, walk.NewFid, walk.Paths)
		}
		if s.OnWalk != nil {
			resp, err := s.OnWalk(con, walk)
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rwalk, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case ClunkRequest:
		if DebugReq {
			fmt.Printf(col(CRecv, "(CLUNK) Fid %0#8x\n"), data.(ClunkRequest).Fid)
		}
		if s.OnClunk != nil {
			err := s.OnClunk(con, data.(ClunkRequest))
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rclunk, msg.Tag, nil))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case OpenRequest:
		open := data.(OpenRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(OPEN) Fid %0#8x Mode %0#2x\n"), open.Fid, open.Mode)
		}
		if s.OnOpen != nil {
			resp, err := s.OnOpen(con, open)
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Ropen, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case ReadRequest:
		read := data.(ReadRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(READ) Fid %0#8x Offset %0#16x Count %0#8x\n"), read.Fid, read.Offset, read.Count)
		}
		if s.OnRead != nil {
			resp, err := s.OnRead(con, read)
			resp = append(le(uint32(len(resp)))[:], resp[:]...)
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rread, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case StatRequest:
		if DebugReq {
			fmt.Printf(col(CRecv, "(STAT) Fid %0#8x\n"), data.(StatRequest).Fid)
		}
		if s.OnStat != nil {
			resp, err := s.OnStat(con, data.(StatRequest))
			if err != nil {
				sendErr(con, msg.Tag, err.Error())
				break
			}
			err = write(con, makeMsg(Rstat, msg.Tag, resp))
			if err != nil {
				s.OnConnError(con, err)
			}
			break
		}
		sendErr(con, msg.Tag, "not implemented")

	case FlushRequest:
		flu := data.(FlushRequest)
		if DebugReq {
			fmt.Printf(col(CRecv, "(FLUSH) Tag %0#8x OldTag %0#8x\n"), msg.Tag, flu.OldTag)
		}
		//TODO abort operation with specified tag
		err := write(con, makeMsg(Rflush, msg.Tag, nil))
		if err != nil {
			s.OnConnError(con, err)
		}

	case UnknownData:
		if DebugReq {
			fmt.Printf(col(CRecv, "(UNKNOWN) Type %d Tag %0#8x Data %x\n"), msg.Type, msg.Tag, data.(UnknownData).Raw)
		}
		sendErr(con, msg.Tag, "unknown command")
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
	if DebugBytes {
		fmt.Printf(col(CBytes, "SEND < %0#x\n"), data)
	}
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
