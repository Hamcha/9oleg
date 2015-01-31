package main

import (
	"./lib9p"
	"errors"
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
	ofs.clients = make(map[net.Conn]Client)
	vfs := new(lib9p.Server)
	vfs.OnConnError = ofs.ConnError
	vfs.OnAttach = ofs.Attach
	vfs.OnWalk = ofs.Walk
	vfs.OnClunk = ofs.Clunk
	vfs.OnOpen = ofs.Open
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
	ofs.clients[con] = Client{
		Fids: make(map[uint32]string),
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

func (ofs *OlegFs) Clunk(con net.Conn, req lib9p.ClunkRequest) error {
	client, ok := ofs.clients[con]
	if !ok {
		return errors.New("Client not found - Please attach first")
	}

	if _, ok = client.Fids[req.Fid]; !ok {
		return errors.New("inexistand fid")
	}
	return nil
}

func (ofs *OlegFs) Open(con net.Conn, req lib9p.OpenRequest) (out lib9p.OpenResponse, err error) {

}
