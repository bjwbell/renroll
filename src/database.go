package main
import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

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
	create table tenants (id integer not null primary key, name text);
	delete from tenants;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		logError(fmt.Sprintf("Couldn't create table, database (" + name + "), ERROR (%q: %s\n)", err, sqlStmt))
		return
	}
}

func dbInsert(databaseName, tenantName string) {
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
	stmt, err := tx.Prepare("insert into tenants(id, name) values(?, ?)")
	if err != nil {
		logError("Couldn't prepare insert in database (" + databaseName + ")" +
			", tenant (" + tenantName + ")")
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(nil, tenantName)
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

	rows, err := db.Query("select id, name from tenants")
	if err != nil {
		logError("Couldn't query database (" + databaseName + ")")
		log.Fatal(err)
	}
	defer rows.Close()
	tenants := []Tenant{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
		tenants = append(tenants, Tenant{Name: name})
	}
	rows.Close()
	return tenants
}
