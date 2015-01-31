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

/* I/O modes */
const (
	MRead   = 0x00
	MWrite  = 0x01
	MRdwr   = 0x02 // ReadWrite
	MExec   = 0x03
	MTrunc  = 0x10
	MRclose = 0x40
)

/* Stat definition and types */

type Stat struct {
	Type   uint16
	Dev    uint32
	Qid    Qid
	Mode   uint32
	Atime  uint32
	Mtime  uint32
	Length uint64
	Name   string
	Uid    string
	Gid    string
	Muid   string
}

const (
	DmDir    = 0x80000000
	DmAppend = 0x40000000
	DmExcl   = 0x20000000
	DmTmp    = 0x04000000
)

/* Messages */

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

type AuthResponse struct {
	Aqid Qid
}

type AttachRequest struct {
	Fid   uint32
	Afid  uint32
	Uname string
	Aname string
}

type AttachResponse struct {
	Qid Qid
}

type WalkRequest struct {
	Fid    uint32
	NewFid uint32
	Paths  []string
}

type WalkResponse struct {
	Qids []Qid
}

type ClunkRequest struct {
	Fid uint32
}

type OpenRequest struct {
	Fid  uint32
	Mode uint8
}

type OpenResponse struct {
	Qid    Qid
	IoUnit uint32
}

type CreateRequest struct {
	Fid        uint32
	Name       string
	Permission uint32
	Mode       uint8
}

type CreateResponse struct {
	Qid    Qid
	IoUnit uint32
}

type ReadRequest struct {
	Fid    uint32
	Offset uint64
	Count  uint32
}

type StatRequest struct {
	Fid uint32
}

type StatResponse struct {
	Stat Stat
}

type WstatRequest struct {
	Fid  uint32
	Stat Stat
}

type FlushRequest struct {
	OldTag uint32
}

type ErrorData struct {
	Message string
}

type UnknownData struct {
	Raw []byte
}
