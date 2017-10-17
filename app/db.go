package app

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gocraft/dbr"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

var (
	TYPERUN string //PRODUCTION, TESTING
	DB      *sql.DB
	err     error
	DbrConn *dbr.Connection
)

func GetDB() *dbr.Connection {
	return DbrConn
}

func openConnection() {
	config_db := GetConfig("database")

	if config_db["protocol"] == "" {
		config_db["protocol"] = "tcp"
	}
	if config_db["host"] == "" {
		config_db["host"] = "127.0.0.1"
	}
	if config_db["port"] == "" {
		config_db["port"] = "3306"
	}
	netAddr := fmt.Sprintf("%s(%s:%s)", config_db["protocol"], config_db["host"], config_db["port"])
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s", config_db["username"], config_db["password"], netAddr, config_db["dbname"])
	DB, err = sql.Open("mysql", dsn)
	//defer DB.Close()
	if err != nil {
		panic("failed to connect database:\n" + err.Error())
	}

	DbrConn, err = dbr.Open("mysql", dsn, nil)

	//defer DB.Close()
	if err != nil {
		panic("failed to connect database:\n" + err.Error())
	}

	err = DbrConn.Ping()
	if err != nil {
		panic("failed to ping database:\n" + err.Error())
	}
}

func GetConnection() {
	if GetDB() == nil {
		//------------------------------------
		// FOR WORK WITH TEST DATABASE
		//
		// go run main.go exec -test
		//
		//------------------------------------
		boolPtr := flag.Bool("test", false, "a bool")
		flag.Parse()
		if *boolPtr == true {
			ConfigDevFileJson = ConfigTestsFileJson
			TYPERUN = "TESTING"
		}
		openConnection()
	}
}

func MakeInsertQuery(table string, myMap map[string]interface{}) (string, []interface{}) {
	keys := make([]string, 0, len(myMap))
	vals := make([]interface{}, 0, len(myMap))
	for key, val := range myMap {
		keys = append(keys, key)
		vals = append(vals, val)
	}
	placeholders := strings.Repeat("?,", len(keys))
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES(%s)",
		table,
		strings.Join(keys, ", "),
		placeholders[:len(placeholders)-1],
	), vals
}

func MakeUpdateQuery(table string, set map[string]interface{}, where map[string]interface{}) (string, []interface{}) {
	keys := make([]string, 0, len(set))
	vals := make([]interface{}, 0, len(set))
	wkeys := make([]string, 0, len(where))
	for key, val := range set {
		keys = append(keys, key+"=?")
		vals = append(vals, val)
	}

	for key, val := range where {
		wkeys = append(wkeys, key+"=?")
		vals = append(vals, val)
	}

	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(keys, ", "),
		strings.Join(wkeys, ", "),
	), vals
}
