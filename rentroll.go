package main
import (
	"log"
	"net/http"
	"html/template"
)

type RentRoll struct {
	Conf Configuration
	Tenants []Tenant	
}

type Tenant struct {
	Name string
}

func rentRollHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrollhandler")
	email := r.FormValue("email")
	tenants := []Tenant{}
	if email == "" {
		log.Print("rentroll - NO EMAIL SET")
		tenants = []Tenant{Tenant{"#1"}}
	} else {
		dbName := r.FormValue("email")
		tenants = dbReadTenants(dbName)
	}
	conf := configuration()
	rentroll := RentRoll{Conf: conf, Tenants: tenants}
	t, _ := template.ParseFiles("rentroll.html")
	log.Print("rentrollhandler - execute")
	t.Execute(w, rentroll)
}

