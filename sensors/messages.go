package sensors

import (
	"log"
	"net/mail"
	"regexp"
	"strings"

	"github.com/tpanum/go-pop3"
)

const (
	maxMessages = 5
)

var (
	messages = []SimpleMessage{}
)

type SimpleMessage struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func callback(number int, uid, data string, err error) (bool, error) {
	r := strings.NewReader(data)
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}

	if len(messages) == maxMessages {
		messages = messages[1:]
	}

	re := regexp.MustCompile("<.+>")
	nameOnly := strings.Replace(m.Header.Get("From"), re.FindString(m.Header.Get("From")), "", -1)

	messages = append(messages, SimpleMessage{
		From:    strings.TrimSpace(nameOnly),
		Message: m.Header.Get("Subject"),
	})

	return false, nil
}

func CheckMessages(server string, ssl bool, username, password string) []SimpleMessage {
	err := pop3.ReceiveMail(server, username, password, ssl, callback)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return messages
}
