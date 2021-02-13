package main

import (
	"gopkg.in/gomail.v2"
	"flag"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "fileName", "default.txt", "default txt file")
	flag.Parse()
	sendEmail(fileName)
}

func sendEmail(filename string) {
	// errMSG, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	return
	// }
	//if (strings.Contains(string(errMSG), "failure")) {
		m := gomail.NewMessage()
		m.SetHeader("From", "blockchainwarning@omnisolu.com")
		m.SetHeader("To", "blockchainwarning@omnisolu.com")
		//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Dappley Web Block Check:")
		m.SetBody("text/html", "<p> Informations in the attachment below </p>")
		m.Attach(filename)

		d := gomail.NewDialer("smtp.gmail.com", 587, "blockchainwarning@omnisolu.com", "01353751")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	//}
}