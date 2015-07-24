package main
import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strconv"
	"time"
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
		AddTenant(r)
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
	Id int
	DbName string
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
		logError("rentroll - NO EMAIL SET")
		tenants = []Tenant{Tenant{
			Id: -1,
			DbName: "",
			Name: "#1",
			Address: "",
			SqFt: 0,
			LeaseStartDate: "",
			LeaseEndDate: "",
			BaseRent: "0",
			Electricity: "0",
			Gas: "0",
			Water: "0",
			SewageTrashRecycle: "0",
			Comments: ""}}
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
		logError(fmt.Sprintf("formatMoney - can't parse money: %v", mon))
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

func addTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("addTenantHandler - begin")
	tenantId, _ := AddTenant(r,)
	w.Write([]byte(strconv.FormatInt(tenantId, 10)))
}

func AddTenant(r *http.Request) (int64, bool) {
	sqFt, _ := strconv.Atoi(r.FormValue("SqFt"))
	return addTenant(r.FormValue("DbName"),
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

func removeTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("removeTenantHandler - begin")
	tenantAction(w, r, dbRemoveTenant)
}

func undoRemoveTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("undoRemoveTenantHandler - begin")
	tenantAction(w, r, dbUndoRemoveTenant)
}

func tenantAction(w http.ResponseWriter, r *http.Request, action func(db string, id int) bool) {
	log.Print("tenantAction - begin")
	success := true
	dbName := ""

	if dbName = r.FormValue("DbName"); dbName == "" {
		logError("Blank DbName")
		success = false
	} else if tenantId, err := strconv.Atoi(r.FormValue("TenantId")); err != nil {
		logError(fmt.Sprintf("Bad TenantId: %v", err))
		success = false
	} else {
		success = action(dbName, tenantId)
	}
	w.Write([]byte(strconv.FormatBool(success)))
}

func addTenant(dbName, name, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) (int64, bool) {
	if name == "" {
		logError("addtenant - NO NAME SET")
		return  -1, false
	}
	log.Print("addtenant - name")
	log.Print(name)
	if dbName == "" {
		logError("addtenant - NO DBNAME SET")
		return -1, false
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
	tenantId, success := dbInsert(dbName, name, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if !success {
		logError("Add tenant, error calling dbInsert")
	}
	return tenantId, success
}

func removeTenant(dbName string, tenantId int) bool {
	return dbRemoveTenant(dbName, tenantId)
}

func undoRemoveTenant(dbName string, tenantId int) bool {
	return dbUndoRemoveTenant(dbName, tenantId)
}
