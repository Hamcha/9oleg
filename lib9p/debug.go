package lib9p

import "fmt"

const (
	DebugReq   = true
	DebugSend  = true
	DebugBytes = true
	DebugANSI  = true
)

const (
	CSend  = "33"
	CRecv  = "32"
	CBytes = "34"
)

func debugWrt(msgType uint8, msgTag uint16, data interface{}) {
	switch msgType {
	case Rversion:
		ver := data.(VersionData)
		fmt.Printf(col(CSend, "R(VERSION) MaxSize %d Version \"%s\"\n"), ver.MaxSize, ver.Version)
	case Rauth:
		fmt.Printf(col(CSend, "R(AUTH) Aqid %0#16x\n"), data.(AuthResponse).Aqid)
	case Rattach:
		fmt.Printf(col(CSend, "R(ATTACH) Qid %v\n"), data.(AttachResponse).Qid)
	case Rwalk:
		fmt.Printf(col(CSend, "R(WALK) Qids %v\n"), data.(WalkResponse).Qids)
	case Ropen:
		open := data.(OpenResponse)
		fmt.Printf(col(CSend, "R(OPEN) Qid %v IoUnit %d\n"), open.Qid, open.IoUnit)
	case Rstat:
		fmt.Printf(col(CSend, "R(STAT) Stat %v\n"), data.(StatResponse).Stat)
	case Rerror:
		fmt.Printf(col(CSend, "R(ERROR) %s\n"), data.(ErrorData).Message)
	case Rread:
		fmt.Printf(col(CSend, "R(READ) - Data (%d bytes) -\n"), dle(data.([]byte)[0:4]))
	case Rclunk:
		fmt.Printf(col(CSend, "R(CLUNK)\n"))
	default:
		fmt.Printf(col(CSend, "R(UNKNOWN) %v"), data)
	}
}

func col(color string, str string) string {
	if !DebugANSI {
		return str
	}

	return "\033[" + color + "m" + str + "\033[0m"
}
