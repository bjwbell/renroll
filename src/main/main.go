package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bjwbell/renroll/src/renroll"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("indexhandler - start")
	index := struct{ Conf renroll.Configuration }{renroll.Config()}
	t, _ := template.ParseFiles("idx.html", "templates/header.html", "templates/topbar.html", "templates/bottombar.html")
	log.Print("indexhandler - execute")
	t.Execute(w, index)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about := struct{ Conf renroll.Configuration }{renroll.Config()}
	t, _ := template.ParseFiles(
		"about.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	about := struct{ Conf renroll.Configuration }{renroll.Config()}
	t, _ := template.ParseFiles(
		"contact.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func unreleasedHandler(w http.ResponseWriter, r *http.Request) {
	conf := struct{ Conf renroll.Configuration }{renroll.Config()}
	conf.Conf.GPlusSigninCallback = "gSettings"
	conf.Conf.FacebookSigninCallback = "fbSettings"
	t, _ := template.ParseFiles(
		"unreleased.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, conf)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	conf := struct{ Conf renroll.Configuration }{renroll.Config()}
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
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/addtenant", renroll.AddTenantHandler)
	http.HandleFunc("/auth/getemail", renroll.GetGPlusEmailHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/createaccount", renroll.CreateAccountHandler)
	http.HandleFunc("/tenantdata", renroll.TenantDataHandler)
	http.HandleFunc("/tenanthistory", renroll.TenantHistoryHandler)
	http.HandleFunc("/tenantsdata", renroll.TenantsDataHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/logerror", renroll.LogErrorHandler)
	http.HandleFunc("/oauth2callback", renroll.Oauth2callback)
	http.HandleFunc("/printinvoices", renroll.PrintInvoicesHandler)
	http.HandleFunc("/removetenant", renroll.RemoveTenantHandler)
	http.HandleFunc("/rentroll", renroll.RentRollHandler)
	http.HandleFunc("/rentrolltemplate", renroll.RentRollTemplateHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.HandleFunc("/signinform", renroll.SigninFormHandler)
	http.HandleFunc("/submit", renroll.SubmitHandler)
	http.HandleFunc("/tenant", renroll.TenantHandler)
	http.HandleFunc("/tenants", renroll.TenantsHandler)
	http.HandleFunc("/updatetenant", renroll.UpdateTenantHandler)
	http.HandleFunc("/undoremovetenant", renroll.UndoRemoveTenantHandler)
	http.HandleFunc("/undoupdatetenant", renroll.UndoUpdateTenantHandler)
	http.HandleFunc("/unreleased", unreleasedHandler)
	http.HandleFunc("/upload", renroll.UploadHandler)

	http.Handle("/", http.FileServer(http.Dir("./")))
	go func() {
		err := http.ListenAndServe(":80", nil)
		if err != nil {
			panic(http.ListenAndServe(":8080", nil))
		}
	}()

	cert := "/etc/letsencrypt/live/renroll.com/cert.pem"
	privkey := "/etc/letsencrypt/live/renroll.com/privkey.pem"
	err := http.ListenAndServeTLS(":443", cert, privkey, nil)
	if err != nil {
		cert = "./src/generate_cert/cert.pem"
		privkey = "./src/generate_cert/key.pem"
		err = http.ListenAndServeTLS(":10443", cert, privkey, nil)
		if err != nil {
			log.Print("HTTPS ListenAndServe: ", err)
		}
	}
}
