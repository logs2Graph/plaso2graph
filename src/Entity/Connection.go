package Entity

import (
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Timestamp int
	Date      time.Time

	SourceIP        string
	SourcePort      int
	DestinationIP   string
	DestinationPort int
	Protocol        string
	Intitiated      bool

	Computer   string
	User       string
	UserDomain string

	ProcessName string
	ProcessId   int
}

type Host struct {
	Domain string
	IP     string
}

func UnionConnections(dest []Connection, src []Connection) []Connection {
	for _, c := range src {
		dest = AddConnection(dest, c)
	}
	return dest
}

func AddConnection(cs []Connection, c Connection) []Connection {
	if c.SourceIP != "Not Found." && c.DestinationIP != "Not Found." {
		cs = append(cs, c)
	}
	return cs
}

func NewConnectionFromSysmon3(evtx EvtxLog) Connection {
	var c = *new(Connection)
	c.Computer = evtx.System.Computer
	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	c.Date = t
	c.Timestamp = int(t.UnixNano())

	c.SourceIP = GetDataValue(evtx, "SourceIp")
	tmp_str := GetDataValue(evtx, "SourcePort")
	if tmp_str != "Not Found." {
		tmp_int, err := strconv.ParseInt(tmp_str, 10, 64)
		handleErr(err)
		c.SourcePort = int(tmp_int)
	}

	c.DestinationIP = GetDataValue(evtx, "DestinationIp")
	tmp_str = GetDataValue(evtx, "DestinationPort")
	if tmp_str != "Not Found." {
		tmp_int, err := strconv.ParseInt(tmp_str, 10, 64)
		handleErr(err)
		c.DestinationPort = int(tmp_int)
	}

	c.Protocol = GetDataValue(evtx, "Protocol")
	c.Intitiated = GetDataValue(evtx, "Initiated") == "true"

	c.ProcessName = GetDataValue(evtx, "Image")
	tmp_str = GetDataValue(evtx, "ProcessId")
	if tmp_str != "Not Found." {
		tmp_int, err := strconv.ParseInt(tmp_str, 10, 64)
		handleErr(err)
		c.ProcessId = int(tmp_int)
	}

	tmp_str = GetDataValue(evtx, "User")
	if tmp_str != "Not Found." {
		splitted := strings.Split(tmp_str, "\\")
		if len(splitted) > 1 {
			c.UserDomain = splitted[0]
			c.User = splitted[1]
		} else {
			c.User = splitted[0]
		}
	}

	c.Computer = evtx.System.Computer
	return c
}
