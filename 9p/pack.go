/*
   Packing utils for sending stuff

   9P uses PASCALish strings, with the first 4 bytes indicating the length
   (each parameter type has a defined number of bytes to specify its length)
   NULL characters (0x00) are forbidden anyway.

   Also, all numbers (including lenght) have to be in Little Endian, so we need
   a shitton of stupid helper functions to correctly serialize them into bytes.
*/

package lib9p

import (
	"reflect"
)

func pack(buf []byte, maxlen uint) []byte {
	length := len(buf)
	var bytes []byte
	switch maxlen {
	case 1:
		bytes = le(uint8(length))
		break
	case 2:
		bytes = le(uint16(length))
		break
	case 4:
		bytes = le(uint32(length))
		break
	case 8:
		bytes = le(uint64(length))
		break
	}
	return append(bytes[:], buf[:]...)
}

func le(value interface{}) []byte {
	bsize := reflect.Size(value)
	out := make([]byte, bsize)
	for i := range out {
		out[i] = uint8(value)
		value <<= 8
	}
	return out
}
