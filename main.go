package main
import (
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
func interestedUser(emailAddress string, loginMethod string) {
	configuration := configuration()
	body := "To: " + configuration.GmailAddress + "\r\nSubject: Renroll Notification Clicked!" + 
		"\r\n\r\nInterested user: " + emailAddress + ".\r\nLogin method: " + loginMethod + "."
	auth := smtp.PlainAuth("", configuration.GmailAddress, configuration.GmailPassword, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, emailAddress,
		[]string{configuration.GmailAddress},[]byte(body))
	if err != nil {
		log.Print("sendEmail - ERROR:")
		log.Fatal(err)
	}   	
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {

	emailAddress := r.FormValue("email")
	interestedUser(emailAddress, "sendEmailHandler")
	http.Redirect(w, r, "/", http.StatusFound)
}

type Index struct {
	Conf Configuration
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("indexhandler - start")
	index := Index{Conf: configuration()}
	t, _ := template.ParseFiles("idx.html", "header-template.html", "fbheader-template.html", "topbar-template.html", "bottombar-template.html")
	log.Print("indexhandler - execute")
	t.Execute(w, index)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: configuration()}
	t, _ := template.ParseFiles("about.html", "header-template.html", "fbheader-template.html", "topbar-template.html", "bottombar-template.html")
	t.Execute(w, about)
}



func main() {
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/sendemail", sendEmailHandler)
	http.HandleFunc("/oauth2callback", oauth2callback)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/auth/getemail", getGPlusEmailHandler)
	http.HandleFunc("/tenants", tenantsHandler)
	http.HandleFunc("/rentroll", rentRollHandler)
	http.HandleFunc("/createaccount", createAccountHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	if http.ListenAndServe(":80", nil) != nil {
		panic(http.ListenAndServe(":8080", nil))
	}
}
