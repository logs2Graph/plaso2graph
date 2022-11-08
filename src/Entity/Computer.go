package Entity

import ()

type Computer struct {
	Name   string
	Domain string
}

func containsComputer(cs []Computer, c Computer) bool {
	for _, v := range cs {
		if v.Name == c.Name && v.Domain == c.Domain {
			return true
		}
	}
	return false
}

func AddComputer(cs []Computer, c Computer) []Computer {
	if containsComputer(cs, c) == false {
		cs = append(cs, c)
	}
	return cs
}

func UnionComputers(dest []Computer, src []Computer) []Computer {
	for _, p := range src {
		dest = AddComputer(dest, p)
	}
	return dest
}

func GetComputer(data []PlasoLog) []Computer {
	var res []Computer

	for _, d := range data {
		switch d.DataType {
		case "windows:winevtx:record":
			res = AddComputer(res, NewComputerFromEvtx(*d.EvtxLog))
		}
	}

	return res
}

func NewComputerFromEvtx(evtx EvtxLog) Computer {
	var res Computer

	res.Name = evtx.System.Computer

	return res
}
