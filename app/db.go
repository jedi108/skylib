package app

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var (
	DB  *sql.DB
	err error
)

func openConnection() {
	config_db := GetConfig("database")
	DB, err = sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", config_db["username"], config_db["password"], config_db["dbname"]))
	if err != nil {
		panic("failed to connect database:\n" + err.Error())
	}
	//defer DB.Close()
}

func GetConnection() {
	if DB == nil {
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