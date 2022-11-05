package Entity

import (
	"encoding/xml"
	"time"
	//"fmt"
	"log"
	. "plaso2graph/master/src"
	"strconv"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Process struct {
	CreatedTime time.Time
	Timestamp   int
	Name        string
	PID         int
	Commandline string

	PPID                 int
	Pprocess_name        string
	Pprocess_commandline string

	User        string
	User_Domain string
	LogonID     int
	Computer    string
	Evidence    []string
}

func containsProcess(ps []Process, p Process) bool {
	for _, v := range ps {
		if v.PID == p.PID && v.PPID == p.PPID && v.Name == p.Name {
			return true
		}
	}
	return false
}

func GetProcesses(data []PlasoLog) []Process {

	var processes []Process

	for _, d := range data {

		if d.EvtxLog != nil {
			if d.EvtxLog.System.EventID == 4688 {
				process := NewProcessFrom4688(*d.EvtxLog)
				l := len(processes)

				// For some reason, plaso may dumplicates evtx entries. Here we filter the duplicates with a quick filter
				if l != 0 && process.PID == processes[l-1].PID && process.PPID == processes[l-1].PPID && process.Name == processes[l-1].Name {
					continue
				}

				processes = append(processes, process)
			}
			/*
				if d.EvtxLog.System.EventID == 1 && strings.Contains(d.EvtxLog.System.Provider.Name, "Sysmon") {
					process := NewProcessFromSysmon(*d.EvtxLog)
					processes = append(processes, process)
					continue
				}
			*/
		}
	}
	return processes
}

func convertOct(s string) int {
	i64, err := strconv.ParseInt(s, 0, 64)
	handleErr(err)
	return int(i64)
}

func NewProcessFrom4688(evtx EvtxLog) Process {
	var process Process
	process.Computer = evtx.System.Computer
	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	process.CreatedTime = t
	process.Timestamp = int(t.UnixNano())
	process.Name = GetDataValue(evtx, "NewProcessName")
	process.PID = convertOct(GetDataValue(evtx, "NewProcessId"))
	process.Commandline = GetDataValue(evtx, "CommandLine")

	process.PPID = convertOct(GetDataValue(evtx, "ProcessId"))
	process.Pprocess_name = GetDataValue(evtx, "ParentProcessName")

	process.User = GetDataValue(evtx, "TargetUserName")
	process.User_Domain = GetDataValue(evtx, "TargetDomainName")
	process.LogonID = convertOct(GetDataValue(evtx, "TargetLogonId"))

	xml_string, err := xml.Marshal(evtx)
	handleErr(err)
	process.Evidence = append(process.Evidence, string(xml_string))

	return process
}

func NewProcessFromSysmon(evtx EvtxLog) Process {
	var process Process
	//TODO: Process From Sysmon 1
	return process
}

func NewProcessFromPrefetch(pf PlasoLog) Process {
	var process Process
	//TODO: Process From Prefetch
	return process
}

func NewProcessFromSRUM(EvtxLog) Process {
	var process Process
	//TODO: ProcessFromSRUM
	return process
}
