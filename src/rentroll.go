package main
import (
	"log"
	"net/http"
	"html/template"
)

type RentRoll struct {
	Conf Configuration
}

func rentRollHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrollhandler - begin")
	conf := configuration()
	conf.GPlusSigninCallback = "gRentRoll"
	conf.FacebookSigninCallback = "fbRentRoll"
	rentroll := RentRoll{Conf: conf}
	t, _ := template.ParseFiles("rentroll.html", "templates/header-template.html", "templates/fbheader-template.html", "templates/topbar-template.html", "templates/bottombar-template.html")
	log.Print("rentrollhandler - execute")
	t.Execute(w, rentroll)
}

type TenantsTemplate struct {
	Conf Configuration
	Tenants []Tenant	
}

type Tenant struct {
	Name string
}

func tenantsHandler(w http.ResponseWriter, r *http.Request) {
	tenants := []Tenant{}
	log.Print("tenantshandler - begin")
	email := r.FormValue("email")
	if email == "" || email == "dummy@dummy.com" {
		log.Print("rentroll - NO EMAIL SET")
		tenants = []Tenant{Tenant{"#1"}}
	} else {
		dbName := email
		tenants = dbReadTenants(dbName)
	}
	t, _ := template.ParseFiles("templates/tenants-template.html")
	log.Print("tenanthandler - execute")
	t.ExecuteTemplate(w, "Tenants", TenantsTemplate{Conf: configuration(), Tenants: tenants})
}

