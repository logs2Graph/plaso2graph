package Entity

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
	"time"
)

type Event struct {
	//General Information
	Date      time.Time
	Timestamp int
	Computer  string
	Evidence  []string

	// Event Information
	Title                 string
	Type                  string
	UserSource            string
	UserDomainSource      string
	UserSourceLogonID     string
	UserDestination       string
	UserDestinationDomain string

	// Event File Information (ex: Sysmon 11,23)
	Filename  string
	FullPath  string
	Extension string
	Process   string
	ProcessId int
}

func AddEvent(cs []Event, c Event) []Event {
	if c.Title != "" {
		cs = append(cs, c)
	}
	return cs
}

func UnionEvents(a []Event, b []Event) []Event {
	for _, v := range b {
		a = AddEvent(a, v)
	}
	return a
}

func constructEvent(evtx EvtxLog) Event {
	e := Event{}

	e.Computer = evtx.System.Computer
	t, err := time.Parse(time.RFC3339Nano, evtx.System.TimeCreated.SystemTime)
	handleErr(err)
	e.Date = t
	e.Timestamp = int(t.UnixNano())
	xml_string, err := xml.Marshal(evtx)
	handleErr(err)
	e.Evidence = append(e.Evidence, string(xml_string))

	return e
}

func NewEventFromEvtx(evtx EvtxLog) Event {
	c := constructEvent(evtx)

	//Get Users (if any)
	c.UserDestination = GetDataValue(evtx, "SubjectUserName")
	c.UserSource = GetDataValue(evtx, "TargetUserName")
	c.UserSourceLogonID = GetDataValue(evtx, "TargetLogonId")

	switch evtx.System.EventID {
	case 4624:
		c.Type = "Logon"
		c.Title = "User " + c.UserSource + " logged in."
		break
	case 4625:
		c.Type = "Failed Logon"
		c.Title = "User " + c.UserSource + " Failed to log in."
		break
	case 4634:
		c.Type = "Logoff"
		c.Title = "User " + c.UserSource + " logged off."
		break
	case 4648:
		c.Type = "Explicit Credential used"
		c.Title = "User " + c.UserSource + " used explicit credentials."
		break
	case 4638:
		c.Type = "User Account Changed"
		c.Title = "User " + c.UserDestination + " changed " + c.UserSource + " account."
		log.Fatal(c.Evidence[0])
		break
	case 4704:
		c.Type = "User Right Assigned"
		c.Title = "User " + c.UserSource + " assigned " + c.UserDestination + " user right."
		log.Fatal(c.Evidence[0])
		break
	case 4705:
		c.Type = "User Right Removed"
		c.Title = "User " + c.UserSource + " removed " + c.UserDestination + " user right."
		log.Fatal(c.Evidence[0])
		break
	case 4720:
		c.Type = "User Account Created"
		c.Title = "User " + c.UserSource + "created account" + c.UserDestination + "."
		break
	case 4722:
		c.Type = "User Account Enabled"
		c.Title = "User " + c.UserSource + " enabled account" + c.UserDestination + "."
		break
	case 4723:
		c.Type = "Password Change Attempted"
		c.Title = "User " + c.UserSource + " attempted to change password."
		break
	case 4724:
		c.Type = "Password Reset Attempted"
		c.Title = "User " + c.UserSource + " attempted to reset " + c.UserDestination + "'s password."
		break
	case 4725:
		c.Type = "User Account Disabled"
		c.Title = "User " + c.UserSource + " disabled account" + c.UserDestination + "."
		break
	case 4726:
		c.Type = "User Account Deleted"
		c.Title = "User " + c.UserSource + " deleted account" + c.UserDestination + "."
		break
	}

	return c
}

func NewEventFromSysmon(e EvtxLog) Event {
	c := constructEvent(e)

	//Parse Domain and Users
	tmp := GetDataValue(e, "ParentUser")
	splitted_user := strings.Split(tmp, "\\")
	if len(splitted_user) > 1 {
		c.UserSource = splitted_user[1]
		c.UserDomainSource = splitted_user[0]
	} else {
		c.UserSource = tmp
	}

	tmp = GetDataValue(e, "User")
	splitted_user = strings.Split(tmp, "\\")
	if len(splitted_user) > 1 {
		c.UserDestination = splitted_user[1]
		c.UserDestinationDomain = splitted_user[0]
	} else {
		c.UserDestination = tmp
	}

	switch e.System.EventID {
	case 7:
		c.Type = "Image Loaded"
		log.Fatal(c.Evidence[0])
		break
	case 9:
		c.Type = "Raw Access Read"
		log.Fatal(c.Evidence[0])
		break
	case 10:
		c.Type = "Process's Memory Access"
		log.Fatal(c.Evidence[0])
		break
	case 11:
		c.Type = "File Created"

		c.Process = GetDataValue(e, "Image")
		c.ProcessId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = GetDataValue(e, "TargetFilename")
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)

		c.Title = "File " + c.FullPath + " created by " + c.Process + " (" + fmt.Sprint(c.ProcessId) + ")."
		break
	case 23:
		c.Type = "File Deleted"

		c.Process = GetDataValue(e, "Image")
		c.ProcessId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = GetDataValue(e, "TargetFilename")
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)

		c.Title = "File " + c.FullPath + " deleted by " + c.Process + " (" + fmt.Sprint(c.ProcessId) + ")."
		break
	}

	return c
}
