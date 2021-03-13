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
	DB         *sql.DB
	log        *log.Logger
	prepareStr string
}

const modelsCodePrefix = "MOD00"

func init() {
	dbLoginUser := os.Getenv("DB_USER")
	dbLoginPassWord := os.Getenv("DB_PWD")
	// fmt.Printf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME"))
	// 這只是連結還沒開始連線，可以先用Ping做測試連線
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbLoginUser, dbLoginPassWord, os.Getenv("DB_HOSTNAME"), os.Getenv("DB_NAME")))
	err = db.Ping()
	if err != nil {
		fmt.Printf("DBConnect: %s", err.Error())
		panic(err.Error())
	}
	l := log.New(os.Stdout, "[sql]", log.LstdFlags)
	DBM = &DBManager{DB: db, log: l}

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
		panic(err.Error())
	}
	l := log.New(os.Stdout, "[sql]", log.LstdFlags)
	DBM = &DBManager{DB: db, log: l}
}

// SQLDebug 印出 sql statement
func (m *DBManager) SQLDebug(args ...interface{}) {
	m.log.Println(m.prepareStr)
	m.log.Println(args...)
}

// SetQuery 設定SQL語法
func (m *DBManager) SetQuery(query string) {
	m.prepareStr = query
}

// GetQuery 取得設定SQL語法
func (m *DBManager) GetQuery() string {
	return m.prepareStr
}

func (m *DBManager) Prepare() (*sql.Stmt, error) {
	return m.DB.Prepare(m.GetQuery())
}
