package email

import (
	"github.com/heesooh/dappley-blockchain-checker/helper"
	"io/ioutil"
	"strings"
	"bufio"
	"time"
	"fmt"
)

func ComposeEmail(fileNames []string) (string, string) {
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
	helper.IsLastErrorExist(fileNames, currTime)

	// If failure exists, update the lastError file and compose failure email; otherwise, create daily digest email.
	if (helper.ContainsFailure(testResult_Bool)) {
		fmt.Println("Test result contains failure.")

		// Update the lastError timestamp to current timestamp
		for index, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")

			//If the test result contains failure, then update
			if testResult_Bool[index] {
				helper.UpdateLastError(serverType, currTime)
			}

			//Create the email message
			testResult_Scanner := bufio.NewScanner(strings.NewReader(string(testResult_Byte[index])))
			emailMessage = helper.FailCaseMessage(testResult_Scanner, emailMessage)
		}
		fmt.Println("Successfully generated email message!")
		return "***Important - Dappley Web Block Check:", emailMessage
	} else {
		fmt.Println("Test result all passes!")
		before, now, after := helper.TimeFrame(currTime)

		//If current time is between the upper & lower bound, check when last error occured
		if (before.Before(now) && after.After(now)) {
			fmt.Println("Current time is between " + before.String() + "~" + after.String())
			//check if all had error 24 before, if true continue, else return "",""
			if (helper.ItsBeen24HrForAll(fileNames, currTime)) {
				fmt.Println("There hasn't been any error in last 24hrs in all three servers.")
				emailMessage = helper.PassCaseMessage(emailMessage, fileNames, currTime)
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