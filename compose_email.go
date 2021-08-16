package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"bufio"
	"time"
	"fmt"
	"log"
	"os"
)

func composeEmail(fileNames []string) (string, string) {
	var emailMessage string
	currTime  := time.Now()
	testResult_Bool := make([]bool, len(fileNames))
	testResult_Byte := make([][]byte, len(fileNames))

	//Read each files and check which one contains the failure message
	for index, fileName := range fileNames {
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			return "", ""
		}
		//Update the content and the status
		testResult_Byte[index] = content
		testResult_Bool[index] = strings.Contains(string(content), "#  failure  detail")
	}

	// Check if the lastError directory exists
	isLastErrorExist(fileNames, currTime)

	// If failure exists, update the lastError file and compose failure email; otherwise, create daily digest email.
	if (containsFailure(testResult_Bool)) {
		fmt.Println("Test result contains failure.")

		// Update the lastError timestamp to current timestamp
		for index, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")

			//If the test result contains failure, then update
			if testResult_Bool[index] {
				UpdateLastError(serverType, currTime)
			}

			//Create the email message
			testResult_Scanner := bufio.NewScanner(strings.NewReader(string(testResult_Byte[index])))
			emailMessage = FailCaseMessage(testResult_Scanner, emailMessage)
		}
		fmt.Println("Successfully generated email message!")
		return "***Important - Dappley Web Block Check:", emailMessage
	} else {
		fmt.Println("Test result all passes!")
		before, now, after := timeFrame(currTime)

		//If current time is between the upper & lower bound, check when last error occured
		if (before.Before(now) && after.After(now)) {
			fmt.Println("Current time is between " + before.String() + "~" + after.String())
			//check if all had error 24 before, if true continue, else return "",""
			if (itsBeen24HrForAll(fileNames, currTime)) {
				fmt.Println("There hasn't been any error in last 24hrs in all three servers.")
				emailMessage = PassCaseMessage(emailMessage, fileNames, currTime)
				head := ("Jenkins Daily Digest is sent every day at 9AM " + 
				"when there hasn't been any error in last 24 Hours in dappley blockchain. \n" + 
				"[Main: http://dappley.dappworks.com/#/dappley/dashboard]\n" + 
				/* "[Mask: http://35.80.10.175/#/dappley/dashboard]\n" + */
				"[Test: http://3.16.250.102/#/dappley/dashboard]\n\n"  +
				"Detailed Info: \n\n")
				tail := ("\n\n▽ Lastest test result below ▽")
				emailMessage = head + emailMessage + tail
				fmt.Println("Successfully generated email message!")
				return "Jenkins Daily Digest:", emailMessage
			}
		}
		fmt.Println("One or more condition did not match.")
		return "", ""
	}
}

//-------------------Helper--------------------
func isExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//If the lastError folder does not exists in the repository, assume that the lastError occured
//24 hours prior to the current testing time then create lastError directory and lastError files.
func isLastErrorExist(fileNames []string, currTime time.Time) {
	lastErrorExists := isExist("../lastError")
	if !lastErrorExists {
		fmt.Println("Creating the lastError directory...")
		err := os.Mkdir("../lastError", os.ModePerm)
		if err != nil {
			fmt.Println("Could not create the lastError directory!")
			panic(err)
		}
		for _, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
			yesterday  := currTime.AddDate(0, 0, -1)
			UpdateLastError(serverType, yesterday)
		}
	}
}

//Returns true when there is at least on true case
func containsFailure(test_results []bool) bool {
	final_result := false
	for _, test_result := range test_results {
		final_result = final_result || test_result
	}
	return final_result
}

//Update/Create a lastError file with the current timestamp
func UpdateLastError(serverType string, currTime time.Time) {
	//Create/OverRide the file
	f, err := os.Create("../lastError/lastError_" + serverType + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	//Write current timestamp to the lastError file
	_, err = f.WriteString(strconv.FormatInt(currTime.Unix(), 10))
	if err != nil {
		log.Fatal(err)
	}
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