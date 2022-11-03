package src

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"strings"
)

type PlasoLog2 struct {
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
	EvtxLog       *EvtxLog
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
