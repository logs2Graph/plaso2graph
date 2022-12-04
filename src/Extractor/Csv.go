package Extractor

import (
	"fmt"
	"log"
	"os"
	. "plaso2graph/master/src/Entity"
)

func InitializeCsvExtractor(args map[string]interface{}) map[string]interface{} {
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	output := args["output"].(string)
	args["output_files"] = map[string]*os.File{}
	args["output_files"].(map[string]*os.File)["user"] = InitializeUserCsv(output)
	args["output_files"].(map[string]*os.File)["process"] = InitializeProcessCsv(output)
	args["output_files"].(map[string]*os.File)["file"] = InitializeFileCsv(output)
	args["output_files"].(map[string]*os.File)["task"] = InitializeTaskCsv(output)
	args["output_files"].(map[string]*os.File)["computer"] = InitializeComputerCsv(output)
	args["output_files"].(map[string]*os.File)["domain"] = InitializeDomainCsv(output)
	args["output_files"].(map[string]*os.File)["webhistory"] = InitializeWebHistoryCsv(output)
	return args
}

func InitializeUserCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/user.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)

	_, err = file.WriteString("Name,Username,Domain,SID,Comments")
	return file
}

func InitializeProcessCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/process.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)

	_, err = file.WriteString("Timestamp,CreatedTime,Filename,FullPath,Commandline,PID,ParentProcessName,ParentProcessCommandline,PPID,User,UserDomain,Computer,LogonID,Evidence")
	handleError(err)

	return file
}

func InitializeFileCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/file.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)
	_, err = file.WriteString("Timestamp,Date,TimestampDesc,Filename,FullPath,Extension,IsAllocated,Evidence")
	return file
}

func InitializeTaskCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/task.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)
	_, err = file.WriteString("Application,Comment,Trigger,User")
	return file
}

func InitializeComputerCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/computer.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)

	_, err = file.WriteString("Name, Domain")
	handleError(err)

	return file
}

func InitializeDomainCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/domain.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)

	_, err = file.WriteString("Name")

	return file
}

func InitializeWebHistoryCsv(output string) *os.File {
	file, err := os.OpenFile(output+"/webhistory.csv", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err)

	_, err = file.WriteString("Timestamp,LastTimeVisited,URL,Title,User,Domain,Path,VisitCount,Evidence")
	handleError(err)
	return file
}

func CsvExtract(data []interface{}, args map[string]interface{}) {
	if args["output"] == nil {
		log.Fatal("Output directory is required")
	}

	if args["verbose"] == nil {
		args["verbose"] = false
	}

	for _, d := range data {
		switch d.(type) {
		case []Process:
			InsertProcessesCsv(d.([]Process), args)
			break
		case []User:
			InsertUsersCsv(d.([]User), args)
			break
		case []File:
			InsertFilesCsv(d.([]File), args)
			break
		case []ScheduledTask:
			InsertTasksCsv(d.([]ScheduledTask), args)
			break
		case []Computer:
			InsertComputersCsv(d.([]Computer), args)
			break
		case []Domain:
			InsertDomainsCsv(d.([]Domain), args)
			break
		case []WebHistory:
			InsertWebHistoriesCsv(d.([]WebHistory), args)
			break
		}
	}

}

func InsertProcessesCsv(processes []Process, args map[string]interface{}) {

	file := args["output_files"].(map[string]*os.File)["process"]

	for _, p := range processes {
		InsertProcessCsv(p, file)
	}
}

func InsertProcessCsv(process Process, file *os.File) {
	var str string

	str += fmt.Sprint(process.Timestamp) + ","
	str += fmt.Sprint(process.CreatedTime) + ","
	str += process.Filename + ","
	str += process.FullPath + ","
	str += process.Commandline + ","
	str += fmt.Sprint(process.PID) + ","

	str += process.ParentProcessName + ","
	str += process.ParentProcessCommandline + ","
	str += fmt.Sprint(process.PPID) + ","

	str += process.User + ","
	str += process.UserDomain + ","
	str += process.Computer + ","
	str += fmt.Sprint(process.LogonID) + ","
	str += fmt.Sprint(process.Evidence)

	_, err := file.WriteString(str)
	handleError(err)
}

func InsertUsersCsv(users []User, args map[string]interface{}) {

	file := args["output_files"].(map[string]*os.File)["user"]

	for _, u := range users {
		InsertUserCsv(u, file)
	}
}

func InsertUserCsv(user User, file *os.File) {
	var str string

	str += user.FullName + ","
	str += user.Username + ","
	str += user.Domain + ","
	str += user.SID + ","
	str += user.Comments
	str += "\n"

	_, err := file.WriteString(str)
	handleError(err)
}

func InsertFilesCsv(files []File, args map[string]interface{}) {

	file := args["output_files"].(map[string]*os.File)["file"]

	for _, f := range files {
		InsertFileCsv(f, file)
	}
}

func InsertFileCsv(file File, f *os.File) {
	var str string

	str += fmt.Sprint(file.Timestamp) + ","
	str += fmt.Sprint(file.Date) + ","
	str += file.TimestampDesc + ","
	str += file.Filename + ","
	str += file.FullPath + ","
	str += file.Extension + ","
	str += fmt.Sprint(file.IsAllocated) + ","
	str += fmt.Sprint(file.Evidence)
	str += "\n"

	_, err := f.WriteString(str)
	handleError(err)
}

func InsertTasksCsv(tasks []ScheduledTask, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["task"]

	for _, t := range tasks {
		InsertTaskCsv(t, file)
	}
}

func InsertTaskCsv(task ScheduledTask, file *os.File) {
	var str string

	str += task.Application + ","
	str += task.Comment + ","
	str += task.Trigger + ","
	str += task.User
	str += "\n"

	_, err := file.WriteString(str)
	handleError(err)
}

func InsertComputersCsv(computers []Computer, args map[string]interface{}) {

	file := args["output_files"].(map[string]*os.File)["computer"]

	for _, c := range computers {
		InsertComputerCsv(c, file)
	}
}

func InsertComputerCsv(computer Computer, file *os.File) {
	var str string

	str += computer.Name + ","
	str += computer.Domain
	str += "\n"

	_, err := file.WriteString(str)
	handleError(err)
}

func InsertWebHistoriesCsv(webhistory []WebHistory, args map[string]interface{}) {

	file := args["output_files"].(map[string]*os.File)["webhistory"]

	for _, w := range webhistory {
		InsertWebHistoryCsv(w, file)
	}
}

func InsertWebHistoryCsv(webhistory WebHistory, file *os.File) {
	var str string

	str += fmt.Sprint(webhistory.Timestamp) + ","
	str += fmt.Sprint(webhistory.LastTimeVisited) + ","
	str += webhistory.Url + ","
	str += webhistory.Title + ","
	str += webhistory.User + ","
	str += webhistory.Domain + ","
	str += webhistory.Path + ","
	str += fmt.Sprint(webhistory.VisitCount)
	str += fmt.Sprint(webhistory.Evidence)
	str += "\n"

	_, err := file.WriteString(str)
	handleError(err)
}

func InsertDomainsCsv(domains []Domain, args map[string]interface{}) {
	file := args["output_files"].(map[string]*os.File)["domain"]
	for _, d := range domains {
		InsertDomainCsv(d, file)
	}
}

func InsertDomainCsv(domain Domain, file *os.File) {
	var str string

	str += domain.Name
	str += "\n"

	_, err := file.WriteString(str)
	handleError(err)
}
