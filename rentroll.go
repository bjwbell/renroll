package main
import (
	"net/http"
	"html/template"
)

type RentRoll struct {
	GoogleClientID string
	GoogleAnalyticsId string
	Tenants []Tenant	
}

type Tenant struct {
	Name string
}

func rentRollHandler(w http.ResponseWriter, r *http.Request) {
	dbName := r.FormValue("email")
	tenants := dbReadTenants(dbName)
	conf := configuration()
	
	rentroll := RentRoll{GoogleAnalyticsId: conf.GoogleAnalyticsId,
		GoogleClientID: conf.GoogleClientID,
		Tenants: tenants}
	
	t, _ := template.ParseFiles("rentroll.html")
	t.Execute(w, rentroll)
}

