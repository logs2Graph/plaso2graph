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
	UserSourceDomain      string
	UserSourceLogonID     string
	UserDestination       string
	UserDestinationDomain string

	// Event File Information (ex: Sysmon 11,23,10, etc)
	Filename        string
	FullPath        string
	Extension       string
	ProcessSource   string
	ProcessSourceId int
	ProcessTarget   string
	ProcessTargetId int

	// For events involving a group
	GroupName   string
	GroupDomain string
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
	xmlString, err := xml.Marshal(evtx)
	handleErr(err)
	e.Evidence = append(e.Evidence, string(xmlString))

	return e
}

func swapUser(c Event) Event {
	temp := c.UserDestination
	c.UserDestination = c.UserSource
	c.UserSource = temp

	temp = c.UserDestinationDomain
	c.UserDestinationDomain = c.UserSourceDomain
	c.UserSourceDomain = temp
	return c
}

func NewEventFromEvtx(evtx EvtxLog) Event {
	c := constructEvent(evtx)

	//Get Users (if any)
	c.UserDestination = GetDataValue(evtx, "SubjectUserName")
	c.UserDestinationDomain = GetDataValue(evtx, "SubjectDomainName")
	c.UserSource = GetDataValue(evtx, "TargetUserName")
	c.UserSourceDomain = GetDataValue(evtx, "TargetDomainName")
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
	case 4649:
		log.Fatal("Event 4649")
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
		//c = swapUser(c)
		c.Title = "User " + c.UserSource + " created account " + c.UserDestination + "."
		break
	case 4722:
		c.Type = "User Account Enabled"
		//c = swapUser(c)
		c.Title = "User " + c.UserSource + " enabled account" + c.UserDestination + "."
		break
	case 4723:
		c.Type = "Password Change Attempted"

		// Switch UserDestination and UserSource
		tmp := c.UserDestination
		c.UserDestination = c.UserSource
		c.UserSource = tmp

		c.Title = "User " + c.UserDestination + " attempted to change " + c.UserSource + "'s password."
		break
	case 4724:
		c.Type = "Password Reset Attempted"

		// Switch UserDestination and UserSource
		tmp := c.UserDestination
		c.UserDestination = c.UserSource
		c.UserSource = tmp

		c.Title = "User " + c.UserDestination + " attempted to reset " + c.UserSource + "'s password."
		break
	case 4725:
		c.Type = "User Account Disabled"
		//c = swapUser(c)
		c.Title = "User " + c.UserDestination + " disabled account" + c.UserSource + "."
		break
	case 4726:
		c.Type = "User Account Deleted"
		c.Title = "User " + c.UserDestination + " deleted account" + c.UserSource + "."
		break
	case 4727:
		c.Type = "Security Enabled Global Group Created"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceLogonID = ""
		c.Title = c.UserDestination + " created group " + c.GroupName + "."
		break
	case 4728:
		c.Type = "Member added to global security group"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = GetDataValue(evtx, "MemberName")
		c.UserSourceLogonID = GetDataValue(evtx, "MemberSid")
		c.Title = "User " + c.UserSource + " added to group " + c.GroupName + " by " + c.UserDestination + "."
		break
	case 4729:
		c.Type = "Member removed from global security group"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = GetDataValue(evtx, "MemberName")
		c.UserSourceLogonID = GetDataValue(evtx, "MemberSid")
		c.Title = "User " + c.UserSource + " removed from group " + c.GroupName + " by " + c.UserDestination + "."
		break
	case 4730:
		c.Type = "Security Enabled Global Group Deleted"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceLogonID = ""
		c.Title = c.UserDestination + " deleted group " + c.GroupName + "."
		break
	case 4731:
		c.Type = "Security Enabled Local Group Created"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceLogonID = ""
		c.Title = c.GroupName + " group was created	 by " + c.UserDestination + "."
		break
	case 4732:
		c.Type = "Member added to Local Security Group"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceDomain = ""
		c.Title = c.UserDestination + " Added to local security group " + c.GroupDomain + "."
		break
	case 4733:
		c.Type = "Member removed from Local Security Group"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = GetDataValue(evtx, "MemberName")
		c.UserSourceLogonID = GetDataValue(evtx, "MemberSid")
		c.Title = c.UserSource + " Removed from local security group " + c.GroupDomain + " by " + c.UserDestination + "."
		break
	case 4734:
		c.Type = "Security Enabled Local Group Deleted"
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceLogonID = ""
		c.Title = c.UserDestination + " deleted group " + c.GroupName + "."
		break
	case 4735:
		c.Type = "Security-enabled local group was changed."
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceDomain = ""
		c.Title = "Global security group " + c.UserSource + " was changed by " + c.UserDestination + "."
		break
	case 4737:
		c.Type = "Security-enabled global group was changed."
		c.GroupName = c.UserSource
		c.GroupDomain = c.UserSourceDomain
		c.UserSource = ""
		c.UserSourceDomain = ""
		c.Title = "Global security group " + c.UserSource + " was changed by " + c.UserDestination + "."
		break
	case 4738:
		c.Type = "User Account Changed"
		// switch User Dest and User Source
		c.Title = "User " + c.UserDestination + " changed " + c.UserSource + " account."
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
		c.UserSourceDomain = splitted_user[0]
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
		c.ProcessSource = strings.ToLower(GetDataValue(e, "Image"))
		c.ProcessSourceId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = strings.ToLower(GetDataValue(e, "ImageLoaded"))
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)
		c.UserSource = ""
		c.UserSourceDomain = ""
		c.Title = "Process " + c.ProcessSource + " loaded file " + c.FullPath + "."
		break
	case 9:
		log.Println("Sysmon 9: Oportunity to TEST")
		c.Type = "Raw Access Read"
		c.ProcessSource = strings.ToLower(GetDataValue(e, "Image"))
		c.ProcessSourceId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = strings.ToLower(GetDataValue(e, "ImageLoaded"))
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)
		c.UserSource = ""
		c.UserSourceDomain = ""
		c.Title = "Process " + c.ProcessSource + " read file " + c.FullPath + "."
		break
	case 10:
		c.Type = "Process's Memory Access"
		c.ProcessSource = strings.ToLower(GetDataValue(e, "SourceImage"))
		c.ProcessSourceId = convertOct(GetDataValue(e, "SourceProcessId"))

		c.ProcessTarget = strings.ToLower(GetDataValue(e, "TargetImage"))
		c.ProcessTargetId = convertOct(GetDataValue(e, "TargetProcessId"))
		c.Title = "Process " + c.ProcessSource + " accessed memory of " + c.ProcessTarget + "."
		break
	case 11:
		c.Type = "File Created"
		c.ProcessSource = strings.ToLower(GetDataValue(e, "Image"))
		c.ProcessSourceId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = GetDataValue(e, "TargetFilename")
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)

		c.Title = "File " + c.FullPath + " created by " + c.ProcessSource + " (" + fmt.Sprint(c.ProcessSourceId) + ")."
		break
	case 23:
		c.Type = "File Deleted"
		c.ProcessSource = strings.ToLower(GetDataValue(e, "Image"))
		c.ProcessSourceId = convertOct(GetDataValue(e, "ProcessId"))
		c.FullPath = GetDataValue(e, "TargetFilename")
		c.Filename = getFilename(c.FullPath)
		c.Extension = getExtension(c.Filename)

		c.Title = "File " + c.FullPath + " deleted by " + c.ProcessSource + " (" + fmt.Sprint(c.ProcessSourceId) + ")."
		break
	}

	return c
}
