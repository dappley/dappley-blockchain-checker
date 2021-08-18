package email

import (
	"github.com/heesooh/dappley-blockchain-checker/helper"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
	"bufio"
	"fmt"
	"log"
)

//Send out the dappley web blockchain test result to recipients specified in the recipients.txt file.
func SendEmail(subject string, emailMessage string, fileNames []string, email string, passWord string) {
	var recipients []string

	file_byte, err := ioutil.ReadFile("recipients.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
	for scanner.Scan() {
		line := scanner.Text()
		if !helper.Valid_email(line) {
			fmt.Println("Invalid email address: \"" + line + "\"")
			continue
		}
		recipients = append(recipients, line)
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", email)
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = mail.FormatAddress(recipient, "")
	}
	mail.SetHeader("To", addresses...)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text", emailMessage)
	for _, fileName := range fileNames {
		mail.Attach(fileName)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, email, passWord)
	if err := d.DialAndSend(mail); err != nil {
		fmt.Println("Unable to send out the email.")
		panic(err)
	}
}