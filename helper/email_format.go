package helper

import (
	"io/ioutil"
	"strconv"
	"strings"
	"bufio"
	"time"
	"log"
)

//Create the email message content for the fail case report.
func FailCaseMessage(scanner *bufio.Scanner, content string) string {
	for i := 1; scanner.Scan(); i++ {
		Message := scanner.Text()
		if i == 1 || i == 2 || i == 4 { continue }
		if i == 3 {
			content += "::" + Message + "::\n"
			continue
		}
		if i == 12 { break }
		content += Message + "\n"
	}
	return content
}

//Create the email message content for the pass case repeort.
func PassCaseMessage(emailMessage string, fileNames []string, currTime time.Time) string {
	for _, fileName := range fileNames {
		serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
		data, err := ioutil.ReadFile("../lastError/lastError_" + serverType + ".txt")
		if err != nil { log.Fatal(err) }
		lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil { log.Fatal(err) }
		Message := ("Dappley - " + serverType + "\n	Last error:    " + 
					time.Unix(lastTimeStamp, 0).String() + "\n	Lastest test: "+ 
					time.Unix((currTime.Unix()), 0).String() + "\n\n")
		emailMessage += Message
	}
	return emailMessage
}