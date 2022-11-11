package Extractor

import (
	"encoding/json"
	"log"
	"os"
	. "plaso2graph/master/src/Entity"
)

func InitializeJsonExtractor(args map[string]interface{}) map[string]interface{} {
	var err error
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	output := args["output"].(string)
	args["output_files"] = map[string]*os.File{}
	args["output_files"].(map[string]*os.File)["process"], err = os.OpenFile(output+"/process.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["user"], err = os.OpenFile(output+"/user.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["file"], err = os.OpenFile(output+"/file.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["task"], err = os.OpenFile(output+"/task.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["computer"], err = os.OpenFile(output+"/computer.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["domain"], err = os.OpenFile(output+"/domain.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	args["output_files"].(map[string]*os.File)["webhistory"], err = os.OpenFile(output+"/webhistory.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err)
	return args
}

func JsonExtract(data []interface{}, args map[string]interface{}) {
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	for _, d := range data {
		switch d.(type) {
		case []Process:
			InsertProcessesJson(d.([]Process), args)
			break
		case []User:
			InsertUsersJson(d.([]User), args)
			break
		case []File:
			InsertFilesJson(d.([]File), args)
			break
		case []ScheduledTask:
			InsertTasksJson(d.([]ScheduledTask), args)
			break
		case []Computer:
			InsertComputersJson(d.([]Computer), args)
			break
		case []Domain:
			InsertDomainsJson(d.([]Domain), args)
			break
		case []WebHistory:
			InsertWebHistoryJson(d.([]WebHistory), args)
			break
		}
	}
}

func InsertProcessesJson(processes []Process, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["process"]

	for _, process := range processes {
		json, err := json.Marshal(process)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}
}

func InsertUsersJson(users []User, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["user"]

	for _, user := range users {
		json, err := json.Marshal(user)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}
}

func InsertFilesJson(files []File, args map[string]interface{}) {
	f := args["output_files"].(map[string]*os.File)["file"]

	for _, file := range files {
		json, err := json.Marshal(file)
		handleError(err)

		_, err = f.WriteString(string(json) + "\n")
		handleError(err)
	}
}

func InsertTasksJson(tasks []ScheduledTask, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["task"]

	for _, task := range tasks {
		json, err := json.Marshal(task)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}

}

func InsertComputersJson(computers []Computer, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["computer"]

	for _, computer := range computers {
		json, err := json.Marshal(computer)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}

}

func InsertDomainsJson(domains []Domain, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["domain"]

	for _, domain := range domains {
		json, err := json.Marshal(domain)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}
}

func InsertWebHistoryJson(webhistory []WebHistory, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["webhistory"]

	for _, web := range webhistory {
		json, err := json.Marshal(web)
		handleError(err)

		_, err = file.WriteString(string(json) + "\n")
		handleError(err)
	}

}
