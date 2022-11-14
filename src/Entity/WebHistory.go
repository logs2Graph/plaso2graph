package Entity

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

type WebHistory struct {
	LastTimeVisited time.Time
	Timestamp       int
	Url             string
	Domain          string
	Path            string
	Title           string
	User            string
	VisitCount      int
	Evidence        []string
}

func AddWebHistory(whs []WebHistory, wh WebHistory) []WebHistory {
	if wh.Url != "" {
		whs = append(whs, wh)
	}
	return whs
}

func UnionWebHistories(dest []WebHistory, src []WebHistory) []WebHistory {
	for _, wh := range src {
		dest = AddWebHistory(dest, wh)
	}
	return dest
}

func NewWebHistoryFromFirefox(pl PlasoLog) WebHistory {
	var wh = *new(WebHistory)
	wh.Url = pl.Url
	wh.Title = pl.Title
	wh.Domain = pl.Host

	var utc, _ = time.LoadLocation("UTC")
	wh.Timestamp = int(pl.Timestamp)
	wh.LastTimeVisited = time.UnixMicro(int64(wh.Timestamp)).In(utc)
	wh.VisitCount = pl.VisitCount

	wh.Evidence = append(wh.Evidence, pl.Message)

	u := NewUserFromPath(pl.Filename)
	if u != nil || u.Name != "" {
		wh.User = u.Name
	} else {
		log.Println("Error parsing user from path: ", pl.Filename)
	}

	return wh
}

func NewWebHistoryFromChrome(pl PlasoLog) WebHistory {
	var wh = *new(WebHistory)

	wh.Url = pl.Url
	wh.Title = pl.Title

	var utc, _ = time.LoadLocation("UTC")
	wh.Timestamp = int(pl.Timestamp)
	wh.LastTimeVisited = time.UnixMicro(int64(pl.Timestamp / 1000000)).In(utc)

	wh.VisitCount = int(pl.TypedCount)

	wh.Evidence = append(wh.Evidence, pl.Message)

	u := NewUserFromPath(pl.Filename)
	if u != nil {
		wh.User = u.Name
	}

	r, err := regexp.Compile("http(?:s|)://(?P<domain>[^/]+)(?P<path>.*)")
	handleErr(err)
	matches := r.FindStringSubmatch(pl.Url)
	if len(matches) == 3 {
		wh.Domain = matches[1]
		wh.Path = matches[2]
	} else {
		log.Println("Error parsing Domain and Path from Url: ", pl.Url, ": ", fmt.Sprint(matches))
	}

	return wh
}
