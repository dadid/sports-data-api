package email

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

type emailconfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"password"`
}

var (
	ch = make(chan *gomail.Message)
	// Emailconf is the universal email configuration
	Emailconf emailconfig
	_         = json.Unmarshal([]byte(os.Getenv("SBD_API_EMAIL_CONFIG")), &Emailconf)
)

func emailDaemon() {
	d := gomail.NewDialer(Emailconf.Host, Emailconf.Port, Emailconf.User, Emailconf.Pass)

	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-ch:
			if !ok {
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					panic(err)
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				log.Println(err)
			}
		// close SMTP server connection if no email was sent in the last 30 seconds
		case <-time.After(30 * time.Second):
			if open {
				if err := s.Close(); err != nil {
					panic(err)
				}
				open = false
			}
		}
	}
}

func createSendEmail() {
	m := gomail.NewMessage()
	m.SetHeader("From", Emailconf.User)
	m.SetHeader("To", "mstockunc@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Email Notification Test.")
	m.SetBody("text/html", "Hello <b>Stock</b>!")
	// m.Attach("C:/Users/Daniel/item.jpg")
	sendEmail(m)
}

func sendEmail(m *gomail.Message) {
	d := gomail.NewDialer(Emailconf.Host, Emailconf.Port, Emailconf.User, Emailconf.Pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}