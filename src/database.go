package main
import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const ActionInsert = "insert"
const ActionModify = "modify"
const ActionRemove = "remove"

func dbCreate(name string) {
	ex, _ := exists("./" + name + ".sqlite")
	if ex {
		logError("Database (" + name + ") already exists, RECREATING")
	}
	os.Remove("./" + name + ".sqlite")
	db, err := sql.Open("sqlite3", "./" + name + ".sqlite")
	if err != nil {
		logError(fmt.Sprintf("Couldn't create database (" + name + "), ERROR: %v", err))
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	create table tenants
(id integer not null primary key,
Action text,
ActionTimeStamp text,
ActionModifyRowId integer,
Name text,
Address text,
SqFt integer,
LeaseStartDate text,
LeaseEndDate text,
BaseRent text,
Electricity text,
Gas text,
Water text,
SewageTrashRecycle text,
Comments text);
	delete from tenants;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logError(fmt.Sprintf("Couldn't create table, database (" + name + "), ERROR (%q: %s\n)", err, sqlStmt))
		return
	}
}

func dbInsert(databaseName, tenantName, address string, sqft int, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments string) {
	ex, _ := exists("./" + databaseName + ".sqlite")
	if ex == false {
		logError("Database (" + databaseName + ") doesnt exist, CREATING")
		dbCreate(databaseName);
	}
	db, err := sql.Open("sqlite3", "./" + databaseName + ".sqlite")
	if err != nil {
		logError("Couldn't open database (" + databaseName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		logError("Couldn't exec begin for database (" + databaseName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into tenants(id, Action, ActionTimeStamp, Name, Address, SqFt, LeaseStartDate, LeaseEndDate, BaseRent, Electricity, Gas, Water, SewageTrashRecycle, Comments) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		logError("Couldn't prepare insert in database (" + databaseName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
	}
	defer stmt.Close()
	var timestamp = time.Now()
	_, err = stmt.Exec(nil, ActionInsert, timestamp, tenantName, address, sqft, start, end, baseRent, electricity, gas, water, sewageTrashRecycle, comments)
	if err != nil {
		logError("Couldn't exec insert in database (" + databaseName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
	}
	tx.Commit()
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func dbReadTenants(databaseName string) []Tenant {
	ex, _ := exists("./" + databaseName + ".sqlite")
	if ex == false {
		logError("Database (" + databaseName + ") doesnt exist, CREATING")
		dbCreate(databaseName);
	}
	db, err := sql.Open("sqlite3", "./" + databaseName + ".sqlite")
	if err != nil {
		logError("Couldn't read database (" + databaseName + ")")
		log.Fatal(err)
		
	}
	defer db.Close()

	rows, err := db.Query(`select
id,
Name,
Address,
SqFt,
LeaseStartDate,
LeaseEndDate,
BaseRent,
Electricity,
Gas,
Water,
SewageTrashRecycle,
Comments 
from tenants`)
	if err != nil {
		logError("Couldn't query database (" + databaseName + ")")
		log.Fatal(err)
	}
	defer rows.Close()
	tenants := []Tenant{}
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
		
		//fmt.Println(id, Name, )
		var tenant = Tenant{
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
		tenants = append(tenants, tenant)
	}
	rows.Close()
	return tenants
}

