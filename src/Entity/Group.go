package Entity

import (
	"encoding/xml"
)

type Group struct {
	Name     string
	Domain   string
	Computer string
	Evidence []string
}

func findGroup(groups []Group, group Group) int {
	for i, g := range groups {
		if g.Name == group.Name && g.Domain == group.Domain {
			return i
		}
	}
	return -1
}

func AddGroup(groups []Group, g Group) []Group {
	i := findGroup(groups, g)
	if i == -1 {
		groups = append(groups, g)
	}
	return groups
}

func UnionGroups(dest []Group, src []Group) []Group {
	for _, g := range src {
		dest = AddGroup(dest, g)
	}
	return dest
}

func NewGroupFromSecurity(evtx EvtxLog) Group {
	g := Group{}
	g.Computer = evtx.System.Computer
	g.Name = GetDataValue(evtx, "TargetUserName")
	g.Domain = GetDataValue(evtx, "TargetDomainName")

	xml_bytes, _ := xml.Marshal(evtx)

	g.Evidence = append(g.Evidence, string(xml_bytes))
	return g
}
