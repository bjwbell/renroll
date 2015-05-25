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
)

type AccessToken struct {
 	Token  string
 	Expiry int64
}

func googleOAuth2Config(domain string) *oauth2.Config {
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
 	conf := googleOAuth2Config(domain(r))
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(oauth2.NoContext, tok)
 	response, err := client.Get("https://www.googleapis.com/plus/v1/people/me?fields=emails")
 	// handle err. You need to change this into something more robust
 	// such as redirect back to home page with error message
 	if err != nil {
 		w.Write([]byte(err.Error()))
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
	log.Print("oauth2callback - response: " + str)
	dec := json.NewDecoder(strings.NewReader(str))
	var m OAuth2Response
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}
	for _, v := range m.Emails {
		log.Print("oauth2callback - email (value, type): " + v.Value + ", " + v.Type)
	}
	
	email := "dummy@dummy.com"
	if len(m.Emails) != 1 {
		log.Print("oauth2callback - NO VALID EMAIL OR TOO MANY")
		
	} else {	
		email = m.Emails[0].Value
	}
	if newAccount == "true" {
		log.Print("oauth2callback - NEW ACCOUNT")
		dbCreate(email)
		dbInsert(email, "#1")
	}
	w.Write([]byte(email))
}


