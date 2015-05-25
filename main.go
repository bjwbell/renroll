package main
import (
	"fmt"
	"encoding/json"
	"net/http"
	"html/template"
	"log"
	"os"
	"net/smtp"
)

type Configuration struct {
	GmailAddress           string
	GmailPassword          string
	GoogleClientId         string
	GoogleClientSecret     string
	GooglePlusScopes       string	
	GoogleAnalyticsId      string
	FacebookScopes         string
	FacebookAppId          string
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

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	configuration := configuration()
	emailAddress := r.FormValue("email")
	body := "To: " + configuration.GmailAddress + "\r\nSubject: Renroll Notification Clicked!" + 
		"\r\n\r\nInterested user " + emailAddress + "."
	auth := smtp.PlainAuth("", configuration.GmailAddress, configuration.GmailPassword, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, emailAddress,
		[]string{configuration.GmailAddress},[]byte(body))
	if err != nil {
		fmt.Println("error1:", err)
		log.Fatal(err)
	}   	
	http.Redirect(w, r, "/", http.StatusFound)
}

type Index struct {
	Conf Configuration
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("indexhandler - start")
	index := Index{Conf: configuration()}
	t, _ := template.ParseFiles("idx.html")
	log.Print("indexhandler - execute")
	t.Execute(w, index)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: configuration()}
	t, _ := template.ParseFiles("about.html")
	t.Execute(w, about)
}



func main() {
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/sendemail", sendEmailHandler)
	http.HandleFunc("/oauth2callback", oauth2callback)
	http.HandleFunc("/rentroll", rentRollHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	if http.ListenAndServe(":80", nil) != nil {
		panic(http.ListenAndServe(":8080", nil))
	}
}
