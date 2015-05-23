package main
import (
	"net/http"
	"html/template"
)

type RentRoll struct {
	GoogleClientID string
	Tenants []Tenant
}

type Tenant struct {
	Name string
}

func rentRollHandler(w http.ResponseWriter, r *http.Request) {
	dbName := r.FormValue("email")
	tenants := dbReadTenants(dbName)
	rentroll := RentRoll{GoogleClientID: configuration().GoogleClientID, Tenants: tenants}
	t, _ := template.ParseFiles("rentroll.html")
	t.Execute(w, rentroll)
}

