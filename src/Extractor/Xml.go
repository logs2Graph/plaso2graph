package Extractor

import (
	"encoding/xml"
	"log"
	"os"
	. "plaso2graph/master/src/Entity"
)

func InitializeXmlExtractor(args map[string]interface{}) map[string]interface{} {
	var err error
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	output := args["output"].(string)
	args["output_files"] = map[string]*os.File{}
	args["output_files"].(map[string]*os.File)["process"], err = os.OpenFile(output+"/process.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["user"], err = os.OpenFile(output+"/user.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["file"], err = os.OpenFile(output+"/file.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["task"], err = os.OpenFile(output+"/task.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["computer"], err = os.OpenFile(output+"/computer.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["domain"], err = os.OpenFile(output+"/domain.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["webhistory"], err = os.OpenFile(output+"/webhistory.xml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	return args
}

func XmlExtract(data []interface{}, args map[string]interface{}) {
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	for _, d := range data {
		switch d.(type) {
		case []Process:
			InsertProcessesXml(d.([]Process), args)
			break
		case []User:
			InsertUsersXml(d.([]User), args)
			break
		case []File:
			InsertFilesXml(d.([]File), args)
			break
		case []ScheduledTask:
			InsertTasksXml(d.([]ScheduledTask), args)
			break
		case []Computer:
			InsertComputersXml(d.([]Computer), args)
			break
		case []Domain:
			InsertDomainsXml(d.([]Domain), args)
			break
		case []WebHistory:
			InsertWebHistoryXml(d.([]WebHistory), args)
			break
		}
	}

}

func InsertProcessesXml(processes []Process, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["process"]

	for _, p := range processes {
		str, err := xml.Marshal(p)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}
}

func InsertUsersXml(users []User, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["user"]

	for _, u := range users {
		str, err := xml.Marshal(u)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}
}

func InsertFilesXml(files []File, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["file"]

	for _, f := range files {
		str, err := xml.Marshal(f)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}
}

func InsertTasksXml(tasks []ScheduledTask, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["file"]

	for _, t := range tasks {
		str, err := xml.Marshal(t)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}

}

func InsertComputersXml(computers []Computer, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["computer"]

	for _, c := range computers {
		str, err := xml.Marshal(c)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}

}

func InsertDomainsXml(domains []Domain, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["domain"]

	for _, d := range domains {
		str, err := xml.Marshal(d)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}

}

func InsertWebHistoryXml(webhistory []WebHistory, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["webhistory"]

	for _, w := range webhistory {
		str, err := xml.Marshal(w)
		handleError(err)
		_, err = file.WriteString(string(str) + "\n")
		handleError(err)
	}
}
