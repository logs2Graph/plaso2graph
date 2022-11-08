package Entity

import ()

type Domain struct {
	Name string
}

func containsDomain(ds []Domain, d Domain) bool {
	for _, v := range ds {
		if v.Name == d.Name {
			return true
		}
	}
	return false
}

func AddDomain(ds []Domain, d Domain) []Domain {
	if containsDomain(ds, d) == false {
		ds = append(ds, d)
	}
	return ds
}

func UnionDomains(dest []Domain, src []Domain) []Domain {
	for _, p := range src {
		dest = AddDomain(dest, p)
	}
	return dest
}

func GetDomain(data []PlasoLog) []Domain {
	var res []Domain

	for _, d := range data {
		switch d.DataType {
		case "windows:evtx:record":
			src, dest := NewDomainFromEvtx(*d.EvtxLog)
			if src != nil {
				res = AddDomain(res, *src)
			}

			if dest != nil {
				res = AddDomain(res, *dest)
			}
			break
		}
	}

	return res
}

func NewDomainFromEvtx(evtx EvtxLog) (*Domain, *Domain) {
	var src, dst *Domain = nil, nil

	tmp := GetDataValue(evtx, "TargetDomainName")
	if tmp != "Not Found." {
		dst = &Domain{Name: tmp}
	}

	tmp = GetDataValue(evtx, "SubjectDomainName")
	if tmp != "Not Found." {
		src = &Domain{Name: tmp}
	}

	return src, dst
}
