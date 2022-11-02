package src

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"strings"
)

type PlasoObject struct {
	Timestamp      float64
	Timestamp_desc string
	Source         string
	Message        string
	Parser         string
	Display_name   string
	xml_string     string
	EvtxLog        *EvtxLog
}

type EvtxLog struct {
	XMLName xml.Name `xml:"Event"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	System  struct {
		Text     string `xml:",chardata"`
		Provider struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name,attr"`
			Guid string `xml:"Guid,attr"`
		} `xml:"Provider"`
		EventID     string `xml:"EventID"`
		Version     string `xml:"Version"`
		Level       string `xml:"Level"`
		Task        string `xml:"Task"`
		Opcode      string `xml:"Opcode"`
		Keywords    string `xml:"Keywords"`
		TimeCreated struct {
			Text       string `xml:",chardata"`
			SystemTime string `xml:"SystemTime,attr"`
		} `xml:"TimeCreated"`
		EventRecordID string `xml:"EventRecordID"`
		Correlation   string `xml:"Correlation"`
		Execution     struct {
			Text      string `xml:",chardata"`
			ProcessID string `xml:"ProcessID,attr"`
			ThreadID  string `xml:"ThreadID,attr"`
		} `xml:"Execution"`
		Channel  string `xml:"Channel"`
		Computer string `xml:"Computer"`
		Security struct {
			Text   string `xml:",chardata"`
			UserID string `xml:"UserID,attr"`
		} `xml:"Security"`
	} `xml:"System"`
	EventData struct {
		Text string `xml:",chardata"`
		Data []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name,attr"`
		} `xml:"Data"`
	} `xml:"EventData"`
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

func ParseFile(path string) []PlasoObject {
	var output []PlasoObject
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

func ParseLine(data string) PlasoObject {
	var output PlasoObject
	var json_obj map[string]interface{}
	json.Unmarshal([]byte(data), &json_obj)

	//fmt.Println(json_obj["timestamp"])

	output = PlasoObject{json_obj["timestamp"].(float64),
		json_obj["timestamp_desc"].(string),
		json_obj["data_type"].(string),
		json_obj["message"].(string),
		json_obj["parser"].(string),
		json_obj["display_name"].(string),
		"", nil}

	if strings.Compare(output.Parser, "winevtx") == 0 {
		//fmt.Println(json_obj["xml_string"])
		output.EvtxLog = ParseEvtx(json_obj["xml_string"].(string))
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
