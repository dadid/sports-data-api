package app

import (
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

// EmailConfig holds the email daemon config
type EmailConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"password"`
}

// EmailSender holds an email config and runs an email sending daemon
type EmailSender struct {
	Conf EmailConfig
}

var ch = make(chan *gomail.Message)

func (e *EmailSender) emailDaemon() {
	d := gomail.NewDialer(e.Conf.Host, e.Conf.Port, e.Conf.User, e.Conf.Pass)

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

// CreateSendEmail sends an email
func (e *EmailSender) CreateSendEmail() {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Conf.User)
	m.SetHeader("To", "dadidnl@gmail.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Email Notification Test.")
	m.SetBody("text/html", "Hello <b>Testing</b>!")
	// m.Attach("C:/Users/Daniel/item.jpg")
	e.sendEmailOnChannel(m)
}

func (e *EmailSender) sendEmail(m *gomail.Message) {
	d := gomail.NewDialer(e.Conf.Host, e.Conf.Port, e.Conf.User, e.Conf.Pass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func (e *EmailSender) sendEmailOnChannel(m *gomail.Message) {
	ch <- m
}