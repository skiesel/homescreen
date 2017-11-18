package sensors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type GotMail struct {
	GotMail bool `json:"gotMail"`
}

var (
	lock    = sync.Mutex{}
	gotMail = GotMail{GotMail: false}
)

func MailSensorReading(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	lock.Lock()
	err := decoder.Decode(&gotMail)
	lock.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	fmt.Fprintf(w, "ok")
}

func CheckMail() GotMail {
	lock.Lock()
	copy := gotMail
	lock.Unlock()
	return copy
}
