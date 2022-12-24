package Entity

import (
	"encoding/json"
	"encoding/xml"
	//"fmt"
	"log"
	"strings"
)

type EvtxLog struct {
	System struct {
		Provider struct {
			Name string `xml:"Name,attr"`
			Guid string `xml:"Guid,attr"`
		} `xml:"Provider"`
		EventID     int    `xml:"EventID"`
		Version     string `xml:"Version"`
		TimeCreated struct {
			Text       string `xml:",chardata"`
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
		Execution struct {
			Text      string `xml:",chardata"`
			ProcessID string `xml:"ProcessID,attr"`
			ThreadID  string `xml:"ThreadID,attr"`
		} `xml:"Execution"`
		EventRecordID string `xml:"EventRecordID"`
		Correlation   string `xml:"Correlation"`
		Channel       string `xml:"Channel"`
		Computer      string `xml:"Computer"`
		Security      struct {
			Text   string `xml:",chardata"`
			UserID string `xml:"UserID,attr"`
		} `xml:"Security"`
	} `xml:"System"`
	EventData struct {
		Data []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name,attr"`
		} `xml:"Data"`
	} `xml:"EventData"`
}

type PlasoLog struct {
	ContainerType string  `json:"__container_type__"`
	Type          string  `json:"__type__"`
	DataType      string  `json:"data_type"`
	DisplayName   string  `json:"display_name"`
	KeyPath       string  `json:"key_path"`
	Message       string  `json:"message"`
	Parser        string  `json:"parser"`
	Path          string  `json:"path"`
	Sha256Hash    string  `json:"sha256_hash"`
	Timestamp     float64 `json:"timestamp"`
	TimestampDesc string  `json:"timestamp_desc"`
	Xml_string    string  `json:"xml_string"`
	Filename      string  `json:"filename"`
	Name          string  `json:"name"`

	// Prefetch
	Executable string `json:"executable"`

	// Link
	EnvVarLocation string `json:"env_var_location"`

	// Bam
	BinaryPath string `json:"binary_path"`

	//Registry
	ValueName string `json:"value_name"`

	//UserAssist
	NumberOfExecutions       int      `json:"number_of_executions"`
	ApplicationFocusCount    int      `json:"application_focus_count"`
	ApplicationFocusDuration int      `json:"application_focus_duration"`
	Entries                  []string `json:"entries"`

	//ShellBags
	ShellItemPath string `json:"shell_item_path"`

	//Job
	Application string `json:"application"`
	Comment     string `json:"comment"`
	Parameters  string `json:"parameters"`

	// PE
	PeType       string `json:"pe_type"`
	ImportedHash string `json:"imphash"`

	// Service
	ServiceType  int    `json:"service_type"`
	StartType    int    `json:"start_type"`
	ErrorControl int    `json:"error_control"`
	ImagePath    string `json:"image_path"`
	ServiceDll   string `json:"service_dll"`
	ObjectName   string `json:"object_name"`

	//Web History
	Url        string `json:"url"`
	Title      string `json:"title"`
	TypedCount int    `json:"typed_count"`
	Host       string `json:"host"`
	VisitCount int    `json:"visit_count"`

	//MFT
	PathHints     []string `json:"path_hints"`
	IsAllocated   bool     `json:"is_allocated"`
	AttributeType string   `json:"attribute_type"`

	//SRUM
	BackgroundBytesRead    int `json:"background_bytes_read"`
	BackgroundBytesWritten int `json:"background_bytes_written"`
	ForegroundBytesRead    int `json:"foreground_bytes_read"`
	ForegroundBytesWritten int `json:"foreground_bytes_written"`

	// SAM
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Comments string `json:"comments"`

	//Evtx
	EvtxLog *EvtxLog
}

func GetDataValue(evtx EvtxLog, name string) string {
	data := evtx.EventData.Data

	for _, d := range data {
		if strings.Compare(d.Name, name) == 0 {
			return d.Text
		}
	}
	return "Not Found."
}

func ParseEntities(data []interface{}, lines []PlasoLog) []interface{} {
	for _, line := range lines {
		t_ps, t_scriptblocks, t_users, t_groups, t_computers, t_domains, t_tasks, t_services, t_webhistories, t_files, t_connections, t_events, t_registries := ParseEntity(line)
		for i, _ := range data {
			switch data[i].(type) {
			case []Process:
				data[i] = UnionProcesses(data[i].([]Process), t_ps)
				break
			case []User:
				data[i] = UnionUsers(data[i].([]User), t_users)
				break
			case []Computer:
				data[i] = UnionComputers(data[i].([]Computer), t_computers)
				break
			case []Domain:
				data[i] = UnionDomains(data[i].([]Domain), t_domains)
				break
			case []ScheduledTask:
				data[i] = UnionScheduledTasks(data[i].([]ScheduledTask), t_tasks)
				break
			case []Service:
				data[i] = UnionServices(data[i].([]Service), t_services)
				break
			case []WebHistory:
				data[i] = UnionWebHistories(data[i].([]WebHistory), t_webhistories)
				break
			case []File:
				data[i] = UnionFiles(data[i].([]File), t_files)
				break
			case []Connection:
				data[i] = UnionConnections(data[i].([]Connection), t_connections)
				break
			case []Event:
				data[i] = UnionEvents(data[i].([]Event), t_events)
				break
			case []Registry:
				data[i] = UnionRegistries(data[i].([]Registry), t_registries)
				break
			case []Group:
				data[i] = UnionGroups(data[i].([]Group), t_groups)
				break
			case []ScriptBlock:
				data[i] = UnionScriptBlocks(data[i].([]ScriptBlock), t_scriptblocks)
				break
			}
		}
	}
	return data
}

func ParseLine(data string) PlasoLog {
	var output PlasoLog
	json.Unmarshal([]byte(data), &output)

	//fmt.Println(json_obj["timestamp"])

	if output.Parser == "winevtx" {
		//fmt.Println(json_obj["xml_string"])
		output.EvtxLog = ParseEvtx(output.Xml_string)
	}

	return output
}

func ParseEvtx(data string) *EvtxLog {
	var evtxLog EvtxLog

	//fmt.Println([]byte(data))
	xml.Unmarshal([]byte(data), &evtxLog)
	//fmt.Println(fmt.Sprint(evtxLog))
	return &evtxLog
}

func ParseEntity(pl PlasoLog) ([]Process, []ScriptBlock, []User, []Group, []Computer, []Domain, []ScheduledTask, []Service, []WebHistory, []File, []Connection, []Event, []Registry) {
	var ps []Process
	var scriptblocks []ScriptBlock
	var users []User
	var computers []Computer
	var domains []Domain
	var tasks []ScheduledTask
	var webhistories []WebHistory
	var files []File
	var connections []Connection
	var events []Event
	var services []Service
	var registries []Registry
	var groups []Group

	switch pl.DataType {
	case "windows:evtx:record":
		if strings.Contains(pl.EvtxLog.System.Provider.Name, "Sysmon") {
			switch pl.EvtxLog.System.EventID {
			case 1:
				ps = AddProcess(ps, NewProcessFromSysmon1(*pl.EvtxLog))
				break
			case 3:
				connections = AddConnection(connections, NewConnectionFromSysmon3(*pl.EvtxLog))
				break
			default:
				event := NewEventFromSysmon(*pl.EvtxLog)
				events = AddEvent(events, event)
			}

		} else {

			// Extract Users from Event Logs
			u1, u2 := newUsersFromSecurity(*pl.EvtxLog)
			users = AddUser(users, u1)
			users = AddUser(users, u2)

			// Extract Computers from Event Logs
			c1 := NewComputerFromEvtx(*pl.EvtxLog)
			computers = AddComputer(computers, c1)

			// Extract Domains from Event Logs
			d1, d2 := NewDomainFromEvtx(*pl.EvtxLog)
			if d1 != nil {
				domains = AddDomain(domains, *d1)
			}
			if d2 != nil {
				domains = AddDomain(domains, *d2)
			}

			switch pl.EvtxLog.System.EventID {
			case 4673:
				log.Fatal("4673 : ", pl.Xml_string)
				break
			case 4627:
				log.Fatal("4627 : ", pl.Xml_string)
				break
			case 4688:
				//Extract Processes from Event Logs
				process := NewProcessFrom4688(*pl.EvtxLog)
				ps = AddProcess(ps, process)
				break
			case 4103:
				// Extract Scheduled Tasks from Event Logs
				scriptblock := NewScriptBlockFrom4103(*pl.EvtxLog)
				scriptblocks = AddScriptBlock(scriptblocks, scriptblock)
				break
			case 4104:
				// Handle Powershell Script Block
				scriptblock := NewScriptBlockFrom4104(*pl.EvtxLog)
				scriptblocks = AddScriptBlock(scriptblocks, scriptblock)
				break
			case 4699:
				log.Println("4699: Found but not parsed - ", pl.Xml_string)
				break
			case 4700:
				log.Println("4700: Found but not parsed - ", pl.Xml_string)
				break
			case 4701:
				log.Println("4701: Found but not parsed - ", pl.Xml_string)
				break
			case 4702:
				log.Println("4702: Found but not parsed - ", pl.Xml_string)
				break
			case 4704:
				log.Println("4704: Found but not parsed - ", pl.Xml_string)
				break
			case 4705:
				log.Println("4705: Found but not parsed - ", pl.Xml_string)
				break
			case 4728:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4729:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4731:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4732:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4733:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4735:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 4737:
				e := NewEventFromEvtx(*pl.EvtxLog)
				g := NewGroupFromSecurity(*pl.EvtxLog)
				events = AddEvent(events, e)
				groups = AddGroup(groups, g)
				break
			case 5131:
				log.Fatal("5131 : ", pl.Xml_string)
			default:
				e := NewEventFromEvtx(*pl.EvtxLog)
				events = AddEvent(events, e)

			}
		}
	case "windows:volume:creation":

		// Extract Process from Prefetch
		if pl.Parser == "prefetch" {
			process := NewProcessFromPrefetchFile(pl)
			ps = AddProcess(ps, process)
		}
		break

	case "windows:prefetch:execution":
		// Extract Process from Prefetch
		process := NewProcessFromPrefetchExecution(pl)
		ps = AddProcess(ps, process)
		break
	case "windows:lnk:link":
		// Extract Process from LNK
		process := NewProcessFromLink(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:registry:amcache":
		// Extract Process from Amcache
		process := NewProcessFromAmCache(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:registry:appcompatcache":
		// Extract Process from AppCompatCache
		process := NewProcessFromAppCompatCache(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:registry:bagmru":
		//log.Fatal("BagMRU: Found but not parsed - ", pl.Xml_string)
		// TODO: Handle MRU with File or Folder?
		break

	case "windows:registry:bam":
		process := NewProcessFromBAM(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:registry:mrulist":
		//log.Fatal("MRUList: Found but not parsed - ", pl.Xml_string)
		break

	case "windows:registry:mrulistex":
		//log.Fatal("MRUListEx: Found but not parsed - ", pl.Xml_string)
		break

	case "windows:srum:application_usage":
		process := NewProcessFromSRUM(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:srum:network_usage":
		//log.Fatal("Network Usage: Found but not parsed - ", pl.Xml_string)
		break

	case "windows:srum:network_connectivity":
		//log.Fatal("Network Connectivity: Found but not parsed - ", pl.Xml_string)
		break

	case "windows:registry:run":
		//log.Fatal("Registry Run: Found but not parsed - ", pl.Xml_string)
		registries = AddRegistry(registries, NewRegistry(pl))
		break

	case "windows:registry:service":
		//log.Fatal("Service: Found but not parsed - ", pl.Xml_string)
		service := NewService(pl)
		services = AddService(services, service)
		break

	case "task_scheduler:task_cache:entry":
		//log.Fatal("Task Cache: Found but not parsed - ", pl.Xml_string)
		// TODO: Add Scheduled Task Runs Entity
		break

	case "windows:registry:sam_users":
		// Create User From SAM Registry
		user := NewUserFromSAM(pl)
		users = AddUser(users, user)
		break

	case "pe":
		file := NewFileFromPE(pl)
		files = AddFile(files, file)
		break

	case "windows:registry:userassist":
		// Extract Process from UserAssist
		process := NewProcessFromUserAssist(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:shell_item:file_entry":
		//Extract Process from ShellBags
		process := NewProcessFromShellBag(pl)
		ps = AddProcess(ps, process)
		break

	case "windows:tasks:job":
		//Extract ScheduledTask from Task Scheduler
		task := NewScheduledTaskFromTask(pl)
		tasks = AddScheduledTask(tasks, task)
		break

	case "firefox:places:page_visited":
		// Extract WebHistory from Firefox
		wh := NewWebHistoryFromFirefox(pl)
		webhistories = AddWebHistory(webhistories, wh)
		break
	case "chrome:history:page_visited":
		//Extract WebHistory from Chrome
		wh := NewWebHistoryFromChrome(pl)
		webhistories = AddWebHistory(webhistories, wh)
		break
	case "fs:stat:ntfs":
		//Extract File from MFT
		//file := NewFileFromMFT(pl)
		//files = AddFile(files, file)
		break

	}

	return ps, scriptblocks, users, groups, computers, domains, tasks, services, webhistories, files, connections, events, registries
}
