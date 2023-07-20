package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	BASE_URL    = "https://www.1secmail.com/api/v1/"
	CHECK_EMAIL = BASE_URL + "?action=getMessages&login=%s&domain=%s"
	READ_EMAIL  = BASE_URL + "?action=readMessage&login=%s&domain=%s&id=%d"
)

type Mail struct {
	Id      int    `json:"id"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	Type     string `json:"contentType"`
}

type MailBody struct {
	Mail
	Attachments []Attachment `json:"attachments"`
	Body        string       `json:"body"`
	TextBody    string       `json:"textBody"`
	HtmlBody    string       `json:"htmlBody"`
}

func GetEmail() string {
	return generateId(18) + "@1secmail.com"
}

func GetMessages(email string) []*Mail {
	var mails []*Mail

	req, err := http.NewRequest("GET", fmt.Sprintf(CHECK_EMAIL, strings.Split(email, "@")[0], strings.Split(email, "@")[1]), nil)
	if err != nil {
		return mails
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mails
	}

	err = json.NewDecoder(resp.Body).Decode(&mails)
	if err != nil {
		return mails
	}

	return mails
}

// ReadMessage will read the first message that matches the given function
// If no messages match the function, it will wait 2 seconds and try again
// This mean that this function will block until a message is found
func ReadMessage(email string, fn func(*Mail) bool) *MailBody {
	var m *Mail

	for _, mail := range GetMessages(email) {
		if fn(mail) {
			m = mail
			break
		}
	}

	if m == nil {
		time.Sleep(time.Second)
		return ReadMessage(email, fn)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(READ_EMAIL, strings.Split(email, "@")[0], strings.Split(email, "@")[1], m.Id), nil)
	resp, _ := http.DefaultClient.Do(req)

	var mail *MailBody
	json.NewDecoder(resp.Body).Decode(&mail)

	return mail
}
