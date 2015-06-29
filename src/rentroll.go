package main
import (
	"log"
	"net/http"
	"html/template"
	"time"
	"strconv"
	"github.com/joiggama/money"
)

type RentRoll struct {
	Conf Configuration
	AsOfDateDay string
	AsOfDateMonth string
	AsOfDateYear  string
	DefaultLeaseStartDate string
	DefaultLeaseEndDate string

}

func rentRollHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrollhandler - begin")
	conf := configuration()
	conf.GPlusSigninCallback = "gRentRoll"
	conf.FacebookSigninCallback = "fbRentRoll"
	month := time.Now().Month()
	day := strconv.Itoa(time.Now().Day())
	year := time.Now().Year()
	start := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year)
	end := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year + 3)
	rentroll := RentRoll{
		Conf: conf,
		AsOfDateDay: day,
		AsOfDateMonth: month.String(),
		AsOfDateYear: strconv.Itoa(time.Now().Year()),
		DefaultLeaseStartDate: start,
		DefaultLeaseEndDate: end,
	}
	if r.FormValue("Name") != "" {
		sqFt, _ := strconv.Atoi(r.FormValue("SqFt"))
		addTenant(r.FormValue("DbName"),
			r.FormValue("Name"),
			r.FormValue("Address"),
			sqFt,
			r.FormValue("LeaseStartDate"),
			r.FormValue("LeaseEndDate"),
			r.FormValue("BaseRent"),
			r.FormValue("Electricity"),
			r.FormValue("Gas"),
			r.FormValue("Water"),
			r.FormValue("SewageTrashRecycle"),
			r.FormValue("Comments"))
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
		formatCurrency(tenants)
	}
	t, _ := template.ParseFiles("templates/tenants.html")
	log.Print("tenanthandler - execute")
	tenantsTemplate := TenantsTemplate{
		Conf: configuration(),
		Tenants: tenants,
	}
	t.ExecuteTemplate(w, "Tenants", tenantsTemplate)
}

func formatCurrency(tenants []Tenant) {
	for i, _ := range tenants {
		tenants[i].BaseRent = formatMoney(tenants[i].BaseRent)
		tenants[i].Electricity = formatMoney(tenants[i].Electricity)
		tenants[i].Gas = formatMoney(tenants[i].Gas)
		tenants[i].Water = formatMoney(tenants[i].Water)
		tenants[i].SewageTrashRecycle = formatMoney(tenants[i].SewageTrashRecycle)
	}
}

func formatMoney(mon string) string {
	if mon == "" {
		mon = "0"
	}
	val, err := strconv.ParseFloat(mon, 32)
	if err != nil {
		log.Print("formatMoney - can't parse money: ")
		log.Print(mon)
		return ""
	}
	return money.New(val)
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

func addTenant(dbName, name, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) {
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
	dbInsert(dbName, name, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)	
}
