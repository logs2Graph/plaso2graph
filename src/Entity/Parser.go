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

	//Registry
	ValueName string `json:"value_name"`

	//UserAssist
	NumberOfExecutions       int `json:"number_of_executions"`
	ApplicationFocusCount    int `json:"application_focus_count"`
	ApplicationFocusDuration int `json:"application_focus_duration"`

	//ShellBags
	ShellItemPath string `json:"shell_item_path"`

	//Job
	Application string `json:"application"`
	Comment     string `json:"comment"`
	Parameters  string `json:"parameters"`

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
		t_ps, t_users, t_computers, t_domains, t_tasks, t_webhistories, t_files, t_connections, t_events := ParseEntity(line)
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

func ParseEntity(pl PlasoLog) ([]Process, []User, []Computer, []Domain, []ScheduledTask, []WebHistory, []File, []Connection, []Event) {
	var ps []Process
	var users []User
	var computers []Computer
	var domains []Domain
	var tasks []ScheduledTask
	var webhistories []WebHistory
	var files []File
	var connections []Connection
	var events []Event

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
			case 4688:
				//Extract Processes from Event Logs
				process := NewProcessFrom4688(*pl.EvtxLog)
				ps = AddProcess(ps, process)
				break
			case 4699:
				log.Fatal("4699: ", pl.Xml_string)
				break
			case 4700:
				log.Fatal("4700: ", pl.Xml_string)
				break
			case 4701:
				log.Fatal("4701: ", pl.Xml_string)
				break
			case 4702:
				log.Fatal("4702: ", pl.Xml_string)
				break
			case 4704:
				log.Fatal("4704: ", pl.Xml_string)
				break
			case 4705:
				log.Fatal("4705: ", pl.Xml_string)
				break
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
		file := NewFileFromMFT(pl)
		files = AddFile(files, file)
		break

	}

	return ps, users, computers, domains, tasks, webhistories, files, connections, events
}
