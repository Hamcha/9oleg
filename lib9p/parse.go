package lib9p

func parseMsg(b []byte) (msg MessageInfo, data interface{}) {
	msg.Length = uint32(dle(b[0:4]))
	msg.Type = uint8(b[4])
	msg.Tag = uint16(dle(b[5:7]))

	switch msg.Type {
	case Tversion, Rversion:
		data = VersionData{
			MaxSize: uint32(dle(b[7:11])),
			Version: dstr(b[11:]),
		}
	case Tauth:
		uname := dstr(b[11:])
		aname := dstr(b[13+len(uname):])
		data = AuthRequest{
			Afid:  uint32(dle(b[7:11])),
			Uname: uname,
			Aname: aname,
		}
	case Tattach:
		uname := dstr(b[15:])
		aname := dstr(b[17+len(uname):])
		data = AttachRequest{
			Fid:   uint32(dle(b[7:11])),
			Afid:  uint32(dle(b[11:15])),
			Uname: uname,
			Aname: aname,
		}
	case Twalk:
		nopaths := uint16(dle(b[15:17]))
		paths := make([]string, nopaths)
		offset := 15
		for i := range paths {
			paths[i] = dstr(b[offset:])
			offset += 2 + len(paths[i])
		}
		data = WalkRequest{
			Fid:    uint32(dle(b[7:11])),
			NewFid: uint32(dle(b[11:15])),
			Paths:  paths,
		}
	case Tclunk:
		data = ClunkRequest{
			Fid: uint32(dle(b[7:11])),
		}
	case Topen:
		data = OpenRequest{
			Fid:  uint32(dle(b[7:11])),
			Mode: uint8(b[11]),
		}
	case Tread:
		data = ReadRequest{
			Fid:    uint32(dle(b[7:11])),
			Offset: uint64(dle(b[11:19])),
			Count:  uint32(dle(b[19:23])),
		}
	case Tflush:
		data = FlushRequest{
			OldTag: uint32(dle(b[7:11])),
		}
	default:
		data = UnknownData{
			Raw: b[7:],
		}
	}
	return
}

func makeMsg(msgType uint8, msgTag uint16, data interface{}) []byte {
	var bytes []byte
	switch data.(type) {
	case VersionData:
		ver := data.(VersionData)
		bytes = append(le(ver.MaxSize)[:], pstr(ver.Version)[:]...)
	case AuthResponse:
		bytes = pqid(data.(AuthResponse).Aqid)
	case AttachResponse:
		bytes = pqid(data.(AttachResponse).Qid)
	case WalkResponse:
		walk := data.(WalkResponse)
		bytes = le(uint16(len(walk.Qids)))
		for _, x := range walk.Qids {
			bytes = append(bytes[:], pqid(x)[:]...)
		}
	case OpenResponse:
		open := data.(OpenResponse)
		bytes = pqid(open.Qid)
		bytes = append(bytes[:], le(open.IoUnit)[:]...)
	case ErrorData:
		bytes = pstr(data.(ErrorData).Message)
	case UnknownData:
		bytes = data.(UnknownData).Raw
	case []byte:
		bytes = data.([]byte)
	case nil:
		bytes = make([]byte, 0)
	}
	length := uint32(7 + len(bytes)) /* Length(4) + Tag(2) + Type(1) = 7 */
	msg := make([]byte, length)
	copy(msg, le(length))
	msg[4] = msgType
	copy(msg[5:7], le(msgTag))
	copy(msg[7:], bytes)
	return msg
}
