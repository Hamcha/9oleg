package main

import (
	"./goleg"
	"./lib9p"
	"errors"
	"fmt"
	"net"
)

type FidData struct {
	Qid  lib9p.Qid
	Path []string
}

type Client struct {
	Fids map[uint32]FidData
}

type OlegFs struct {
	vfs     *lib9p.Server
	db      goleg.Database
	clients map[net.Conn]Client
}

func makeFs(dbdir string, dbname string) *OlegFs {
	/* Make OlegFs instance */
	ofs := new(OlegFs)
	ofs.clients = make(map[net.Conn]Client)

	/* Open OlegDB database */
	var err error
	ofs.db, err = goleg.Open(dbdir, dbname, goleg.F_APPENDONLY|goleg.F_LZ4|goleg.F_SPLAYTREE|goleg.F_AOL_FFLUSH)
	if err != nil {
		panic(err.Error())
	}

	/* Make VFS */
	vfs := new(lib9p.Server)
	vfs.OnConnError = ofs.ConnError
	vfs.OnAttach = ofs.Attach
	vfs.OnWalk = ofs.Walk
	vfs.OnClunk = ofs.Clunk
	vfs.OnOpen = ofs.Open
	vfs.OnRead = ofs.Read
	vfs.OnStat = ofs.Stat

	ofs.vfs = vfs
	return ofs
}

func (ofs *OlegFs) ConnError(con net.Conn, err error) {
	fmt.Println(err.Error())
}

func (ofs *OlegFs) Attach(con net.Conn, req lib9p.AttachRequest) (out lib9p.AttachResponse, err error) {
	path := make([]string, 0)

	out.Qid, _ = ofs.getQid(path)

	ofs.clients[con] = Client{
		Fids: make(map[uint32]FidData),
	}

	ofs.clients[con].Fids[req.Fid] = FidData{
		Qid:  out.Qid,
		Path: path,
	}
	return
}

func (ofs *OlegFs) Walk(con net.Conn, req lib9p.WalkRequest) (out lib9p.WalkResponse, err error) {
	client, ok := ofs.clients[con]
	if !ok {
		err = errors.New(lib9p.ErrDenied)
		return
	}
	if _, ok = client.Fids[req.Fid]; !ok {
		err = errors.New(lib9p.ErrUnknownFid)
		return
	}

	current := FidData{
		Qid:  client.Fids[req.Fid].Qid,
		Path: client.Fids[req.Fid].Path[:],
	}

	out.Qids = make([]lib9p.Qid, len(req.Paths))
	for i := range out.Qids {
		fmt.Printf("Walking to %s..\n", req.Paths[i])

		switch req.Paths[i] {
		case ".":
			out.Qids[i] = current.Qid
		case "..":
			nlen := len(current.Path) - 1
			if nlen < 0 {
				nlen = 0
			}
			current.Qid, err = ofs.getQid(current.Path[:nlen])
			out.Qids[i] = current.Qid
		default:
			current.Path = append(current.Path, req.Paths[i])
			out.Qids[i], err = ofs.getQid(current.Path)
			if err != nil {
				return
			}
		}
	}

	client.Fids[req.NewFid] = current
	return
}

func (ofs *OlegFs) Clunk(con net.Conn, req lib9p.ClunkRequest) error {
	client, ok := ofs.clients[con]
	if !ok {
		return errors.New(lib9p.ErrDenied)
	}

	if _, ok = client.Fids[req.Fid]; !ok {
		return errors.New(lib9p.ErrUnknownFid)
	}
	delete(client.Fids, req.Fid)
	return nil
}

func (ofs *OlegFs) Open(con net.Conn, req lib9p.OpenRequest) (out lib9p.OpenResponse, err error) {
	out.Qid = lib9p.Qid{
		Type:    lib9p.QtDir,
		Version: 1,
		PathId:  0,
	}
	out.IoUnit = 4096
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

func (ofs *OlegFs) getQid(path []string) (lib9p.Qid, error) {
	var qid lib9p.Qid
	var err error
	if len(path) < 1 {
		qid = lib9p.Qid{
			Type:    lib9p.QtDir,
			Version: 1,
			PathId:  0,
		}
	} else {
		err = errors.New(lib9p.ErrNotFound)
	}
	return qid, err
}
