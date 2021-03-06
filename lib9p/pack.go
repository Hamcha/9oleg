/*
   Packing/Unpacking utils for sending stuff

   9P uses PASCALish strings, with the first 4 bytes indicating the length
   (each parameter type has a defined number of bytes to specify its length)
   NULL characters (0x00) are forbidden anyway.

   Also, all numbers (including length) have to be in Little Endian, so we need
   a shitton of stupid helper functions to correctly serialize them into bytes.
*/

package lib9p

func pack(buf []byte, maxlen uint) []byte {
	length := len(buf)
	var bytes []byte
	switch maxlen {
	case 1:
		bytes = le(uint8(length))
	case 2:
		bytes = le(uint16(length))
	case 4:
		bytes = le(uint32(length))
	case 8:
		bytes = le(uint64(length))
	}
	return append(bytes[:], buf[:]...)
}

func le(value interface{}) []byte {
	var bsize int
	var ivalue uint64
	switch value := value.(type) {
	case uint8, int8:
		bsize = 1
		ivalue = uint64(value.(uint8))
	case uint16, int16:
		bsize = 2
		ivalue = uint64(value.(uint16))
	case uint32, int32:
		bsize = 4
		ivalue = uint64(value.(uint32))
	case uint64, int64:
		bsize = 8
		ivalue = value.(uint64)
	}
	out := make([]byte, bsize)
	for i := range out {
		out[i] = uint8(ivalue)
		ivalue >>= 8
	}
	return out
}

func dle(value []byte) (out uint64) {
	for i := range value {
		out |= uint64(uint64(value[i]) << (8 * uint64(i)))
	}
	return
}

func dstr(value []byte) string {
	length := uint16(dle(value[0:2]))
	if length == 0 {
		return ""
	}
	return string(value[2 : 2+length])
}

func pstr(str string) []byte {
	return pack([]byte(str), 2)
}

func pqid(qid Qid) []byte {
	buf := make([]byte, 13)
	buf[0] = qid.Type
	copy(buf[1:5], le(qid.Version))
	copy(buf[5:13], le(qid.PathId))
	return buf
}

func pstat(stat Stat) []byte {
	nlen := len(stat.Name)
	ulen := len(stat.Uid)
	glen := len(stat.Gid)
	mlen := len(stat.Muid)
	length := 49 + nlen + ulen + glen + mlen
	sbytes := make([]byte, length)
	copy(sbytes[0:2], le(uint16(length-2)))
	copy(sbytes[2:4], le(stat.Type))
	copy(sbytes[4:8], le(stat.Dev))
	copy(sbytes[8:21], pqid(stat.Qid))
	copy(sbytes[21:25], le(stat.Mode))
	copy(sbytes[25:29], le(stat.Atime))
	copy(sbytes[29:33], le(stat.Mtime))
	copy(sbytes[33:41], le(stat.Length))
	s := 41
	e := s + nlen + 2
	copy(sbytes[s:e], pstr(stat.Name))
	s = e
	e = s + ulen + 2
	copy(sbytes[s:e], pstr(stat.Uid))
	s = e
	e = s + glen + 2
	copy(sbytes[s:e], pstr(stat.Gid))
	s = e
	e = s + mlen + 2
	copy(sbytes[s:e], pstr(stat.Muid))
	return sbytes
}
