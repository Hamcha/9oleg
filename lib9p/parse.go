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
	default:
		data = UnknownData{
			Raw: b[7:],
		}
	}
	return
}
