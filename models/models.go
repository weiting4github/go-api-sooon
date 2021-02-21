package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // for "database/sql"
	"log"
	"os"
)

// DBM singleton instance connection pool
var DBM *DBManager

// DBManager sqldb 管理器
type DBManager struct {
	DB *sql.DB
}

const modelsCodePrefix = "MOD00"

func init() {
	dbLoginUser := os.Getenv("DB_USER")
	dbLoginPassWord := os.Getenv("DB_PWD")
	// fmt.Printf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME"))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME")))
	err = db.Ping()
	if err != nil {
		fmt.Printf("DBConnect: %s", err.Error())
		panic(err.Error())
	}
	DBM = &DBManager{DB: db}

}

/*NewDBConnect 連線DB func 若DB掛了可以用這組create*/
func (m *DBManager) NewDBConnect() {
	dbLoginUser := os.Getenv("DB_USER")
	dbLoginPassWord := os.Getenv("DB_PWD")
	// fmt.Printf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME"))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal("Failed to init db:", err)
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("DBConnect: %s", err.Error())

	}
	DBM = &DBManager{DB: db}
}
