package Entity

import (
	"encoding/xml"
	//"fmt"
	//"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	CreatedTime time.Time
	Timestamp   int
	FullPath    string
	Filename    string
	PID         int
	Commandline string

	PPID                     int
	ParentProcessName        string
	ParentProcessCommandline string

	User             string
	UserDomain       string
	ParentUser       string
	ParentUserDomain string
	LogonID          int
	Computer         string
	Sha256Hash       string
	Evidence         []string
}

func containsProcess(ps []Process, p Process) bool {
	for _, v := range ps {
		if v.PID == p.PID && v.PPID == p.PPID && v.FullPath == p.FullPath {
			return true
		}
	}
	return false
}

func AddProcess(ps []Process, p Process) []Process {
	if p.Filename != "" {
		ps = append(ps, p)
	}
	return ps
}

func UnionProcesses(dest []Process, src []Process) []Process {
	for _, p := range src {
		dest = AddProcess(dest, p)
	}
	return dest
}

func removeProcess(array []Process, index int) []Process {
	array[index] = array[len(array)-1]
	return array[:len(array)-1]
}

// Merge Last 2 * batch_size process
func MergeProcesses(processes []Process, approx int) []Process {

	for i := 0; i < len(processes); i++ {
		var markedToRemove []int
		for j := 0; j < len(processes); j++ {
			// We merge process if they have the same Filename and have a timestamp approximatly close
			if i != j && processes[i].Filename == processes[j].Filename && processes[j].Timestamp-approx < processes[i].Timestamp && processes[i].Timestamp < processes[j].Timestamp+approx {
				processes[i] = mergeProcess(processes[i], processes[j])
				// We mark the process that we have merged to be removed. (we don't mess with indexes in J's for loop)
				markedToRemove = append(markedToRemove, j)
			}
		}

		sort.Sort(sort.Reverse(sort.IntSlice(markedToRemove)))
		// Here we remove merged Process
		for _, index := range markedToRemove {
			processes = removeProcess(processes, index)
		}
	}
	return processes
}

func mergeProcess(dest Process, src Process) Process {
	if dest.Commandline == "" {
		dest.Commandline = src.Commandline
	}

	if dest.FullPath == "" {
		dest.Filename = src.FullPath
	}

	if dest.LogonID == 0 {
		dest.LogonID = src.LogonID
	}

	if dest.PID == 0 {
		dest.PID = src.PID
	}

	if dest.PPID == 0 {
		dest.PPID = src.PID
	}

	if dest.User == "" {
		dest.User = src.User
	}

	if dest.UserDomain == "" {
		dest.UserDomain = src.UserDomain
	}

	if dest.ParentProcessCommandline == "" {
		dest.ParentProcessCommandline = src.ParentProcessCommandline
	}

	if dest.ParentProcessName == "" {
		dest.ParentProcessName = src.ParentProcessName
	}

	dest.Evidence = append(dest.Evidence, src.Evidence...)

	return dest
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
	process.FullPath = GetDataValue(evtx, "NewProcessName")
	process.PID = convertOct(GetDataValue(evtx, "NewProcessId"))
	process.Commandline = GetDataValue(evtx, "CommandLine")
	splitted_path := strings.Split(process.FullPath, "\\")
	process.Filename = splitted_path[len(splitted_path)-1]

	process.PPID = convertOct(GetDataValue(evtx, "ProcessId"))
	process.ParentProcessName = GetDataValue(evtx, "ParentProcessName")

	process.User = GetDataValue(evtx, "TargetUserName")
	process.UserDomain = GetDataValue(evtx, "TargetDomainName")
	process.LogonID = convertOct(GetDataValue(evtx, "TargetLogonId"))

	xml_string, err := xml.Marshal(evtx)
	handleErr(err)
	process.Evidence = append(process.Evidence, string(xml_string))

	return process
}

func NewProcessFromSysmon1(evtx EvtxLog) Process {
	var process Process
	process.Computer = evtx.System.Computer
	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	process.CreatedTime = t
	process.Timestamp = int(t.UnixNano())

	process.FullPath = GetDataValue(evtx, "Image")
	splitted_path := strings.Split(process.FullPath, "\\")
	process.Filename = splitted_path[len(splitted_path)-1]
	process.Commandline = GetDataValue(evtx, "CommandLine")
	process.PID = convertOct(GetDataValue(evtx, "ProcessId"))

	//Parse Hash
	tmp := GetDataValue(evtx, "Hashes")
	hashes := strings.Split(tmp, ",")
	for _, hash := range hashes {
		if strings.Contains(hash, "SHA256") {
			process.Sha256Hash = strings.Split(hash, "=")[1]
		}
	}

	process.PPID = convertOct(GetDataValue(evtx, "ParentProcessId"))
	process.ParentProcessName = GetDataValue(evtx, "ParentImage")
	process.ParentProcessCommandline = GetDataValue(evtx, "ParentCommandLine")

	//Parse Domain and Users
	tmp = GetDataValue(evtx, "ParentUser")
	splitted_user := strings.Split(tmp, "\\")
	if len(splitted_user) > 1 {
		process.ParentUser = splitted_user[1]
		process.ParentUserDomain = splitted_user[0]
	} else {
		process.User = tmp
	}

	tmp = GetDataValue(evtx, "User")
	splitted_user = strings.Split(tmp, "\\")
	if len(splitted_user) > 1 {
		process.User = splitted_user[1]
		process.UserDomain = splitted_user[0]
	} else {
		process.User = tmp
	}

	xml_string, err := xml.Marshal(evtx)
	handleErr(err)
	process.Evidence = append(process.Evidence, string(xml_string))

	return process
}

func NewProcessFromPrefetchFile(pf PlasoLog) Process {
	var process = *new(Process)

	process.Evidence = append(process.Evidence, pf.DisplayName)

	prefetch_file := getFilename(pf.DisplayName)
	process.Filename = strings.Split(prefetch_file, "-")[0]

	var utc, _ = time.LoadLocation("UTC")
	process.Timestamp = int(pf.Timestamp)
	process.CreatedTime = time.UnixMicro(int64(pf.Timestamp)).In(utc)
	// {"__container_type__": "event", "__type__": "AttributeContainer", "data_type": "windows:volume:creation", "date_time": {"__class_name__": "Filetime", "__type__": "DateTimeValues", "timestamp": 132902282486494590}, "device_path": "\\VOLUME{01d829ebf972357e-10f97ebe}", "display_name": "OS:/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Windows/prefetch/CHROME.EXE-AED7BA3D.pf", "filename": "/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Windows/prefetch/CHROME.EXE-AED7BA3D.pf", "inode": "-", "message": "\\VOLUME{01d829ebf972357e-10f97ebe} Serial number: 0x10F97EBE Origin: CHROME.EXE-AED7BA3D.pf", "origin": "CHROME.EXE-AED7BA3D.pf", "parser": "prefetch", "pathspec": {"__type__": "PathSpec", "location": "/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Windows/prefetch/CHROME.EXE-AED7BA3D.pf", "type_indicator": "OS"}, "serial_number": 284786366, "sha256_hash": "1ac73a0134a784d92b04eb054bf1d4e950c1270d223606a94857332d0d136645", "timestamp": 1645754648649459, "timestamp_desc": "Creation Time"}
	return process
}

func NewProcessFromUserAssist(pl PlasoLog) Process {
	var process = *new(Process)
	//{"__container_type__": "event", "__type__": "AttributeContainer", "application_focus_count": 0, "application_focus_duration": 0, "data_type": "windows:registry:userassist", "date_time": {"__class_name__": "NotSet", "__type__": "DateTimeValues", "string": "Not set"}, "display_name": "OS:/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Users/coren/NTUSER.DAT", "entry_index": 1, "filename": "/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Users/coren/NTUSER.DAT", "inode": "-", "key_path": "HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\UserAssist\\{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}\\Count", "message": "[HKEY_CURRENT_USER\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\UserAssist\\{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}\\Count] UserAssist entry: 1 Value name: UEME_CTLCUACount:ctor Count: 0 Application focus count: 0 Application focus duration: 0", "number_of_executions": 0, "parser": "winreg/userassist", "pathspec": {"__type__": "PathSpec", "location": "/home/csoulet/Desktop/Workspace/projects/Plaso_test/C/Users/coren/NTUSER.DAT", "type_indicator": "OS"}, "sha256_hash": "5f8ceb89f8aa2fea7e6ad649197af0d8c703de4bb6aa91f05973aac026a7c5ca", "timestamp": 0, "timestamp_desc": "Last Time Executed", "value_name": "UEME_CTLCUACount:ctor"}

	// Assign Time and Timestamp
	var utc, _ = time.LoadLocation("UTC")
	process.Timestamp = int(pl.Timestamp)
	process.CreatedTime = time.UnixMicro(int64(pl.Timestamp)).In(utc)

	// Add Evidence
	process.Evidence = append(process.Evidence, pl.Message)

	//Parse ValueName
	// if ValueName contains ".exe" it is a path, if not it is a Application Name
	match, _ := regexp.MatchString("\\.exe$", pl.ValueName)
	if match {
		process.FullPath = pl.ValueName
		process.Filename = getFilename(pl.ValueName)
	} else {
		process.Filename = pl.ValueName
	}

	return process
}

func NewProcessFromShellBag(pl PlasoLog) Process {
	var process Process

	// Assign Time and Timestamp
	var utc, _ = time.LoadLocation("UTC")
	process.Timestamp = int(pl.Timestamp)
	process.CreatedTime = time.UnixMicro(int64(pl.Timestamp)).In(utc)

	// Add Evidence
	process.Evidence = append(process.Evidence, pl.Message)

	//Parse ShellItemName
	//<My Computer> C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe
	if strings.Index(pl.ShellItemPath, ".exe") == -1 {
		return *new(Process)
	}

	index := strings.Index(pl.ShellItemPath, " ")
	if index >= 0 {
		process.FullPath = pl.ShellItemPath[index:]
	} else {
		process.FullPath = pl.ShellItemPath
	}
	process.Filename = getFilename(process.FullPath)

	return process
}
