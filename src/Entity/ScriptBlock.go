package Entity

import (
	"encoding/xml"
	//	"fmt"
	"strconv"
	"time"
)

type ScriptBlock struct {
	Timestamp     int
	Date          time.Time
	Computer      string
	ScriptBlockID string
	ProcessID     int
	Text          string
	MessageNumber int
	MessageTotal  int
	Path          string
	Context       string
	Evidence      string
}

func findScriptBlock(scriptBlocks []ScriptBlock, pwshScriptBlock ScriptBlock) int {
	for i, p := range scriptBlocks {
		if p.ScriptBlockID == pwshScriptBlock.ScriptBlockID {
			return i
		}
	}
	return -1
}

func AddScriptBlock(scriptBlocks []ScriptBlock, p ScriptBlock) []ScriptBlock {

	i := findScriptBlock(scriptBlocks, p)
	if i == -1 {
		scriptBlocks = append(scriptBlocks, p)
	}
	return scriptBlocks
}

func UnionScriptBlocks(dest []ScriptBlock, src []ScriptBlock) []ScriptBlock {
	for _, p := range src {
		dest = AddScriptBlock(dest, p)
	}
	return dest
}

func NewScriptBlockFrom4104(evtx EvtxLog) ScriptBlock {
	var s ScriptBlock

	s.Computer = evtx.System.Computer
	xmlByte, _ := xml.Marshal(evtx)
	s.Evidence = string(xmlByte)

	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	s.Date = t
	s.Timestamp = int(t.UnixMicro())
	i64, err := strconv.ParseInt(evtx.System.Execution.ProcessID, 0, 64)
	s.ProcessID = int(i64)

	if len(evtx.EventData.Data) == 0 {
		return s
	}
	intStr := GetDataValue(evtx, "MessageNumber")
	if intStr == "Not Found." {
		s.Text = GetDataValue(evtx, "Payload")
		s.Context = GetDataValue(evtx, "ContextInfo")
		return s
	}

	i64, err = strconv.ParseInt(intStr, 0, 64)
	handleErr(err)
	s.MessageNumber = int(i64)

	i64, err = strconv.ParseInt(GetDataValue(evtx, "MessageTotal"), 0, 64)
	handleErr(err)
	s.MessageTotal = int(i64)

	s.ScriptBlockID = GetDataValue(evtx, "ScriptBlockId")
	s.Text = GetDataValue(evtx, "ScriptBlockText")

	s.Path = GetDataValue(evtx, "Path")

	return s
}

func NewScriptBlockFrom4103(evtx EvtxLog) ScriptBlock {
	var s ScriptBlock

	s.Computer = evtx.System.Computer
	xmlByte, _ := xml.Marshal(evtx)
	s.Evidence = string(xmlByte)

	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	s.Date = t
	s.Timestamp = int(t.UnixMicro())
	i64, err := strconv.ParseInt(evtx.System.Execution.ProcessID, 0, 64)
	s.ProcessID = int(i64)

	if len(evtx.EventData.Data) == 0 {
		return s
	}

	s.Text = GetDataValue(evtx, "Payload")
	s.Context = GetDataValue(evtx, "ContextInfo")

	return s
}
