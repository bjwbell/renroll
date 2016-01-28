package renroll

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func TenantHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenantHandler - begin")
	conf := Config()
	conf.GPlusSigninCallback = "gTenant"
	conf.FacebookSigninCallback = "fbTenant"
	tenant := struct{ Conf Configuration }{Config()}
	t, _ := template.ParseFiles(
		"tenant.html",
		"templates/header.html",
		"templates/topbar.html",
		"templates/bottombar.html")
	log.Print("tenantHandler - execute")
	t.Execute(w, tenant)
}

func TenantDataHandler(w http.ResponseWriter, r *http.Request) {
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

func TenantsDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenanstDataHandler - begin")
	dbName := r.FormValue("DbName")
	tenants := map[string]Tenant{}
	for tenantId, tenant := range dbReadTenants(dbName) {
		tenants[strconv.Itoa(tenantId)] = tenant
	}
	bytes, err := json.Marshal(tenants)
	if err != nil {
		logError(fmt.Sprintf("Error serializing tenants to json, ERR: %v", err))
	}
	w.Write(bytes)
}

func TenantHistoryHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("tenantHistoryHandler - begin")
	dbName := r.FormValue("DbName")
	tenantId, _ := strconv.Atoi(r.FormValue("TenantId"))
	actions := dbTenantHistory(dbName, tenantId)
	bytes, err := json.Marshal(actions)
	if err != nil {
		logError(fmt.Sprintf("Error serializing tenant history to json, ERR: %v", err))
	}
	w.Write(bytes)
}
