package Entity

import (
	"strings"
	"time"
)

type File struct {
	Filename      string
	FullPath      string
	Timestamp     int
	TimestampDesc string
	IsAllocated   bool
	Date          time.Time
	Extension     string
	Evidence      []string
}

func AddFile(files []File, f File) []File {
	if f.Filename != "" {
		files = append(files, f)
	}
	return files
}

func UnionFiles(dest []File, src []File) []File {
	for _, f := range src {
		dest = AddFile(dest, f)
	}
	return dest
}

func NewFileFromMFT(pl PlasoLog) File {
	var f = *new(File)
	if len(pl.PathHints) > 0 {
		f.FullPath = pl.PathHints[0]
	}
	// Parse the filename from the full path
	tmp_split := strings.Split(f.FullPath, "\\")
	f.Filename = tmp_split[len(tmp_split)-1]

	// Parse the extension from the filename
	tmp_split = strings.Split(f.Filename, ".")
	f.Extension = tmp_split[len(tmp_split)-1]

	var utc, _ = time.LoadLocation("UTC")
	f.Timestamp = int(pl.Timestamp)
	f.Date = time.UnixMicro(int64(pl.Timestamp / 1000)).In(utc)

	f.TimestampDesc = pl.TimestampDesc
	f.Evidence = append(f.Evidence, pl.Message)
	f.IsAllocated = pl.IsAllocated

	return f
}
