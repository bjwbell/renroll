// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
import (
	"fmt"
	"encoding/json"
	"net/http"
	"log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"strings"
	"time"
	"html/template"
)

func googleOAuth2Config() *oauth2.Config {
	appConf := configuration()
	conf := &oauth2.Config{
 		ClientID:     appConf.GoogleClientId,
		ClientSecret: appConf.GoogleClientSecret,
 		RedirectURL:  "postmessage",
		Scopes:       []string{},
		Endpoint: google.Endpoint,
 	}
	return conf
}

func readHttpBody(response *http.Response) string {
	fmt.Println("Reading body")
 	bodyBuffer := make([]byte, 5000)
 	var str string
 	count, err := response.Body.Read(bodyBuffer)
 	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {
 		if err != nil {

 		}
 		str += string(bodyBuffer[:count])
 	}
 	return str
 }

func oauth2callback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		log.Print("oauth2callback - NO CODE")
		email := "dummy@dummy.com"
		w.Write([]byte(email))
		return
	}
	log.Print("oauth2callback - url: " + r.URL.RawQuery)
	log.Print("oauth2callback - code: " + code)
	
	newAccount := r.FormValue("new_account")
	//email := getGPlusEmail(code)
	token := getGPlusToken(r)
	email := getGPlusEmail(token)
	if newAccount == "true" {
		log.Print("oauth2callback - NEW ACCOUNT")
		interestedUser(email, "oauth2callback")
		dbCreate(email)
		dbInsert(email, "#1")
	}
	w.Write([]byte(email))
}

func getGPlusToken(r *http.Request) oauth2.Token {
        accessToken := r.FormValue("access_token")
        if accessToken == "" {
		log.Fatal("getGPlusToken - NO ACCESS TOKEN")
	}
	return oauth2.Token{AccessToken: accessToken,
		TokenType: "Bearer",
		RefreshToken: "",
		Expiry: time.Now().Add(time.Hour)}
}

func getGPlusEmail(tok oauth2.Token) string {
 	conf := googleOAuth2Config()
	client := conf.Client(oauth2.NoContext, &tok)
 	response, err := client.Get("https://www.googleapis.com/plus/v1/people/me?fields=emails")
 	// handle err. You need to change this into something more robust
 	// such as redirect back to home page with error message
 	if err != nil {
		log.Print("getGPlusEmail - COULDN'T GET PROFILE INFO WITH TOK, ERR:")
		log.Fatal(err)
 	}
 	str := readHttpBody(response)
	type Email struct {
		Value string
		Type string
	}
	type OAuth2Response struct {
		Kind string
		Emails []Email
		Id string
	}
	log.Print("getemail - response: " + str)
	dec := json.NewDecoder(strings.NewReader(str))
	var m OAuth2Response
	if err := dec.Decode(&m); err != nil {
		log.Print("getGPlusEmail - COULDN'T DECODE OAUTH2 RESPONSE, ERR:")
		log.Fatal(err)
	}
	for _, v := range m.Emails {
		log.Print("getemail - email (value, type): " + v.Value + ", " + v.Type)
	}
	
	email := "dummy@dummy.com"
	if len(m.Emails) != 1 {
		log.Print("getemail - NO VALID EMAIL OR TOO MANY")
		
	} else {	
		email = m.Emails[0].Value
	}
	return email
}

func getGPlusEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := getGPlusEmail(getGPlusToken(r))
	w.Write([]byte(email))
}


func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	loginMethod := r.FormValue("loginmethod")
	if email == "" {
		log.Print("createAccountHandler - NO EMAIL")
		success := "false"
		w.Write([]byte(success))
		return
	}
	log.Printf("createAccountHandler - NEW ACCOUNT: %s", email)
	interestedUser(email, loginMethod)
	dbCreate(email)
	dbInsert(email, "#1")
	success := "true"
	w.Write([]byte(success))
}


type SigninForm struct {
	Conf Configuration
}

func signinFormHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("signinform - begin")
	t, _ := template.ParseFiles("signinform-template.html")
	log.Print("signinform - execute")
	t.Execute(w, SigninForm{Conf: configuration()})
}
