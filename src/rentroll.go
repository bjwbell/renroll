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
	if r.FormValue("Name") != "" {
		addTenant(r.FormValue("DbName"), r.FormValue("Name"))
	}
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
	DefaultLeaseStartDate string
	DefaultLeaseEndDate string
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
	month := time.Now().Month()
	day := strconv.Itoa(time.Now().Day())
	year := time.Now().Year()
	start := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year)
	end := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year + 3)
	tenantsTemplate := TenantsTemplate{
		Conf: configuration(),
		Tenants: tenants,
		AsOfDateDay: day,
		AsOfDateMonth: month.String(),
		AsOfDateYear: strconv.Itoa(time.Now().Year()),
		DefaultLeaseStartDate: start,
		DefaultLeaseEndDate: end,
	}
	t.ExecuteTemplate(w, "Tenants", tenantsTemplate)
}

func rentRollTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrolltemplate - begin")
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

func addTenant(dbName string, name string) {
	if name == "" {
		log.Print("addtenant - NO NAME SET")
		return
	}
	log.Print("addtenant - name")
	log.Print(name)
	if dbName == "" {
		log.Print("addtenant - NO DBNAME SET")
		return
	}
	log.Print("addtenant - dbname")
	log.Print(dbName)

	/*address := r.FormValue("address")
	sqft := r.FormValue("sqft")
	start := r.FormValue("leasestartdate")
	end := r.FormValue("leaseenddate")
	base := r.FormValue("baserent")
	electricity := r.FormValue("electricity")
	gas := r.FormValue("gas")
	water := r.FormValue("water")
	sewagetrashrecycle := r.FormValue("sewagetrashrecycle")
	comments := r.FormValue("comments")*/
	log.Print("addtenant - execute")
	dbInsert(dbName, name)
	
}
