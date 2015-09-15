package main

import (
	"net/http"
	"html/template"
	"log"
	"renroll"
)

type Index struct {
	Conf renroll.Configuration
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("indexhandler - start")
	index := Index{Conf: renroll.Config()}
	t, _ := template.ParseFiles("idx.html", "templates/header.html", "templates/topbar.html", "templates/bottombar.html")
	log.Print("indexhandler - execute")
	t.Execute(w, index)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: renroll.Config()}
	t, _ := template.ParseFiles(
		"about.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	about := Index{Conf: renroll.Config()}
	t, _ := template.ParseFiles(
		"contact.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	t.Execute(w, about)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	conf := Index{Conf: renroll.Config()}
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

	http.Handle("/", http.FileServer(http.Dir("./")))
	if http.ListenAndServe(":80", nil) != nil {
		panic(http.ListenAndServe(":8080", nil))
	}
}
