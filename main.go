package main

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
	"strconv"
	"bufio"
	"time"
	"flag"
	"log"
	"fmt"
	"os"
)

//--------------------Core--------------------

func main() {
	//initialize flags
	var email, passWord, main, test, mask string
	flag.StringVar(&email, "email", "default@example.com", "Email address of the sender")
	flag.StringVar(&passWord, "passWord", "default_password", "Password of the sender email")
	flag.StringVar(&test, "test", "log_test.txt", "newman log file from http://54.176.241.99: Test Server.")
	flag.StringVar(&mask, "mask", "log_mask.txt", "newman log file from http://35.80.10.175: Mask Chain Server.")
	flag.StringVar(&main, "main", "log_main.txt", "newman log file from http://dappley.dappworks.com: Main Server")
	flag.Parse()

	fileNames := []string{main, mask, test}

	//Email content
	fmt.Println("Creating Email Message...")
	subject, emailMessage := makeMessage(fileNames)
	if (subject != "" && emailMessage != "") {
		sendEmail(subject, emailMessage, fileNames, email, passWord)
	} else {
		fmt.Println("Unable to send Email!")
	}
}

func makeMessage(fileNames []string) (string, string){
	//Initialize variables
	var emailMessage string
	currTime  := time.Now()
	testResult_Bool := make([]bool, 3)
	testResult_Byte := make([][]byte, 3)

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

	//If at least one of the file contains the failure message, update the time stamp of last error
	//and create the email message with detailed info
	if (containFail(testResult_Bool)) {
		fmt.Println("Test result contains failure.")
		//Create lastError.txt recording current timestamp
		for index, fileName := range fileNames {
			serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")

			//If the test result contains failure, then update
			if testResult_Bool[index] {
				UpdateLastError(serverType, currTime)
			}

			//Create the email message
			CommitScanner := bufio.NewScanner(strings.NewReader(string(testResult_Byte[index])))
			emailMessage = FailCaseMessage(CommitScanner, emailMessage)
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
				"[Mask: http://35.80.10.175/#/dappley/dashboard]\n" + 
				"[Test: http://54.176.241.99/#/dappley/dashboard]\n\n"  +
				"Detailed Info: \n\n")
				tail := ("\n\n▽ Lastest test result below ▽")
				emailMessage = head + emailMessage + tail
				fmt.Println("Successfully generated email message!")
				return "Jenkins Daily Digest:", emailMessage
			}
		}
		fmt.Println("One or more condition does not match.")
		return "", ""
	}
}

func sendEmail(subject string, emailMessage string, fileNames []string, email string, passWord string) {
	gmail := gomail.NewMessage()
		gmail.SetHeader("From", email)
		gmail.SetHeader("To", "blockchainwarning@omnisolu.com") /*, 
							  "wulize1994@gmail.com", 
							  "rshi@omnisolu.com", 
							  "ilshiyi@omnisolu.com")*/
		gmail.SetHeader("Subject", subject)
		gmail.SetBody("text", emailMessage)

		for _, fileName := range fileNames {
			gmail.Attach(fileName)
		}

		d := gomail.NewDialer("smtp.gmail.com", 587, email, passWord)

		if err := d.DialAndSend(gmail); err != nil {
			fmt.Println("Unable to send out the email.")
			panic(err)
		}

	fmt.Println("Email sent!")
}

//-------------------Helper--------------------

//Returns true when there is at least on true case
func containFail(booList []bool) bool {
	return booList[0] || booList[1] || booList[2]
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

func timeFrame(currTime time.Time) (before, now, after time.Time) {
	//Create layout for time.Time
	layout := "15:04:05"

	//Create lower bound
	lowerBound := "08:55:00"
	before, err := time.Parse(layout, lowerBound)
	if err != nil {
		log.Fatal(err)
	}

	//Create upper bound
	upperBound := "09:55:00"
	after, err = time.Parse(layout, upperBound)
	if err != nil {
		log.Fatal(err)
	}

	//Create current time
	curr := currTime.Format("15:04:05")
	now, err = time.Parse(layout, curr)
	if err != nil {
		log.Fatal(err)
	}
	
	return
}

func itsBeen24HrForAll(fileNames []string, currTime time.Time) bool {
	itsBeen24Hr   := true
	currTimeStamp := currTime.Unix()
	for _, fileName := range fileNames {
		serverType := strings.TrimSuffix(strings.TrimPrefix(fileName, "log_"), ".txt")
		data, err := ioutil.ReadFile("../lastError/lastError_" + serverType + ".txt")
		if err != nil {
			UpdateLastError(serverType, currTime)
			itsBeen24Hr = itsBeen24Hr && false
			continue
		}
		lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		if (currTimeStamp - lastTimeStamp >= 86400) {
			itsBeen24Hr = itsBeen24Hr && true
			continue
		}
		itsBeen24Hr = itsBeen24Hr && false
	}
	return itsBeen24Hr
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