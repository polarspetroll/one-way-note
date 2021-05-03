package main

import (
	"fmt"
	"net/http"
	"app"
	"os"
)

var (
	dbusr  = os.Getenv("DBUSR")
	dbpwd  = os.Getenv("DBPWD")
	dbname = os.Getenv("DBNAME")
	dbaddr = os.Getenv("DBADDR")
	lport = os.Getenv("LPORT")
)

func main() {
	app.DBc = fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", dbusr, dbpwd, dbaddr, dbname)
	app.Domain = os.Getenv("DOMAIN")
	http.HandleFunc("/", app.Index)
	http.HandleFunc("/t/", app.GetText)
	http.ListenAndServe(":" + lport, nil)

}
