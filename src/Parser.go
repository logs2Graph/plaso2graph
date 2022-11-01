package src

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type PlasoObject struct {
	Datetime      time.Time
	Datetime_desc string
	Source        string
	Message       string
	Parser        string
	Display_name  string
	Tag           string
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ParseFile(path string) []PlasoObject {
	var output []PlasoObject
	file, _ := os.Open(path)
	scanner := bufio.NewScanner(file)
	i := 0

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		if i == 0 {
			i += 1
			continue
		}
		obj := ParseLine(scanner.Text())
		output = append(output, obj)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return output
}

func ParseLine(data string) PlasoObject {
	var output PlasoObject
	obj := strings.Split(data, ",")
	t, _ := time.Parse(time.RFC3339, obj[0])
	output = PlasoObject{t, obj[1], obj[3], obj[4], obj[5], obj[6], obj[7]}

	return output
}

func PlasoObjectPrint(data PlasoObject) {
	fmt.Println("Source: ", data.Source)
	fmt.Println("Message: ", data.Message)
	fmt.Println("Parser: ", data.Parser)
	fmt.Println("Display Name: ", data.Display_name)
}
