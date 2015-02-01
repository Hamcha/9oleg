package main

import (
	"./lib9p"
	"errors"
	"fmt"
	"net"
)

type Client struct {
	Pwd  []string
	Fids map[uint32][]string
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
	vfs.OnRead = ofs.Read
	vfs.OnStat = ofs.Stat
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
		Pwd:  make([]string, 0),
		Fids: make(map[uint32][]string),
	}

	ofs.clients[con].Fids[req.Fid] = ofs.clients[con].Pwd
	return
}

func (ofs *OlegFs) Walk(con net.Conn, req lib9p.WalkRequest) (out lib9p.WalkResponse, err error) {
	client, ok := ofs.clients[con]
	if !ok {
		err = errors.New("Client not found - Please attach first")
		return
	}
	if _, ok = client.Fids[req.Fid]; !ok {
		err = errors.New("nonexistant fid")
		return
	}

	out.Qids = make([]lib9p.Qid, len(req.Paths))
	for i := range out.Qids {
		out.Qids[i] = lib9p.Qid{
			Type:    lib9p.QtDir,
			Version: 1,
			PathId:  0,
		}
		fmt.Printf("Walking to %s..\n", req.Paths[i])
	}

	client.Fids[req.NewFid] = client.Pwd
	return
}

func (ofs *OlegFs) Clunk(con net.Conn, req lib9p.ClunkRequest) error {
	client, ok := ofs.clients[con]
	if !ok {
		return errors.New("Client not found - Please attach first")
	}

	if _, ok = client.Fids[req.Fid]; !ok {
		return errors.New("nonexistant fid")
	}
	delete(client.Fids, req.Fid)
	return nil
}

func (ofs *OlegFs) Open(con net.Conn, req lib9p.OpenRequest) (out lib9p.OpenResponse, err error) {
	client, ok := ofs.clients[con]
	if !ok {
		err = errors.New("Client not found - Please attach first")
		return
	}

	out.Qid = lib9p.Qid{
		Type:    lib9p.QtDir,
		Version: 1,
		PathId:  0,
	}
	out.IoUnit = 2048

	client.Fids[req.Fid] = client.Pwd[:]
	return
}

func (ofs *OlegFs) Read(con net.Conn, req lib9p.ReadRequest) (b []byte, err error) {
	b = make([]byte, 0)
	return
}

func (ofs *OlegFs) Stat(con net.Conn, req lib9p.StatRequest) (out lib9p.StatResponse, err error) {
	out = lib9p.StatResponse{
		Stat: lib9p.Stat{
			Qid: lib9p.Qid{
				Type:    lib9p.QtDir,
				Version: 1,
				PathId:  0,
			},
			Mode:   lib9p.DmDir,
			Atime:  0,
			Mtime:  0,
			Length: 0,
			Name:   "/",
			Uid:    "none",
			Gid:    "none",
			Muid:   "none",
		},
	}
	return
}
