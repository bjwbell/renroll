package renroll

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joiggama/money"
)

// the rentroll struct holds...
type RentRoll struct {
	Conf                  Configuration
	AsOfDateDay           string
	AsOfDateMonth         string
	AsOfDateYear          string
	DefaultLeaseStartDate string
	DefaultLeaseEndDate   string
}

func RentRollHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrollhandler - begin")
	conf := Config()
	conf.GPlusSigninCallback = "gRentRoll"
	conf.FacebookSigninCallback = "fbRentRoll"
	month := time.Now().Month()
	day := strconv.Itoa(time.Now().Day())
	year := time.Now().Year()
	start := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year)
	end := strconv.Itoa(int(month)) + "/" + day + "/" + strconv.Itoa(year+3)

	rentroll := RentRoll{
		Conf:                  conf,
		AsOfDateDay:           day,
		AsOfDateMonth:         month.String(),
		AsOfDateYear:          strconv.Itoa(year),
		DefaultLeaseStartDate: start,
		DefaultLeaseEndDate:   end,
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
	Conf    Configuration
	Tenants []Tenant
	Day     string
	Month   string
	Year    string
}

type Tenant struct {
	Id                 int
	DbName             string
	Name               string
	Address            string
	SqFt               int
	LeaseStartDate     string
	LeaseEndDate       string
	BaseRent           string
	Electricity        string
	Gas                string
	Water              string
	SewageTrashRecycle string
	Total              string
	Comments           string
}

// The tenants handler
func TenantsHandler(w http.ResponseWriter, r *http.Request) {
	tenants := []Tenant{}
	log.Print("tenantshandler - begin")
	email := r.FormValue("email")
	if email == "" || email == "dummy@dummy.com" {
		logError("rentroll - NO EMAIL SET")
		tenants = []Tenant{Tenant{
			Id:                 -1,
			DbName:             "",
			Name:               "#1",
			Address:            "",
			SqFt:               0,
			LeaseStartDate:     "",
			LeaseEndDate:       "",
			BaseRent:           "0",
			Electricity:        "0",
			Gas:                "0",
			Water:              "0",
			SewageTrashRecycle: "0",
			Total:              "0",
			Comments:           ""}}
	} else {
		dbName := email
		tenants = dbReadSortedTenants(dbName)
		formatCurrency(tenants)
	}
	t, _ := template.ParseFiles("templates/tenants.html")
	log.Print("tenanthandler - execute")
	sort.Sort(ByTenantId(tenants))
	month := time.Now().Month()
	day := strconv.Itoa(time.Now().Day())
	year := time.Now().Year()
	tenantsTemplate := TenantsTemplate{
		Conf:    Config(),
		Tenants: tenants,
		Day:     day,
		Month:   month.String(),
		Year:    strconv.Itoa(year),
	}
	t.ExecuteTemplate(w, "Tenants", tenantsTemplate)
}

func formatCurrency(tenants []Tenant) {
	for idx, _ := range tenants {
		tenant := tenants[idx]
		tenant.Total = formatMoney(
			strconv.FormatFloat(parseMoney(tenant.BaseRent)+
				parseMoney(tenant.Electricity)+
				parseMoney(tenant.Gas)+
				parseMoney(tenant.Water)+
				parseMoney(tenant.SewageTrashRecycle), 'f', -1, 64))

		tenant.BaseRent = formatMoney(tenant.BaseRent)
		tenant.Electricity = formatMoney(tenant.Electricity)
		tenant.Gas = formatMoney(tenant.Gas)
		tenant.Water = formatMoney(tenant.Water)
		tenant.SewageTrashRecycle = formatMoney(tenant.SewageTrashRecycle)

		tenants[idx] = tenant
	}
}
func parseMoney(mon string) float64 {
	mon = strings.Replace(mon, "$", "", -1)
	mon = strings.Replace(mon, ",", "", -1)
	val, err := strconv.ParseFloat(mon, 64)
	if err != nil {
		logError(fmt.Sprintf("formatMoney - can't parse money: %v", mon))
		return -1
	}
	return val

}
func formatMoney(mon string) string {
	if mon == "" {
		mon = "0"
	}
	val, err := strconv.ParseFloat(mon, 64)
	if err != nil {
		logError(fmt.Sprintf("formatMoney - can't parse money: %v", mon))
		return ""
	}
	return money.New(val)
}

func RentRollTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("rentrolltemplate - begin")
	conf := Config()
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

func AddTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("addTenantHandler - begin")
	tenantId, _ := AddTenant(r)
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

func UpdateTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("updateTenantHandler - begin")
	success := UpdateTenant(r)
	w.Write([]byte(strconv.FormatBool(success)))
}

func UpdateTenant(r *http.Request) bool {
	sqFt, _ := strconv.Atoi(r.FormValue("SqFt"))
	tenantId, _ := strconv.Atoi(r.FormValue("TenantId"))
	return updateTenant(r.FormValue("DbName"),
		tenantId,
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

func RemoveTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("removeTenantHandler - begin")
	tenantAction(w, r, dbRemoveTenant)
}

func UndoRemoveTenantHandler(w http.ResponseWriter, r *http.Request) {
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
		return -1, false
	}
	log.Print("addtenant - name")
	log.Print(name)
	if dbName == "" {
		logError("addtenant - NO DBNAME SET")
		return -1, false
	}
	log.Print("addtenant - dbname")
	log.Print(dbName)

	log.Print("addtenant - execute")
	tenantId, success := dbInsert(dbName, name, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if !success {
		logError("Add tenant, error calling dbInsert")
	}
	return tenantId, success
}

func updateTenant(dbName string, tenantId int, name, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) bool {
	if name == "" {
		logError("addtenant - NO NAME SET")
		return false
	}
	if dbName == "" {
		logError("addtenant - NO DBNAME SET")
		return false
	}
	success := dbUpdate(dbName, tenantId, name, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if !success {
		logError("update tenant, error calling dbUpdate")
	}
	return success
}

func UndoUpdateTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("undoUpdateTenantHandler - begin")
	tenantAction(w, r, dbUndoUpdateTenant)
}
