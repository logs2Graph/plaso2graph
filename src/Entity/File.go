package Entity

type File struct {
	Filename      string
	FullPath      string
	Birth_date    int
	Modified_date int
	Accessed_date int
	Changed_date  int
}

func AddFile(files []File, f File) []File {
	if f.Filename != "" {
		files = append(files, f)
	}
	return files
}

/*
func NewFileFromMFT(pl PlasoLog) File {
	var f = *new(File)
	f.Filename = pl.Filename
	f.FullPath = pl.Directory
	f.Birth_date = pl.Birth_date
	f.Modified_date = pl.Modified_date
	f.Accessed_date = pl.Accessed_date
	f.Changed_date = pl.Changed_date
	return f
}
*/
