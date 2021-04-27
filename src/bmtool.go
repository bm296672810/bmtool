package main

import (
	"log"
	"os"
	"project/bmtool/src/mysql"
)

var ILogger = log.New(os.Stdout, "[INFO]", log.LstdFlags|log.Lshortfile)
var WLogger = log.New(os.Stdout, "[WARNING]", log.LstdFlags|log.Lshortfile)
var ELogger = log.New(os.Stdout, "[ERROR]", log.LstdFlags|log.Lshortfile)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 3 {
		ELogger.Print(`
		Plase use: bmtool ${sql_file} ${sql dirver name} ${package name}
		The dirver name is : mysql sqlite \n`)
		ELogger.Fatalln("cmd line args error")
	}

	sqlFile := argsWithoutProg[0]
	sqlDirver := argsWithoutProg[1]
	packageName := argsWithoutProg[2]
	// sqlFile := "../test.sql"
	// sqlDirver := "mysql"
	// packageName := "userdata"
	if sqlDirver == "mysql" {
		mysql.ProcessSql(sqlFile, packageName)
	}

}
