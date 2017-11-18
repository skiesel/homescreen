package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/skiesel/homescreen/sensors"
)

type Config struct {
	WUKey            string
	WULocation       string
	GAPIKey          string
	GCalID           string
	HeadlinesRSSFeed string
	Pop3MailServer   string
	MailServerSSL    bool
	Pop3MailUsername string
	Pop3MailPassword string
}

const (
	indexFile = "templates/index.html"
)

var (
	configFilename = flag.String("config", "config.json", "config file to use")
	index          *template.Template
	config         Config
	err            error
)

func getConfig() {
	configFile, err := os.Open(*configFilename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err = index.ExecuteTemplate(w, "index.html", config)
	if err != nil {
		log.Fatal(err)
	}
}

func getMailFlag(w http.ResponseWriter, r *http.Request) {
	mail := sensors.CheckMail()
	js, err := json.Marshal(mail)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getHeadlines(w http.ResponseWriter, r *http.Request) {
	headlines := sensors.CheckHeadlines(config.HeadlinesRSSFeed)
	js, err := json.Marshal(headlines)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	messages := sensors.CheckMessages(config.Pop3MailServer, config.MailServerSSL,
		config.Pop3MailUsername, config.Pop3MailPassword)
	js, err := json.Marshal(messages)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	flag.Parse()
	getConfig()

	index, err = template.New(indexFile).Delims("[[", "]]").ParseFiles(indexFile)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/mail", getMailFlag)
	http.HandleFunc("/headlines", getHeadlines)
	http.HandleFunc("/messages", getMessages)

	http.HandleFunc("/sensor/mail", sensors.MailSensorReading)

	http.ListenAndServe(":8080", nil)
}
