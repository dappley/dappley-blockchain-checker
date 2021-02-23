package main

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
	"bufio"
	"flag"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "fileName", "default.txt", "default txt file")
	flag.Parse()
	sendEmail(fileName)
}

func sendEmail(filename string) {
	result, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	if (strings.Contains(string(result), "failure")) {
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
		m := gomail.NewMessage()
		m.SetHeader("From", "blockchainwarning@omnisolu.com")
		m.SetHeader("To", "blockchainwarning@omnisolu.com", "wulize1994@gmail.com", "rshi@omnisolu.com", "ilshiyi@omnisolu.com")
		//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "***Important - Dappley Web Block Check:")
		m.SetBody("text", emailMessage)
		m.Attach(filename)

		d := gomail.NewDialer("smtp.gmail.com", 587, "blockchainwarning@omnisolu.com", "01353751")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}
}