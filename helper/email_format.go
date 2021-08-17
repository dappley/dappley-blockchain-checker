package helper

import (
	"io/ioutil"
	"strconv"
	"strings"
	"bufio"
	"time"
	"log"
)

//Returns true when there is at least on true case
func ContainsFailure(test_results []bool) bool {
	final_result := false
	for _, test_result := range test_results {
		final_result = final_result || test_result
	}
	return final_result
}

func FailCaseMessage(scanner *bufio.Scanner, content string) string {
	i := 0
	for scanner.Scan() {
		Message := scanner.Text()
		i++
		if i == 1 || i == 2 || i == 4 {
			continue
		}
		if i == 3 {
			Message = "::" + Message + "::\n"
			content += Message
			continue
		}
		if i == 12 {
			break
		}
		Message = Message + "\n"
		content += Message
	}
	return content
}

func PassCaseMessage(emailMessage string, fileNames []string, currTime time.Time) string {
	for _, fileName := range fileNames {
		serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
		data, err := ioutil.ReadFile("../lastError/lastError_" + serverType + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		
		lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
		Message := ("Dappley - " + serverType + "\n	Last error:    " + time.Unix(lastTimeStamp, 0).String() +
		"\n	Lastest test: " + time.Unix((currTime.Unix()), 0).String() + "\n\n")
		emailMessage += Message
	}

	return emailMessage
}