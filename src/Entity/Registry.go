package Entity

import (
	"time"
)

type Registry struct {
	LastModificationTime       time.Time
	LastModifictationTimestamp int
	Hive                       string
	Path                       string
	Entries                    []string
}

func AddRegistry(registries []Registry, r Registry) []Registry {
	if r.Path != "" {
		registries = append(registries, r)
	}
	return registries
}

func UnionRegistries(dest []Registry, src []Registry) []Registry {
	for _, r := range src {
		dest = AddRegistry(dest, r)
	}
	return dest
}

func NewRegistry(pl PlasoLog) Registry {
	var r = *new(Registry)

	var utc, _ = time.LoadLocation("UTC")
	r.LastModifictationTimestamp = int(pl.Timestamp)
	r.LastModificationTime = time.UnixMicro(int64(pl.Timestamp)).In(utc)

	r.Hive = pl.Filename
	r.Path = pl.KeyPath
	r.Entries = pl.Entries
	return r
}
