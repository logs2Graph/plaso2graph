package Entity

import (
	"log"
	"strings"
)

func getFilename(full_path string) string {
	splitted_string := strings.Split(full_path, "\\")
	return splitted_string[len(splitted_string)-1]
}

func getExtension(filename string) string {
	splitted_string := strings.Split(filename, ".")
	if len(splitted_string) > 1 {
		return splitted_string[len(splitted_string)-1]
	}
	return ""
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
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
