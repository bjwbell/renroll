package main
import (
	"encoding/json"
	"net/http"
	"html/template"
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

func configuration() Configuration {
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

func sendAdminEmail(emailAddress string, subject string, body string) {
	configuration := configuration()
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

func interestedUser(emailAddress string, loginMethod string) {
	subject := "Renroll Notification Clicked (" + emailAddress + ")"
	body := "Interested user: " + emailAddress + ".\r\n" +
		"Login method: " + loginMethod + "."
	sendAdminEmail(emailAddress, subject, body)
}

func logErrorHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("logerrorhandler - begin")
	error := r.FormValue("error")
	sendAdminEmail(configuration().GmailAddress, "Renroll JS Error", error)
}

func logError(error string) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	trace := fmt.Sprintf("%s", buf)
	msg := "Go Error\r\nError Message: " + error + "\r\n\r\nStack Trace:\r\n" + trace
	log.Print(error)
	sendAdminEmail(configuration().GmailAddress, "Renroll Go Error", msg)
}

type Index struct {
	Conf Configuration
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("indexhandler - start")
	index := Index{Conf: configuration()}
	t, _ := template.ParseFiles("idx.html", "templates/header.html", "templates/topbar.html", "templates/bottombar.html")
	log.Print("indexhandler - execute")
	t.Execute(w, index)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: configuration()}
	t, _ := template.ParseFiles(
		"about.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: configuration()}
	t, _ := template.ParseFiles(
		"contact.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	conf := Index{Conf: configuration()}
	conf.Conf.GPlusSigninCallback = "gSettings"
	conf.Conf.FacebookSigninCallback = "fbSettings"
	t, _ := template.ParseFiles(
		"settings.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, conf)
}

func main() {
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/logerror", logErrorHandler)
	http.HandleFunc("/oauth2callback", oauth2callback)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/auth/getemail", getGPlusEmailHandler)
	http.HandleFunc("/tenants", tenantsHandler)
	http.HandleFunc("/rentroll", rentRollHandler)
	http.HandleFunc("/rentrolltemplate", rentRollTemplateHandler)
	http.HandleFunc("/createaccount", createAccountHandler)
	http.HandleFunc("/signinform", signinFormHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	if http.ListenAndServe(":80", nil) != nil {
		panic(http.ListenAndServe(":8080", nil))
	}
}
