package Entity

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	//"fmt"
	"log"
	"os"
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

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ParseFile(path string) []PlasoLog {
	var output []PlasoLog
	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)
	i := 0

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if i == 0 {
			i += 1
			continue
		}
		obj := ParseLine(scanner.Text())
		output = append(output, obj)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return output
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

// ParseEntity parses the entities from the given log array.
func ParseEntity(pl []PlasoLog) ([]Process, []User, []Computer, []Domain, []ScheduledTask, []WebHistory) {
	var ps []Process
	var users []User
	var computers []Computer
	var domains []Domain
	var tasks []ScheduledTask
	var webhistories []WebHistory

	var batch_size = 100

	for _, p := range pl {

		// Merge Processes with the same name and timestamp every batch_size times.
		if len(ps)%batch_size == 0 && len(ps) != 0 {
			ps = MergeLastProcesses(ps, batch_size, 100)
		}

		switch p.DataType {
		case "windows:evtx:record":
			if strings.Contains(p.EvtxLog.System.Provider.Name, "Sysmon") {
				//TODO: Sysmon
			} else {

				// Extract Users from Event Logs
				u1, u2 := newUsersFromSecurity(*p.EvtxLog)
				users = AddUser(users, u1)
				users = AddUser(users, u2)

				// Extract Computers from Event Logs
				c1 := NewComputerFromEvtx(*p.EvtxLog)
				computers = AddComputer(computers, c1)

				// Extract Domains from Event Logs
				d1, d2 := NewDomainFromEvtx(*p.EvtxLog)
				if d1 != nil {
					domains = AddDomain(domains, *d1)
				}
				if d2 != nil {
					domains = AddDomain(domains, *d2)
				}

				switch p.EvtxLog.System.EventID {
				case 4688:
					//Extract Processes from Event Logs
					process := NewProcessFrom4688(*p.EvtxLog)
					ps = AddProcess(ps, process)
					break
				}
			}
			break
		case "windows:volume:creation":

			// Extract Process from Prefetch
			if p.Parser == "prefetch" {
				process := NewProcessFromPrefetchFile(p)
				ps = AddProcess(ps, process)
			}
			break
		case "windows:registry:userassist":

			// Extract Process from UserAssist
			process := NewProcessFromUserAssist(p)
			ps = AddProcess(ps, process)
			break
		case "windows:shell_item:file_entry":

			//Extract Process from ShellBags
			process := NewProcessFromShellBag(p)
			ps = AddProcess(ps, process)
			break
		case "windows:tasks:job":

			//Extract ScheduledTask from Task Scheduler
			task := NewScheduledTaskFromTask(p)
			tasks = AddScheduledTask(tasks, task)
			break
		case "firefox:places:page_visited":
			// Extract WebHistory from Firefox
			wh := NewWebHistoryFromFirefox(p)
			webhistories = AddWebHistory(webhistories, wh)
			break
		case "chrome:history:page_visited":
			//Extract WebHistory from Chrome
			wh := NewWebHistoryFromChrome(p)
			webhistories = AddWebHistory(webhistories, wh)
			break
		}

	}

	ps = MergeLastProcesses(ps, batch_size, 1000)

	return ps, users, computers, domains, tasks, webhistories
}
