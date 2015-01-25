package lib9p

type MessageInfo struct {
	Length uint32
	Type   uint8
	Tag    uint16
}

type VersionData struct {
	MaxSize uint32
	Version string
}

type AuthRequest struct {
	Afid  uint32
	Uname string
	Aname string
}

type AttachRequest struct {
	Fid   uint32
	Afid  uint32
	Uname string
	Aname string
}

type ErrorData struct {
	Message string
}

type UnknownData struct {
	Raw []byte
}

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
		bytes = append(le(uint32(ver.MaxSize))[:], pstr(ver.Version)[:]...)
	case ErrorData:
		bytes = pstr(data.(ErrorData).Message)
	case UnknownData:
		bytes = data.(UnknownData).Raw
	}
	length := uint32(7 + len(bytes)) /* Length(4) + Tag(2) + Type(1) = 7 */
	msg := make([]byte, length)
	copy(msg, le(length))
	msg[4] = msgType
	copy(msg[5:7], le(msgTag))
	copy(msg[7:], bytes)
	return msg
}
