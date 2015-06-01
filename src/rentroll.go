package main
import (
	"log"
	"net/http"
	"html/template"
	"time"
	"strconv"
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
	t, _ := template.ParseFiles(
		"rentroll.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	log.Print("rentrollhandler - execute")
	t.Execute(w, rentroll)
}

type TenantsTemplate struct {
	Conf Configuration
	Tenants []Tenant
	AsOfDateDay string
	AsOfDateMonth string
	AsOfDateYear  string
}

type Tenant struct {
	Name string
	Address string
	SqFt int
	LeaseStartDate string
	LeaseEndDate string
	BaseRent string
	Electricity string
	Gas string
	Water string
	SewageTrashRecycle string
	Comments string
}

func tenantsHandler(w http.ResponseWriter, r *http.Request) {
	tenants := []Tenant{}
	log.Print("tenantshandler - begin")
	email := r.FormValue("email")
	if email == "" || email == "dummy@dummy.com" {
		log.Print("rentroll - NO EMAIL SET")
		tenants = []Tenant{Tenant{"#1", "", 0, "", "", "", "", "", "", "", ""}}
	} else {
		dbName := email
		tenants = dbReadTenants(dbName)
	}
	t, _ := template.ParseFiles("templates/tenants.html")
	log.Print("tenanthandler - execute")
	tenantsTemplate := TenantsTemplate{
		Conf: configuration(),
		Tenants: tenants,
		AsOfDateDay: strconv.Itoa(time.Now().Day()),
		AsOfDateMonth: time.Now().Month().String(),
		AsOfDateYear: strconv.Itoa(time.Now().Year()),
	}
	t.ExecuteTemplate(w, "Tenants", tenantsTemplate)
}

func rentRollTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrolltemplatehandler - begin")
	conf := configuration()
	conf.GPlusSigninCallback = "gRentRollTemplate"
	conf.FacebookSigninCallback = "fbRentRollTemplate"
	rentroll := RentRoll{Conf: conf}
	t, _ := template.ParseFiles(
		"rentrolltemplate.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	log.Print("rentrollhandler - execute")
	t.Execute(w, rentroll)
}
