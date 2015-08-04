package main
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strconv"
)

type TenantTemplate struct {
	Conf Configuration

}

func tenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenantHandler - begin")
	conf := configuration()
	conf.GPlusSigninCallback = "gTenant"
	conf.FacebookSigninCallback = "fbTenant"
	tenant := TenantTemplate {
		Conf: conf,
	};
	t, _ := template.ParseFiles(
		"tenant.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	log.Print("tenantHandler - execute")
	t.Execute(w, tenant)
}

func tenantDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenantDataHandler - begin")
	dbName := r.FormValue("DbName")
	tenantId, _ := strconv.Atoi(r.FormValue("TenantId"))
	tenant := dbReadTenant(dbName, tenantId)
	bytes, err := json.Marshal(tenant)
	if err != nil {
		logError(fmt.Sprintf("Error serializing tenant to json, ERR: %v", err))
	}
	w.Write(bytes)
}

func tenantsDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenanstDataHandler - begin")
	dbName := r.FormValue("DbName")
	tenants := map[string]Tenant{ }
	for tenantId, tenant := range dbReadTenants(dbName) {
		tenants[strconv.Itoa(tenantId)] = tenant
	}
	bytes, err := json.Marshal(tenants)
	if err != nil {
		logError(fmt.Sprintf("Error serializing tenants to json, ERR: %v", err))
	}
	w.Write(bytes)
}
