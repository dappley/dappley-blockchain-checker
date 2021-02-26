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

func main() {
	var fileName string
	var email string
	var passWord string
	flag.StringVar(&fileName, "fileName", "default.txt", "default txt file")
	flag.StringVar(&email, "email", "default@example.com", "default email address")
	flag.StringVar(&passWord, "passWord", "default_password", "default password")
	flag.Parse()
	subject, emailMessage := makeMessage(fileName)
	if (subject != "" && emailMessage != "") {
		fmt.Println("Both subject and emailMessages are not empty.")
		sendEmail(subject, emailMessage, fileName, email, passWord)
	}
}

func makeMessage(filename string) (string, string){
	//Create current timestamp
	timeStamp := time.Now().Unix()

	//Read the test result
	result, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", ""
	}

	if (strings.Contains(string(result), "#  failure  detail")) {
		fmt.Println("Test result contains failure.")
		//Create lastError.txt recording current timestamp
		f, err := os.Create("../lastError/lastError.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		//Write current timestamp to lastError.txt
		_, err = f.WriteString(strconv.FormatInt(timeStamp, 10))
		if err != nil {
			log.Fatal(err)
		}

		//Create the email message
		emailMessage := ""
		CommitScanner := bufio.NewScanner(strings.NewReader(string(result)))
		for CommitScanner.Scan() {
			Message := CommitScanner.Text()
			if (strings.Contains(Message, "┌─────────────────────────┬──────────────────┬──────────────────┐")) {
				break
			}
			Message = Message + "\n"
			emailMessage += Message
		}
		return "***Important - Dappley Web Block Check:", emailMessage
	} else {
		//Create layout for time.Time
		layout := "15:04:05"

		//Create lower bound
		lowerBound := "08:55:00"
		before, err := time.Parse(layout, lowerBound)
		if err != nil {
			log.Fatal(err)
		}

		//Create upper bound
		upperBound := "19:55:00"
		after, err := time.Parse(layout, upperBound)
		if err != nil {
			log.Fatal(err)
		}

		//Create current time
		curr := time.Now().Format("15:04:05")
		now, err := time.Parse(layout, curr)
		if err != nil {
			log.Fatal(err)
		}

		//If current time is between the upper and lower bound
		//Then check when last error was generated
		if (before.Before(now) && after.After(now)) {
			fmt.Println("Current time is between the upper bound and the lower bound.")
			//Read the lastError.txt
			data, err := ioutil.ReadFile("../lastError/lastError.txt")
			if err != nil {
				log.Fatal(err)
			}
			//Convert string timestamp to int64
			lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(lastTimeStamp)
			//If last error happend 24 before then create the email message
			if (timeStamp - lastTimeStamp >= 86400) {  //if the time is 9AM
				fmt.Println("Last error occured 24 hours before.")

				emailMessage := ("Jenkins Daily Digest is sent every day at 9AM " + 
								"when there hasn't been any error in last 24 Hours in dappley blockchain. " + 
								"(http://dappley.dappworks.com/#/dappley/dashboard)\n\n"  +
								"Detailed Info: \n	Last error:    " + 
								time.Unix(lastTimeStamp, 0).String() + "\n	Lastest test: " + 
								time.Unix(timeStamp, 0).String() + "\n\n▽ Lastest test result below ▽")
				return "Jenkins Daily Digest:", emailMessage
			}
		}

		fmt.Println("One or more requirements did not meet to send the \"Jenkins Daily Digest\".")
		return "", ""
	}
}

func sendEmail(subject string, emailMessage string, fileName string, email string, passWord string) {
	gmail := gomail.NewMessage()
		gmail.SetHeader("From", email)
		gmail.SetHeader("To", "blockchainwarning@omnisolu.com") /*, 
							  "wulize1994@gmail.com", 
							  "rshi@omnisolu.com", 
							  "ilshiyi@omnisolu.com") */
		gmail.SetHeader("Subject", subject)
		gmail.SetBody("text", emailMessage)
		gmail.Attach(fileName)

		d := gomail.NewDialer("smtp.gmail.com", 587, email, passWord)

		if err := d.DialAndSend(gmail); err != nil {
			fmt.Println("Unable to send out the email.")
			panic(err)
		}

	fmt.Println("Email sent!")
}