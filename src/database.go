package main
import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

const ActionInsert = "insert"
const ActionModify = "modify"
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
Action text, ActionTenantId integer, ActionTimeStamp text, 
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

func dbInsert(dbName, tenantName, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) bool {

	if !dbExists(dbName) {
		return false
	}
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
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTimeStamp, Name, Address, SqFt, LeaseStartDate, LeaseEndDate, BaseRent, Electricity, Gas, Water, SewageTrashRecycle, Comments) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	defer stmt.Close()
	var timestamp = time.Now()
	_, err = stmt.Exec(nil, ActionInsert, timestamp, tenantName, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if err != nil {
		logError("Couldn't exec insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	if err = tx.Commit(); err != nil {
		logError("Couldn't exec insert in database (" + dbName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
		return false
	}
	return true
}

func dbReadTenants(dbName string) []Tenant {
	if !dbExists(dbName) {
		logError("dbReadTenants: CREATING Database (" + dbName + ")")
		dbCreate(dbName);
	}
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't read database (" + dbName + ")")
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(`select
                               id,
                               Action,
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
		var id, SqFt int;
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
		Comments string;

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
			&Comments);
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
			Comments: Comments};
		tenants1 = append(tenants1, tenant)
	}
	rows.Close()
	rows2, err := db.Query(`select
                               ActionTenantId 
                               from tenants where
                               Action='remove' and ActionTenantId is not null`)
	if err != nil {
		logError("Couldn't query database (" + dbName + ")")
		log.Fatal(err)
	}
	defer rows2.Close()
	removedIds := []int{}
	for rows2.Next() {
		var id int;
		rows2.Scan(&id)
		removedIds = append(removedIds, id)
	}
	rows2.Close()
	tenants := []Tenant{}
	for _, tenant := range tenants1 {
		removed := false
		for _, tenantId := range removedIds {
			if tenantId == tenant.Id {
				removed = true
				break
			}
		}
		if !removed {
			tenants = append(tenants, tenant)
		} 
	}
	return tenants
}

func dbRemoveTenant(dbName string, tenantId int) bool {
	return dbTenantAction(dbName, ActionRemove, tenantId)
}

func dbUndoRemoveTenant(dbName string, tenantId int) bool {
	return dbTenantAction(dbName, ActionUndoRemove, tenantId)
}

func dbTenantAction(dbName string, action string, tenantId int) bool {
	if !dbExists(dbName) {
		return false;
	}	
	db, err := sql.Open("sqlite3", "./" + dbName + ".sqlite")
	if err != nil {
		logError("Couldn't open database (" + dbName + ")" +
			", tenantId (" + strconv.Itoa(tenantId) + ")")
		return false
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		logError("Couldn't exec begin for database (" + dbName + ")" +
			", tenantId (" + strconv.Itoa(tenantId) + ")")
		return false
	}
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTenantId, ActionTimeStamp) values(?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare remove tenant in database (" +
			dbName + ")" + ", tenantId (" +
			strconv.Itoa(tenantId) + ")")
		return false
	}
	defer stmt.Close()
	var timestamp = time.Now()
	_, err = stmt.Exec(nil, ActionRemove, tenantId, timestamp)
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
