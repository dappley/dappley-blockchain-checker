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
	"os"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "fileName", "default.txt", "default txt file")
	flag.Parse()
	sendEmail(fileName)
}

func sendEmail(filename string) {
	timeStamp := time.Now().Unix()

	result, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	if (strings.Contains(string(result), "failure")) {
		f, err := os.Create("lastError.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.WriteString(strconv.FormatInt(timeStamp, 10))
		if err != nil {
			log.Fatal(err)
		}

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
		failEmail := gomail.NewMessage()
		failEmail.SetHeader("From", "blockchainwarning@omnisolu.com")
		failEmail.SetHeader("To", "blockchainwarning@omnisolu.com", 
								  "wulize1994@gmail.com", 
								  "rshi@omnisolu.com", 
								  "ilshiyi@omnisolu.com")
		//failEmail.SetAddressHeader("Cc", "dan@example.com", "Dan")
		failEmail.SetHeader("Subject", "***Important - Dappley Web Block Check:")
		failEmail.SetBody("text", emailMessage)
		failEmail.Attach(filename)

		d := gomail.NewDialer("smtp.gmail.com", 587, "blockchainwarning@omnisolu.com", "01353751")

		if err := d.DialAndSend(failEmail); err != nil {
			panic(err)
		}
	} else {
		layout := "15:04:05"

		lowerBound := "08:55:00"
		before, err := time.Parse(layout, lowerBound)
		if err != nil {
			log.Fatal(err)
		}

		upperBound := "09:55:00"
		after, err := time.Parse(layout, upperBound)
		if err != nil {
			log.Fatal(err)
		}

		curr := time.Now().Format("15:04:05")
		now, err := time.Parse(layout, curr)
		if err != nil {
			log.Fatal(err)
		}

		if (before.Before(now) && after.After(now)) {
			data, err := ioutil.ReadFile("lastError.txt")
			if err != nil {
				log.Fatal(err)
			}
			lastTimeStamp, err := strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			if (timeStamp - lastTimeStamp <= 86400) {  //if the time is 9AM
				dailyDigest := gomail.NewMessage()
				dailyDigest.SetHeader("From", "blockchainwarning@omnisolu.com")
				dailyDigest.SetHeader("To", "blockchainwarning@omnisolu.com",
											"wulize1994@gmail.com", 
											"rshi@omnisolu.com", 
											"ilshiyi@omnisolu.com")
				dailyDigest.SetHeader("Subject", "Jenkins Daily Digest:")
				emailMessage := ("Jenkins Daily Digest is sent every day at 9AM " + 
								"when there hasn't been any error in last 24 Hours in dappley blockchain. " + 
								"(http://dappley.dappworks.com/#/dappley/dashboard)\n\n"  +
								"Detailed Info: \n	Last error:    " + 
								time.Unix(lastTimeStamp, 0).String() + "\n	Lastest test: " + 
								time.Unix(timeStamp, 0).String() + "\n\n▽ Lastest test result below ▽")
				dailyDigest.SetBody("text", emailMessage)
				dailyDigest.Attach(filename)
	
				d := gomail.NewDialer("smtp.gmail.com", 587, "blockchainwarning@omnisolu.com", "01353751")
	
				if err := d.DialAndSend(dailyDigest); err != nil {
					panic(err)
				}
			}
		}
		}
}