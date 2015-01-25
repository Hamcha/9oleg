/*
   Data types and enums
*/

package lib9p

/* 9P settings */
const (
	Version     = "9P2000"
	DefaultPort = 564
	Debug       = true
)

/* Fcall errors */
const (
	NoTag = 0xffff
	NoFid = 0xffffffff
	NoUid = 0xffffffff
)

/* Fcall types */
const (
	Topenfd  = 98
	Ropenfd  = 99
	Tversion = 100
	Rversion = 101
	Tauth    = 102
	Rauth    = 103
	Tattach  = 104
	Rattach  = 105
	Terror   = 106 /* wait, what? */
	Rerror   = 107
	Tflush   = 108
	Rflush   = 109
	Twalk    = 110
	Rwalk    = 111
	Topen    = 112
	Ropen    = 113
	Tcreate  = 114
	Rcreate  = 115
	Tread    = 116
	Rread    = 117
	Twrite   = 118
	Rwrite   = 119
	Tclunk   = 120
	Rclunk   = 121
	Tremove  = 122
	Rremove  = 123
	Tstat    = 124
	Rstat    = 125
	Twstat   = 126
	Rwstat   = 127
)

/* Qid definition and types */

type Qid struct {
	Type    uint8
	Version uint32
	PathId  uint64
}

const (
	QtFile    = 0x00
	QtLink    = 0x01 // 9P2000.u
	QtSymlink = 0x02 // 9P2000.u
	QtTmp     = 0x04
	QtAuth    = 0x08 // for Tauth/Rauth
	QtMount   = 0x10
	QtExcl    = 0x20
	QtAppend  = 0x40
	QtDir     = 0x80
)
