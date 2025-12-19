package db

import (
	"database/sql"
	// "fmt"
	"log"
	_ "github.com/denisenkom/go-mssqldb" // MSSQL driver
)

var DB *sql.DB

func InitDB() {
	
	// connString := fmt.Sprintf(
    //     "server=DESKTOP-TT6NHQT;user id=sa;password=YourStrongPass;port=1433;database=erp;encrypt=disable",
    // )
	connString := "server=222.255.144.197;user id=dms;password=Khcn248@@;port=56712;database=DMS;encrypt=true;TrustServerCertificate=true"
	var err error
	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging DB: ", err.Error())
	}
	log.Println("Connected to SQL Server!")
	
	
}