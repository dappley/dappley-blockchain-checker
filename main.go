package main

import (
	"flag"
	"fmt"
)

func main() {
	// Maskchain server is paused temporarily
	// var email, passWord, main, test, mask string
	// flag.StringVar(&mask, "mask", "log_mask.txt", "newman log file from http://35.80.10.175: Mask Chain Server.")
	// fileNames := []string{main, mask, test}

	// Create flags
	var email, passWord, main, test string
	flag.StringVar(&email,    "email",    "default@example.com", "Email address of the sender")
	flag.StringVar(&passWord, "passWord", "default_password",    "Password of the sender email")
	flag.StringVar(&test,     "test", "default_test.txt", "newman log file from http://3.16.250.102: Test Server.")
	flag.StringVar(&main,     "main", "default_main.txt", "newman log file from http://dappley.dappworks.com: Main Server")
	flag.Parse()

	// Create the server test result email
	fileNames := []string{main, test}
	fmt.Println("Creating Email Message...")
	subject, emailMessage := composeEmail(fileNames)

	// Send out the email message
	if (subject != "" && emailMessage != "") {
		sendEmail(subject, emailMessage, fileNames, email, passWord)
		fmt.Println("Email sent!")
	}
}