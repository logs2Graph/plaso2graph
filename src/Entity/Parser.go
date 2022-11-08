package Entity

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
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

func ParseFile(path string) ([]Process, []User, []Computer, []Domain, []ScheduledTask, []WebHistory) {
	var ps []Process
	var users []User
	var computers []Computer
	var domains []Domain
	var tasks []ScheduledTask
	var webhistories []WebHistory
	var lines []PlasoLog
	// Batch size must be 200 for now, because smaller batch size will miss some duplicate entities (because of async nature of goroutines)
	var batch_size = 10000

	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	var wg sync.WaitGroup
	for scanner.Scan() {
		line := ParseLine(scanner.Text())
		lines = append(lines, line)
		if len(lines) == batch_size {
			wg.Wait()
			go GoParseEntity(&wg, &ps, &users, &computers, &domains, &tasks, &webhistories, lines)
			lines = *new([]PlasoLog)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	go GoParseEntity(&wg, &ps, &users, &computers, &domains, &tasks, &webhistories, lines)
	lines = *new([]PlasoLog)

	// There is race condition here, so we need to add some delay
	time.Sleep(2 * time.Millisecond)
	fmt.Println("Waiting for goroutines to finish Before Final Merge...")
	wg.Wait()
	fmt.Println("Warning: Merging Process, may take a while...")
	ps = MergeLastProcesses(ps, len(ps), 1000)

	return ps, users, computers, domains, tasks, webhistories
}

func GoParseEntity(wg *sync.WaitGroup, ps *[]Process, users *[]User, computers *[]Computer, domains *[]Domain, tasks *[]ScheduledTask, webhistories *[]WebHistory, lines []PlasoLog) {
	wg.Add(1)
	for _, line := range lines {
		t_ps, t_users, t_computers, t_domains, t_tasks, t_webhistories := ParseEntity(line)
		*ps = UnionProcesses(*ps, t_ps)
		*users = UnionUsers(*users, t_users)
		*computers = UnionComputers(*computers, t_computers)
		*domains = UnionDomains(*domains, t_domains)
		*tasks = UnionScheduledTasks(*tasks, t_tasks)
		*webhistories = UnionWebHistories(*webhistories, t_webhistories)
	}
	wg.Done()
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

func ParseEntity(pl PlasoLog) ([]Process, []User, []Computer, []Domain, []ScheduledTask, []WebHistory) {
	var ps []Process
	var users []User
	var computers []Computer
	var domains []Domain
	var tasks []ScheduledTask
	var webhistories []WebHistory

	switch pl.DataType {
	case "windows:evtx:record":
		if strings.Contains(pl.EvtxLog.System.Provider.Name, "Sysmon") {
			//TODO: Sysmon
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
			}
		}
		break
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
	}

	return ps, users, computers, domains, tasks, webhistories
}

// ParseEntity parses the entities from the given log array.
func ParseEntities(pl []PlasoLog) ([]Process, []User, []Computer, []Domain, []ScheduledTask, []WebHistory) {
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

		t_ps, t_users, t_computers, t_domains, t_tasks, t_webhistories := ParseEntity(p)
		ps = UnionProcesses(ps, t_ps)
		users = UnionUsers(users, t_users)
		computers = UnionComputers(computers, t_computers)
		domains = UnionDomains(domains, t_domains)
		tasks = UnionScheduledTasks(tasks, t_tasks)
		webhistories = UnionWebHistories(webhistories, t_webhistories)

	}

	ps = MergeLastProcesses(ps, batch_size, 1000)

	return ps, users, computers, domains, tasks, webhistories
}
