package renroll

import (
	"encoding/json"
	"net/http"
	"log"
	"os"
	"net/smtp"
	"fmt"
	"runtime"
)

type Configuration struct {
	GmailAddress           string
	GmailPassword          string
	GoogleClientId         string
	GoogleClientSecret     string
	GooglePlusScopes       string
	GPlusSigninCallback    string
	GoogleAnalyticsId      string
	FacebookScopes         string
	FacebookAppId          string
	FacebookSigninCallback    string
}

func Config() Configuration {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}

func domain(r *http.Request) string {
	return r.Host
}

func SendAdminEmail(emailAddress string, subject string, body string) {
	configuration := Config()
	content :=
		"To: " + configuration.GmailAddress + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body
	auth := smtp.PlainAuth("", configuration.GmailAddress, configuration.GmailPassword, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, emailAddress,
		[]string{configuration.GmailAddress},[]byte(content))
	if err != nil {
		log.Print("sendEmail - ERROR:")
		log.Fatal(err)
	}
}

func InterestedUser(emailAddress string, loginMethod string) {
	subject := "Renroll Notification Clicked (" + emailAddress + ")"
	body := "Interested user: " + emailAddress + ".\r\n" +
		"Login method: " + loginMethod + "."
	SendAdminEmail(emailAddress, subject, body)
}

func LogErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("logerrorhandler - begin")
	error := r.FormValue("error")
	SendAdminEmail(Config().GmailAddress, "Renroll JS Error", error)
}

func logError(error string) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	trace := fmt.Sprintf("%s", buf)
	msg := "Go Error\r\nError Message: " + error + "\r\n\r\nStack Trace:\r\n" + trace
	log.Print("logError:\n" + error)
	SendAdminEmail(Config().GmailAddress, "Renroll Go Error", msg)
}
