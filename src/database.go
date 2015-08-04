package main
import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sort"
	"strconv"
	"time"
)

const ActionInsert = "insert"
const ActionUpdate = "update"
const ActionUndoUpdate = "undoupdate"
const ActionRemove = "remove"
const ActionUndoRemove = "undoremove"

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func dbExists(name string) bool {
	success, _ := exists("./" + name + ".sqlite")
	if !success {
		logError("dbExists: database (" + name + ") doesnt exist")
	}
	return success
}

func dbCreate(name string) bool {
	if dbExists(name) {
		logError("Database (" + name + ") already exists, RECREATING")
		os.Remove("./" + name + ".sqlite")
	}
	db, err := sql.Open("sqlite3", "./" + name + ".sqlite")
	if err != nil {
		logError(fmt.Sprintf("Couldn't create database (" +
			name + "), ERROR: %v", err))
		log.Fatal(err)
		return false
	}
	defer db.Close()
	sqlStmt := `
	create table tenants
(id integer not null primary key,
Action text, ActionTenantId integer, ActionRowId integer, ActionTimeStamp text, 
Name text, Address text, SqFt integer,
LeaseStartDate text, LeaseEndDate text, 
BaseRent text, Electricity text, Gas text, Water text, SewageTrashRecycle text,
Comments text);
	delete from tenants;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logError(fmt.Sprintf("Couldn't create table, database (" +
			name + "), ERROR (%q: %s\n)", err, sqlStmt))
		return false
	}
	return true
}

func dbInsert(dbName, tenantName, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) (int64, bool) {

	if !dbExists(dbName) {
		return -1, false
	}
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't open database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return -1, false
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		logError("Couldn't exec begin for database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return -1, false
	}
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTimeStamp, Name, Address, SqFt, LeaseStartDate, LeaseEndDate, BaseRent, Electricity, Gas, Water, SewageTrashRecycle, Comments) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return -1, false
	}
	defer stmt.Close()
	var timestamp = time.Now()
	result, err := stmt.Exec(nil, ActionInsert, timestamp, tenantName, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if err != nil {
		logError("Couldn't exec insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return -1, false
	}
	if err = tx.Commit(); err != nil {
		logError("Couldn't exec insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return -1, false
	}
	id, _ := result.LastInsertId()
	return id, true
}


func dbUpdate(dbName string, tenantId int, tenantName, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) bool {

	if !dbExists(dbName) {
		return false
	}
	fmt.Println("dbUpdate - tenantId:")
	fmt.Println(tenantId)
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't open database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		logError("Couldn't exec begin for database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTenantId, ActionTimeStamp, Name, Address, SqFt, LeaseStartDate, LeaseEndDate, BaseRent, Electricity, Gas, Water, SewageTrashRecycle, Comments) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare update insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	defer stmt.Close()
	var timestamp = time.Now()
	_, err = stmt.Exec(nil, ActionUpdate, tenantId, timestamp, tenantName, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if err != nil {
		logError("Couldn't exec update insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	if err = tx.Commit(); err != nil {
		logError("Couldn't exec update insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	return true
}

func dbReadTenants(dbName string) map[int]Tenant {
	if !dbExists(dbName) {
		logError("dbReadTenants: CREATING Database (" + dbName + ")")
		dbCreate(dbName)
	}
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't read database (" + dbName + ")")
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`select
                               id,                               
                               Name, Address, SqFt,
                               LeaseStartDate, LeaseEndDate,
                               BaseRent, Electricity, Gas, Water, SewageTrashRecycle,
                               Comments  from tenants where Action='insert' and ActionTenantId is null`)
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	defer rows.Close()
	tenants1 := []Tenant{}
	for rows.Next() {
		var id, SqFt int
		var
		Name,
		Address,
		LeaseStartDate,
		LeaseEndDate,
		BaseRent,
		Electricity,
		Gas,
		Water,
		SewageTrashRecycle,
		Comments string

		rows.Scan(
			&id,
			&Name,
			&Address,
			&SqFt,
			&LeaseStartDate,
			&LeaseEndDate,
			&BaseRent,
			&Electricity,
			&Gas,
			&Water,
			&SewageTrashRecycle,
			&Comments)
		
		var tenant = Tenant{
			Id: id,
			DbName: dbName,
			Name: Name,
			Address: Address,
			SqFt: SqFt,
			LeaseStartDate: LeaseStartDate,
			LeaseEndDate: LeaseEndDate,
			BaseRent: BaseRent,
			Electricity: Electricity,
			Gas: Gas,
			Water: Water,
			SewageTrashRecycle: SewageTrashRecycle,
			Comments: Comments}
		tenants1 = append(tenants1, tenant)
	}
	rows.Close()
	removedIds := dbRemovedTenantIds(dbName)
	tenants := map[int]Tenant{}
	for _, tenant := range tenants1 {
		removed := false
		for _, tenantId := range removedIds {
			if tenantId == tenant.Id {
				removed = true
				break
			}
		}
		if !removed {
			tenants[tenant.Id] = tenant
			if new, err := dbUpdatedTenantValues(dbName, tenant.Id); err {
				tenants[tenant.Id] = new
			}
		} 
	}
	return tenants
}

func dbReadTenant(dbName string, tenantId int) Tenant {
	tenants := dbReadTenants(dbName)
	tenant, ok := tenants[tenantId]
	if ok != true {
		logError(fmt.Sprintf("Error getting tenantId: %v", tenantId))
	}
	return tenant
}

func dbRemovedTenantIds (dbName string) []int {
	if !dbExists(dbName) {
		logError("dbReadTenants: CREATING Database (" + dbName + ")")
		dbCreate(dbName)
	}
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't read database (" + dbName + ")")
		log.Fatal(err)
	}
	defer db.Close()
	tenants, err := db.Query(`select
                               id 
                               from tenants where
                               Action='` + ActionInsert + `'`)
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	defer tenants.Close()

	rows, err := db.Query(`select
                               ActionTenantId 
                               from tenants where
                               Action='` + ActionRemove + `' and ActionTenantId is not null`)
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	defer rows.Close()
	rows2, err := db.Query(`select
                               ActionTenantId 
                               from tenants where
                               Action='` + ActionUndoRemove + `' and ActionTenantId is not null`)
	defer rows2.Close()
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	removedIds := []int{}
	removesIds := []int{}
	undoIds := []int{}
	for rows.Next() {
		var id int
		rows.Scan(&id)
		removesIds = append(removesIds, id)
	}
	for rows2.Next() {
		var id int
		rows2.Scan(&id)
		undoIds = append(undoIds, id)
	}
	for tenants.Next() {
		var tenantId, removes, undo int
		tenants.Scan(&tenantId)
		for _, removedId := range removesIds {
			if removedId == tenantId {
				removes = removes + 1
			}
		}
		for _, undoId := range undoIds {
			if undoId == tenantId {
				undo = undo + 1
			}
		}
		if removes > undo {
			removedIds = append(removedIds, tenantId)
		}
	}
	return removedIds
}

func dbUpdatedTenantValues(dbName string, tenantId int) (Tenant, bool) {
	if !dbExists(dbName) {
		logError("dbReadTenants: CREATING Database (" + dbName + ")")
		dbCreate(dbName)
	}
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't read database (" + dbName + ")")
		log.Fatal(err)
	}
	defer db.Close()
	rows2, err := db.Query(`select
                               ActionRowId
                               from tenants where
                               Action='` + ActionUndoUpdate + `' and ActionTenantId is not null and ActionTenantId=` + strconv.Itoa(tenantId))
	defer rows2.Close()
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	undoIds := []int{}
	for rows2.Next() {
		var id int
		rows2.Scan(&id)
		undoIds = append(undoIds, id)
	}
	rows2.Close()
	rows, err := db.Query(`select
                               id,
                               Name, Address, SqFt,
                               LeaseStartDate, LeaseEndDate,
                               BaseRent, Electricity, Gas, Water, SewageTrashRecycle,
                               Comments  from tenants where Action='` + ActionUpdate + `' and ActionTenantId=` + strconv.Itoa(tenantId))
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	defer rows.Close()
	tenants1 := []Tenant{}
	for rows.Next() {
		var id int
		var SqFt int
		var
		Name,
		Address,
		LeaseStartDate,
		LeaseEndDate,
		BaseRent,
		Electricity,
		Gas,
		Water,
		SewageTrashRecycle,
		Comments string

		rows.Scan(
			&id,
			&Name,
			&Address,
			&SqFt,
			&LeaseStartDate,
			&LeaseEndDate,
			&BaseRent,
			&Electricity,
			&Gas,
			&Water,
			&SewageTrashRecycle,
			&Comments)
		
		var tenant = Tenant{
			Id: tenantId,
			DbName: dbName,
			Name: Name,
			Address: Address,
			SqFt: SqFt,
			LeaseStartDate: LeaseStartDate,
			LeaseEndDate: LeaseEndDate,
			BaseRent: BaseRent,
			Electricity: Electricity,
			Gas: Gas,
			Water: Water,
			SewageTrashRecycle: SewageTrashRecycle,
			Comments: Comments}
		searchIdx := sort.IntSlice(undoIds).Search(id)
		if searchIdx >= len(undoIds) || undoIds[searchIdx] != id {
			tenants1 = append(tenants1, tenant)
		}
	}

	
	if len(tenants1) > 0 {
		fmt.Println("Updated Tenant Values - tenantId:")
		fmt.Println(tenantId)
		fmt.Println("new values:")
		fmt.Println(tenants1[len(tenants1) - 1])
		return tenants1[len(tenants1) - 1], true
	} else {
		return Tenant{}, false
	}
}

func dbRemoveTenant(dbName string, tenantId int) bool {
	return dbTenantAction(dbName, ActionRemove, ActionInsert, tenantId)
}

func dbUndoRemoveTenant(dbName string, tenantId int) bool {
	log.Print("dbUndoRemoveTenant")
	return dbTenantAction(dbName, ActionUndoRemove, ActionRemove, tenantId)
}

func dbUndoUpdateTenant(dbName string, tenantId int) bool {
	log.Print("dbUndoRemoveTenant")
	return dbTenantAction(dbName, ActionUndoUpdate, ActionUpdate, tenantId)
}

func dbTenantAction(dbName string, action string, prevAction string, tenantId int) bool {
	if !dbExists(dbName) {
		return false
	}	
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't open database (" + dbName + ")" +
			", tenantId (" + strconv.Itoa(tenantId) + ")")
		return false
	}
	defer db.Close()

	rows2, err := db.Query(`select
                               id
                               from tenants where
                               Action='` + prevAction + `' and ((ActionTenantId is not null and ActionTenantId=` + strconv.Itoa(tenantId) + `) or (ActionTenantId is null and id=` + strconv.Itoa(tenantId) + `))`)
	defer rows2.Close()
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	var prevActionRowId int
	for rows2.Next() {
		rows2.Scan(&prevActionRowId)
	}
	rows2.Close()
	
	tx, err := db.Begin()
	if err != nil {
		logError("Couldn't exec begin for database (" + dbName + ")" +
			", tenantId (" + strconv.Itoa(tenantId) + ")")
		return false
	}
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTenantId, ActionRowId, ActionTimeStamp) values(?, ?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare remove tenant in database (" +
			dbName + ")" + ", tenantId (" +
			strconv.Itoa(tenantId) + ")")
		return false
	}
	defer stmt.Close()
	var timestamp = time.Now()
	_, err = stmt.Exec(nil, action, tenantId, prevActionRowId, timestamp)
	if err != nil {
		logError("Couldn't exec remove tenant in database (" +
			dbName + ")" + ", tenantId (" +
			strconv.Itoa(tenantId) + ")")
		log.Fatal(err)
		return false
	}
	if err = tx.Commit(); err != nil {
		logError("Couldn't exec remove tenant in database (" +
			dbName + ")" + ", tenantId (" +
			strconv.Itoa(tenantId) + ")")
		return false
	}
	return true
}
