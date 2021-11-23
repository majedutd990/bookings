package main

import (
	"fmt"
	"github.com/majedutd990/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"io/ioutil"
	"strings"
	"time"
)

func listenForMail() {

	//	let's listen to our mail channel
	//  m := <-app.MailChan
	//  the above code is solution but a bad one
	//	it is like calling it every single time we want to send a message
	//	so what is the purpose of chan
	//	we need some procedure to run in background and listen non-stop to catch an email and send it
	//	a-synchronization
	//	 we fire an anonymous function in the background
	//we want it to listen all the time for in-coming data
	go func() {

		for {
			msg := <-app.MailChan
			sendMailMsg(msg)
		}
	}()
}

func sendMailMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	//we do not want it to be active all the time
	server.KeepAlive = false
	//	sensible time outs
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	//	 it is a dummy server if we use real one we need username and password too
	//	 we have server now let's set the client
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
		return
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
		return
	} else {
		errorLog.Println("Email Sent!")
	}
}
