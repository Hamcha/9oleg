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
