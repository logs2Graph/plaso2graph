package Entity

import (
	"time"
	//"fmt"
	"strings"
)

type User struct {
	LastPasswordChange    time.Time
	LastPasswordTimestamp int
	FullName              string
	Username              string
	Comments              string
	SID                   string // Windows
	Domain                string // Windows
}

func findUser(users []User, user User) int {
	for i, u := range users {
		if u.FullName == user.FullName || u.Username == user.FullName || u.FullName == user.Username {
			return i
		}
	}
	return -1
}

func mergeUser(dest User, src User) User {
	if src.Username != "" && dest.Username == "" {
		dest.FullName = src.FullName
		dest.Username = src.Username
		dest.Comments = src.Comments
	}

	if dest.FullName == "" {
		dest.FullName = dest.Username
	}
	return dest
}

func AddUser(users []User, u *User) []User {
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

func UnionUsers(dest []User, src []User) []User {
	for _, u := range src {
		dest = AddUser(dest, &u)
	}
	return dest
}

func newUsersFromSecurity(evtx EvtxLog) (*User, *User) { //Best Effort
	var u1, u2 = new(User), new(User)

	s_name := GetDataValue(evtx, "SubjectUserName")
	s_Domain := GetDataValue(evtx, "SubjectDomainName")
	// If there is no Name, There is no user
	if s_name != "Not Found." {
		u1.FullName = strings.ToLower(s_name)
		u1.Domain = strings.ToLower(s_Domain)
	} else {
		u1 = nil
	}

	t_name := GetDataValue(evtx, "TargetUserName")
	t_Domain := GetDataValue(evtx, "TargetDomainName")
	// If there is no Name, There is no user
	if t_name != "Not Found." {
		u2.FullName = strings.ToLower(t_name)
		u2.Domain = strings.ToLower(t_Domain)
	} else {
		u2 = nil
	}

	return u1, u2
}

func NewUserFromPath(path string) *User {
	var u = new(User)
	if strings.Contains(path, "Users") {
		splitted := strings.Split(path, "\\")
		if len(splitted) == 1 {
			splitted = strings.Split(path, "/")
		}
		//fmt.Println(splitted)
		u.FullName = strings.ToLower(splitted[2])

	} else {
		return nil // Not a user
	}
	return u
}

func NewUserFromSAM(pl PlasoLog) *User {
	var user = new(User)

	user.Comments = pl.Comments
	user.FullName = pl.FullName
	user.Username = pl.Username
	user.LastPasswordTimestamp = int(pl.Timestamp)
	user.LastPasswordChange = time.Unix(int64(pl.Timestamp), 0)

	return user
}
