package Entity

import (
	"encoding/xml"
	"time"
)

type Change struct {
	Date              time.Time
	Timestamp         int
	Title             string
	Type              string
	Description       string
	UserSource        string
	UserSourceLogonID string
	UserDestination   string
	Computer          string
	Evidence          []string
}

func AddChange(cs []Change, c Change) []Change {
	if c.Title != "" {
		cs = append(cs, c)
	}
	return cs
}

func NewChangeFromEvtx(e EvtxLog) Change {
	c := Change{}

	c.Computer = e.System.Computer
	t, err := time.Parse(time.RFC3339Nano, e.System.TimeCreated.SystemTime)
	handleErr(err)
	c.Date = t
	c.Timestamp = int(t.UnixNano())

	//Get Users (if any)
	c.UserDestination = GetDataValue(e, "SubjectUserName")
	c.UserSource = GetDataValue(e, "TargetUserName")
	c.UserSourceLogonID = GetDataValue(e, "TargetLogonId")

	xml_string, err := xml.Marshal(e)
	handleErr(err)
	c.Evidence = append(c.Evidence, string(xml_string))

	switch e.System.EventID {
	case 4638:
		c.Type = "User Account Changed"
		break
	case 4704:
		c.Type = "User Right Assigned"
		break
	case 4705:
		c.Type = "User Right Removed"
		break
	case 4720:
		c.Type = "User Account Created"
		break
	case 4722:
		c.Type = "User Account Enabled"
		break
	case 4723:
		c.Type = "Password Change Attempted"
		break
	case 4724:
		c.Type = "Password Reset Attempted"
		break
	case 4725:
		c.Type = "User Account Disabled"
		break
	case 4726:
		c.Type = "User Account Deleted"
		break
	}

	return c
}
