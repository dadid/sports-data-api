package app

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

type emailConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"password"`
}

var (
	ch = make(chan *gomail.Message)
	// EmailConf is the universal email configuration
	EmailConf emailConfig
	_         = json.Unmarshal([]byte(os.Getenv("SBD_API_EMAIL_CONFIG")), &EmailConf)
)

func (s *Server) emailDaemon() {
	d := gomail.NewDialer(EmailConf.Host, EmailConf.Port, EmailConf.User, EmailConf.Pass)

	var sc gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-ch:
			if !ok {
				return
			}
			if !open {
				if sc, err = d.Dial(); err != nil {
					panic(err)
				}
				open = true
			}
			if err := gomail.Send(sc, m); err != nil {
				log.Println(err)
			}
		// close SMTP server connection if no email was sent in the last 30 seconds
		case <-time.After(30 * time.Second):
			if open {
				if err := sc.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		}
	}
}

func (s *Server) createSendEmail() {
	m := gomail.NewMessage()
	m.SetHeader("From", EmailConf.User)
	m.SetHeader("To", "mstockunc@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Email Notification Test.")
	m.SetBody("text/html", "Hello <b>Stock</b>!")
	// m.Attach("C:/Users/Daniel/item.jpg")
	s.sendEmail(m)
}

func (s *Server) sendEmail(m *gomail.Message) {
	d := gomail.NewDialer(EmailConf.Host, EmailConf.Port, EmailConf.User, EmailConf.Pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}