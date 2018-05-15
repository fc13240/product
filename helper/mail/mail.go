package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

const sender = "service@51helper.com"
const password = "Lixuegang123"
const hostname = "smtp.exmail.qq.com"
const template = "To:%s\r\nFrom:%s\r\nSubject:%s\r\nContent-Type:text/html;charset=UTF-8\r\n\r\n%s"

func SendMail(username, subject, body string) error {
	msg := fmt.Sprintf(template, username, sender, subject, body)

	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(fmt.Sprint(hostname,":25"), auth, sender, []string{username}, []byte(msg))

	if err != nil {
		log.Printf("smtp.SendMail(%s) failed: %s", username, err)
	}
	return nil
}

func Send(username, subject, body string) error {
	msg := fmt.Sprintf(template, username, "51helper", subject, body)

	auth := smtp.PlainAuth("", sender, password, hostname)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         hostname,
	}

	conn, err := tls.Dial("tcp", hostname+":465", tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, hostname)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(sender); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(username); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

	return nil
}
