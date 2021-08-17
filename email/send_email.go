package email

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/mail"
	"strings"
	"bufio"
	"fmt"
	"log"
)

func SendEmail(subject string, emailMessage string, fileNames []string, email string, passWord string) {
	var recipients []string

	file_byte, err := ioutil.ReadFile("recipients.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
	for scanner.Scan() {
		line := scanner.Text()
		if !valid_email(line) {
			fmt.Println("Invalid email address: \"" + line + "\"")
			continue
		}
		recipients = append(recipients, line)
	}

	gmail := gomail.NewMessage()
	gmail.SetHeader("From", email)
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = gmail.FormatAddress(recipient, "")
	}
	gmail.SetHeader("To", addresses...)
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
}

//-------------------Helper--------------------
func valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}