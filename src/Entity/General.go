package Entity

import (
	"log"
	"strings"
)

func getFilename(full_path string) string {
	splittedString := strings.Split(full_path, "\\")
	return splittedString[len(splittedString)-1]
}

func getExtension(filename string) string {
	splittedString := strings.Split(filename, ".")
	if len(splittedString) > 1 {
		return splittedString[len(splittedString)-1]
	}
	return ""
}

func handleErr(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func GetUsernameFromPath(path string) string {
	if strings.Contains(path, "Users") {
		splitted := strings.Split(path, "\\")
		if len(splitted) == 1 {
			splitted = strings.Split(path, "/")
		}
		if len(splitted) > 2 {
			return splitted[2]
		} else {
			return ""
		}
	}
	return ""
}
