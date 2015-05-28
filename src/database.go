package main
import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func dbCreate(name string) {
	os.Remove("./" + name + ".sqlite")

	db, err := sql.Open("sqlite3", "./" + name + ".sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table tenants (id integer not null primary key, name text);
	delete from tenants;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func dbInsert(databaseName, tenantName string) {
	db, err := sql.Open("sqlite3", "./" + databaseName + ".sqlite")
	if err != nil {
		log.Print("dbInsert - FATAL ERROR:")
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Print("dbInsert - FATAL ERROR:")
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into tenants(id, name) values(?, ?)")
	if err != nil {
		log.Print("dbInsert - FATAL ERROR:")
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(nil, tenantName)
	if err != nil {
		log.Print("dbInsert - FATAL ERROR:")
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
		dbCreate(databaseName);
	}
	db, err := sql.Open("sqlite3", "./" + databaseName + ".sqlite")
	if err != nil {
		log.Print("dbReadTenants - FATAL ERROR:")
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, name from tenants")
	if err != nil {
		log.Print("dbReadTenants - FATAL ERROR:")
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
