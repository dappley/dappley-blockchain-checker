package email

import (
	"github.com/heesooh/dappley-blockchain-checker/helper"
	"io/ioutil"
	"strings"
	"bufio"
	"time"
	"fmt"
)

const (
	head = ("Jenkins Daily Digest is sent every day at 9AM " + 
			"when there hasn't been any error in last 24 Hours in dappley blockchain. \n" + 
			"[Main: http://dappley.dappworks.com/#/dappley/dashboard]\n" + 
			/* "[Mask: http://35.80.10.175/#/dappley/dashboard]\n" + */
			"[Test: http://3.16.250.102/#/dappley/dashboard]\n\n"  +
			"Detailed Info: \n\n")

	tail = ("\n\n▽ Lastest test result below ▽")
)
//Composes the email message for the dappley web blockchain test result report.
func ComposeEmail(fileNames []string) (subject string, emailMessage string) {
	currTime  := time.Now()
	testResult_Bool := make([]bool,   len(fileNames))
	testResult_Byte := make([][]byte, len(fileNames))

	//Read each files and check which one contains the failure message.
	for index, fileName := range fileNames {
		content, err := ioutil.ReadFile(fileName)
		if err != nil { return }
		//Update the content and the status
		testResult_Byte[index] = content
		testResult_Bool[index] = strings.Contains(string(content), "#  failure  detail")
	}

	//Check if the lastError directory exists, if not create one.
	helper.IsLastErrorExist(fileNames, currTime)

	//If failure exists, update the lastError file and compose a failure email; otherwise, create a daily digest email.
	if (helper.ContainsFailure(testResult_Bool)) {
		fmt.Println("Test result contains failure!")
		subject = "***Important - Dappley Web Block Check:"

		// Update the lastError timestamp to current timestamp.
		for index, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
			//If the test result contains failure, then update the lastError.txt with the current timestamp.
			if testResult_Bool[index] { helper.UpdateLastError(serverType, currTime) }
			//Create the email message.
			testResult_Scanner := bufio.NewScanner(strings.NewReader(string(testResult_Byte[index])))
			emailMessage = helper.FailCaseMessage(testResult_Scanner, emailMessage)
		}
	} else {
		fmt.Println("All tests passed!")
		before, now, after := helper.TimeFrame(currTime)

		//If current time is between the upper & lower bound, check when the last error occured.
		if (before.Before(now) && after.After(now)) {
			//If the last error occured 24 hours or more ago, create message; otherwise do nothing.
			if (helper.ItsBeen24HrForAll(fileNames, currTime)) {
				subject = "Jenkins Daily Digest:"
				fmt.Println("There hasn't been any error in last 24hrs in all three servers.")
				emailMessage = head + helper.PassCaseMessage(emailMessage, fileNames, currTime) + tail
			} else {
				fmt.Println("Errors occured less than 24 hours ago!")
			}
		} else {
			fmt.Println("Email is only deployed between 8:55AM to 9:55AM.\nCurrent time:", now)
		}
	}
	return
}