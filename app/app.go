package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
)
type tmpl struct {
  Stat bool
  Val string
  Text string
}
var DBc, Domain string
func Index(w http.ResponseWriter, r *http.Request) {
  var out tmpl
  temp, err := template.ParseFiles("template/index.html")
  ErrCheck(err)
  DB, err := sql.Open("mysql", DBc)
  ErrCheck(err)
  if r.Method == "GET" {
    out.Stat = true
    temp.Execute(w, out)
    return
  }
  r.ParseForm()
  text := r.PostForm.Get("text")
  if len(text) > 65520 {
    out.Text = "ERROR!text too long."
    temp.Execute(w, out)
    return
  }
  path := RandPath()
  if Insert(text, path, DB) > 1 {
    out.Text = "Server Error"
    temp.Execute(w, out)
    return
  }
  out.Text = "http://" + Domain + "/t/" + path
  temp.Execute(w, out)
}

func GetText(w http.ResponseWriter, r *http.Request) {
  var out tmpl
  temp, err := template.ParseFiles("template/index.html")
  ErrCheck(err)
  DB, err := sql.Open("mysql", DBc)
  ErrCheck(err)
  if r.Method == "POST" {
    http.Error(w, http.StatusText(405), 405)
    return
  }
  u := ParseURL(r.URL.Path)
  if u == "" {
    http.NotFound(w, r)
    return
  }
  text := Query(u, DB)
  if text == "" {
    http.NotFound(w, r)
    return
  }
  if Delete(u, DB) > 1 {
    http.Error(w, http.StatusText(500), 500)
    return
  }
  out.Val = template.HTMLEscapeString(text)
  temp.Execute(w, out)
}

func Query(path string, DB *sql.DB) string {
	q, err := DB.Query("SELECT text FROM ownote WHERE path=?", path)
	ErrCheck(err)
	if q.Next() == false {
		return ""
	}
	var text string
	q.Scan(&text)
	return text
}

func Delete(path string, DB *sql.DB) int64 {
	d, err := DB.Prepare("DELETE FROM ownote WHERE path=?")
	ErrCheck(err)
	e, err := d.Exec(path)
	ErrCheck(err)
	a, err := e.RowsAffected()
	ErrCheck(err)
	return a
}

func Insert(text string, path string, DB *sql.DB) int64 {
	i, err := DB.Prepare("INSERT INTO ownote(text, path) VALUES(?, ?)")
	ErrCheck(err)
	e, err := i.Exec(text, path)
	ErrCheck(err)
	a, err := e.RowsAffected()
	ErrCheck(err)
	return a
}

func RandPath() string {
	r := make([]byte, 10)
	rand.Read(r)
	return hex.EncodeToString(r)
}

func ParseURL(path string) string {
	if len(path) <= 3 {
		return ""
	}
	if string(path[len(path)-1]) == "/" {
		return string(path[3 : len(path)-1])
	} else {
		return string(path[3:])
	}
}

func ErrCheck(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
