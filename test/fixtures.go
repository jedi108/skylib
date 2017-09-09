package test

import (
	"fmt"
	"io/ioutil"
	"skylib/app"
	"strings"
)

func RunFixtures() {
	files, err := ioutil.ReadDir(DirFixtures)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		sqlByte, err := loadFile(DirFixtures + f.Name())
		CheckErr(err)
		splitSql(string(sqlByte))
		CheckErr(err)
	}
}

func loadFile(filename string) ([]byte, error) {
	if isVerboseFixtures {
		PrintLn(filename)
	}
	return ioutil.ReadFile(filename)
}

func splitSql(sql string) {
	s := strings.Split(sql, ";")
	for _, vv := range s {
		if len(vv) > 1 {
			runQuery(vv)
		}
	}
}

func runQuery(sql string) {
	var pointers string
	rows, err := app.DB.Query(sql)
	CheckErr(err)
	for rows.Next() {
		rows.Scan(&pointers)
		if isVerboseFixtures {
			PrintLn(pointers)
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		PrintLn("__________________________________")
		PrintLn(err.Error())
		T.Error(err)
		T.Fail()
	}
}

func PrintLn(stringUser string) {
	if &T != nil {
		fmt.Println(stringUser)
		return
	}
	T.Log(stringUser)
}
