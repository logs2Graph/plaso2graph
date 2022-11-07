package Entity

import (
	. "plaso2graph/master/src"
	"strings"
)

type User struct {
	Name   string
	SID    string // Windows
	Domain string // Windows
}

func findUser(users []User, user User) int {
	for i, u := range users {
		if u.Name == user.Name {
			return i
		}
	}
	return -1
}

func mergeUser(dest User, src User) User {
	if dest.SID == "Not Found." {
		dest.SID = src.SID
	}
	return dest
}

func addUser(users []User, u *User) []User {
	if u != nil {
		i := findUser(users, *u)
		if i == -1 {
			users = append(users, *u)
		} else {
			users[i] = mergeUser(users[i], *u)
		}
	}
	return users
}

func GetUsers(data []PlasoLog) []User {
	var users []User

	for _, d := range data {
		var u1, u2 *User
		if d.EvtxLog != nil {

			if strings.Contains(d.EvtxLog.System.Provider.Name, "Sysmon") {
				switch d.EvtxLog.System.EventID {
				case 1:
					u1, u2 = newUsersFromSysmon1(*d.EvtxLog)
					break
				default:
					u1, u2 = newUsersFromSysmonDefault(*d.EvtxLog)
				}
			} else {
				u1, u2 = newUsersFromSecurity(*d.EvtxLog)
			}
			users = addUser(users, u1)
			users = addUser(users, u2)
		}
	}
	return users
}

func newUsersFromSecurity(evtx EvtxLog) (*User, *User) { //Best Effort
	//TODO: User From Security Default
	var u1, u2 = new(User), new(User)

	s_name := GetDataValue(evtx, "SubjectUserName")
	s_SID := GetDataValue(evtx, "SubjectUserSid")
	s_Domain := GetDataValue(evtx, "SubjectDomainName")
	// If there is no Name, There is no user
	if s_name != "Not Found." {
		u1.Name = s_name
		u1.SID = s_SID
		u1.Domain = s_Domain
	} else {
		u1 = nil
	}

	t_name := GetDataValue(evtx, "TargetUserName")
	t_SID := GetDataValue(evtx, "TargetUserSid")
	t_Domain := GetDataValue(evtx, "TargetDomainName")
	// If there is no Name, There is no user
	if t_name != "Not Found." {
		u2.Name = t_name
		u2.SID = t_SID
		u2.Domain = t_Domain
	} else {
		u2 = nil
	}

	return u1, u2
}

func newUsersFromSysmon1(evtx EvtxLog) (*User, *User) {
	//TODO: User From Sysmon 1
	return nil, nil
}

func newUsersFromSysmonDefault(evtx EvtxLog) (*User, *User) { //Default
	//TODO: User From Sysmon Default
	return nil, nil
}
