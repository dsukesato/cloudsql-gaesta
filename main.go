package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	db = DB()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/user/", userHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func DB() *sql.DB {
	var (
		connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
		user           = mustGetenv("CLOUDSQL_USER")
		dbName         = os.Getenv("CLOUDSQL_DATABASE_NAME") // NOTE: dbName may be empty
		password       = os.Getenv("CLOUDSQL_PASSWORD")      // NOTE: password may be empty
		socket         = os.Getenv("CLOUDSQL_SOCKET_PREFIX")
	)

	/*var user = "root"
	var password = "test"
	var dbName = "pbl_test"*/

	// /cloudsql is used on App Engine.
	if socket == "" {
		socket = "/cloudsql"
	}

	// MySQL Connection, comment out to use PostgreSQL.
	// connection string format: USER:PASSWORD@unix(/cloudsql/PROJECT_ID:REGION_ID:INSTANCE_ID)/[DB_NAME]

	/*dbURI := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", user, password, dbName)
	conn, err := sql.Open("mysql", dbURI)

	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}

	return conn*/

	dbURI := fmt.Sprintf("%s:%s@unix(%s/%s)/%s?parseTime=true", user, password, socket, connectionName, dbName)
	conn, err := sql.Open("mysql", dbURI)

	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}

	return conn

	// PostgreSQL Connection, uncomment to use.
	// connection string format: user=USER password=PASSWORD host=/cloudsql/PROJECT_ID:REGION_ID:INSTANCE_ID/[ dbname=DB_NAME]
	// dbURI := fmt.Sprintf("user=%s password=%s host=/cloudsql/%s dbname=%s", user, password, connectionName, dbName)
	// conn, err := sql.Open("postgres", dbURI)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	rows, err := db.Query("show databases")
	if err != nil {
		log.Printf("Could not query db: %v", err)
		http.Error(w, "Internal Error", 500)
		return
	}
	defer rows.Close()

	buf := bytes.NewBufferString("Databases:\n")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			log.Printf("Could not scan result: %v", err)
			http.Error(w, "Internal Error", 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}
	w.Write(buf.Bytes())
}

type User struct {
	id        int
	name      string
	password  string
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	rows, err := db.Query("select * from user where id = ?", 1)
	if err != nil {
		log.Printf("Could not query db: %v", err)
		http.Error(w, "Internal Error", 500)
		return
	}
	defer rows.Close()

	buf := bytes.NewBufferString("USER:\n")
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.id, &user.name, &user.password, &user.createdAt, &user.updatedAt, &user.deletedAt); err != nil {
			log.Printf("Could not scan result: %v", err)
			http.Error(w, "Internal Error", 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", user.name)
	}
	w.Write(buf.Bytes())
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
