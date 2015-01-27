package main

import (
	"./lib9p"
	"fmt"
	"net"
)

type Client struct {
	Fids map[uint32]string
}

type OlegFs struct {
	clients map[net.Conn]Client
}

func makeFs(dbdir string, dbname string) *lib9p.Server {
	ofs := new(OlegFs)
	vfs := new(lib9p.Server)
	vfs.OnConnError = ofs.ConnError
	vfs.OnAttach = ofs.Attach
	vfs.OnWalk = ofs.Walk
	vfs.OnClunk = ofs.Clunk
	return vfs
}

func (ofs *OlegFs) ConnError(con net.Conn, err error) {
	fmt.Println(err.Error())
}

func (ofs *OlegFs) Attach(con net.Conn, req lib9p.AttachRequest) (out lib9p.AttachResponse, err error) {
	out.Qid = lib9p.Qid{
		Type:    lib9p.QtDir,
		Version: 1,
		PathId:  0,
	}
	return
}

func (ofs *OlegFs) Walk(con net.Conn, req lib9p.WalkRequest) (out lib9p.WalkResponse, err error) {
	out.NoQids = req.NoPaths
	out.Qids = make([]lib9p.Qid, req.NoPaths)
	for i := range out.Qids {
		out.Qids[i] = lib9p.Qid{
			Type:    lib9p.QtDir,
			Version: 1,
			PathId:  0,
		}
	}
	return
}

func (ofs *OlegFs) Clunk(con net.Conn, req lib9p.ClunkRequest) (err error) {
	return
}
