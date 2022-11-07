package Entity

import (
	. "plaso2graph/master/src"
)

type WebHistory struct {
	LastTime       string
	Timestamp      int
	Url            string
	Domain         string
	Title          string
	DownloadedFile string
	User           string
	Evidence       []string
}

func NewWebHistoryFromMozilla(pl PlasoLog) WebHistory {
	var wh = *new(WebHistory)

	return wh
}

func NewWebHistoryFromChrome(pl PlasoLog) WebHistory {
	var wh = *new(WebHistory)

	return wh
}
