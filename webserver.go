// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"net/smtp"
	"encoding/json"
	"os"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("submit.html")
	if err != nil {
		return
	}
	fmt.Fprintf(w, "%s", body)
}

type Configuration struct {
	GmailAddress     string
	GmailPassword    string
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error1:", err)
		log.Fatal(err)
	}
	emailAddress := r.FormValue("email")
	body := "To: " + configuration.GmailAddress + "\r\nSubject: Renroll Notification Clicked!" + 
		"\r\n\r\nInterested user " + emailAddress + "."
	auth := smtp.PlainAuth("", configuration.GmailAddress, configuration.GmailPassword, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, emailAddress,
		[]string{configuration.GmailAddress},[]byte(body))
	if err != nil {
		fmt.Println("error1:", err)
		log.Fatal(err)
	}   	
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/sendemail", sendEmailHandler)
	http.Handle("/", http.FileServer(http.Dir("./")))
	panic(http.ListenAndServe(":80", nil))
}
