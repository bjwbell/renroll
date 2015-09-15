package renroll
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/joiggama/money"
	"github.com/jung-kurt/gofpdf"
)

type RentRoll struct {
	Conf Configuration
	AsOfDateDay string
	AsOfDateMonth string
	AsOfDateYear  string
	DefaultLeaseStartDate string
	DefaultLeaseEndDate string

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
	Total string
	Comments string
}

type ByTenantId []Tenant

func (t ByTenantId) Len() int           { return len(t) }
func (t ByTenantId) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTenantId) Less(i, j int) bool { return t[i].Id < t[j].Id }

func TenantsHandler(w http.ResponseWriter, r *http.Request) {
	tenants := map[int]Tenant{}
	log.Print("tenantshandler - begin")
	email := r.FormValue("email")
	if email == "" || email == "dummy@dummy.com" {
		logError("rentroll - NO EMAIL SET")
		tenants = map[int]Tenant{-1: Tenant{
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
			Total: "0",
			Comments: ""}}
	} else {
		dbName := email
		tenants = dbReadTenants(dbName)
		formatCurrency(tenants)
	}
	t, _ := template.ParseFiles("templates/tenants.html")
	log.Print("tenanthandler - execute")
	tenants1 := []Tenant{}
	for _, tenant := range tenants {
		tenants1 = append(tenants1, tenant)
	}
	sort.Sort(ByTenantId(tenants1))
	tenantsTemplate := TenantsTemplate{
		Conf: Config(),
		Tenants: tenants1,
	}
	t.ExecuteTemplate(w, "Tenants", tenantsTemplate)
}

func formatCurrency(tenants map[int]Tenant) {
	for id, _ := range tenants {
		tenant := tenants[id]
		tenant.Total = formatMoney(
			strconv.FormatFloat(parseMoney(tenant.BaseRent) +
				parseMoney(tenant.Electricity) +
				parseMoney(tenant.Gas) +
				parseMoney(tenant.Water) +
				parseMoney(tenant.SewageTrashRecycle),'f', -1, 64))
		
		tenant.BaseRent = formatMoney(tenant.BaseRent)
		tenant.Electricity = formatMoney(tenant.Electricity)
		tenant.Gas = formatMoney(tenant.Gas)
		tenant.Water = formatMoney(tenant.Water)
		tenant.SewageTrashRecycle = formatMoney(tenant.SewageTrashRecycle)
		
		tenants[id] = tenant
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

func UpdateTenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("updateTenantHandler - begin")
	success := UpdateTenant(r,)
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
		return  false
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

func PrintInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("printInvoicesHandler - begin")
	pdf := gofpdf.New("P", "in", "Letter", "")
	pdf.SetFont("Arial", "", 12)
	header := []string{"Base Rent", "Electricity", "Gas", "Water", "Sewage/Trash/Rec.", "Total"}
	widths := []float64{1.0, 1.0, 1.0, 1.0, 1.5, 1.0}
	dbName := r.FormValue("DbName")
	if dbName == "" {
		logError("printinvoices - no dbname set")
		return
	}
	tenants := dbReadTenants(dbName)
	formatCurrency(tenants)
	for _, tenant := range tenants {
		pdf.AddPage()
		h := 0.4
		pdf.CellFormat(5, h, "Name: " + tenant.Name + " (#" + strconv.Itoa(tenant.Id) + ")", "", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.Ln(-1)
		for i, str := range header {
			pdf.CellFormat(widths[i], h, str, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
		pdf.CellFormat(widths[0], h, tenant.BaseRent, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[1], h, tenant.Electricity, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[2], h, tenant.Gas, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[3], h, tenant.Water, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[4], h, tenant.SewageTrashRecycle, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[5], h, tenant.Total, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}
	err := pdf.OutputFileAndClose("invoices.pdf")
	if err != nil {
		logError("Unable to print invoices - write file")
		return
	}
	invoices, err := ioutil.ReadFile("invoices.pdf")
	if err != nil {
		logError("Unable to print invoices - read file")
		return
	}
	log.Print("printInvoicesTenantHandler - end")
	w.Write(invoices)
}
