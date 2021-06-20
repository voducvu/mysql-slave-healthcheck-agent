package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

var (
	dsn string
	Version = "0.0.2"
	replicationLagLimit = 5
)

func main() {
	var port int
	var showVersion bool
	flag.IntVar(&port, "port", 5000, "http listen port number")
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:3306)/?charset=utf8", "MySQL DSN")
	flag.IntVar(&replicationLagLimit, "lag-limit", 5, "replication lag limit")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("version %s\n", Version)
		return
	}

	log.Printf("Listing port %d", port)
	log.Printf("dsn %s", dsn)
	log.Printf("lag-limit %ds", replicationLagLimit)

	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", dsn)
	defer db.Close()

	if err != nil {
		serverError(w, err)
		return
	}
	rows, err := db.Query("SHOW SLAVE STATUS")
	if err != nil {
		serverError(w, err)
		return
	}
	if !rows.Next() {
		serverError(w, errors.New("No slave status"))
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	for i, _ := range values {
		var v sql.RawBytes
		values[i] = &v
	}

	err = rows.Scan(values...)
	if err != nil {
		serverError(w, err)
		return
	}

	slaveInfo := make(map[string]interface{})
	for i, name := range columns {
		bp := values[i].(*sql.RawBytes)
		vs := string(*bp)
		vi, err := strconv.ParseInt(vs, 10, 64)
		if err != nil {
			slaveInfo[name] = vs
		} else {
			slaveInfo[name] = vi
		}
	}

	if (slaveInfo["Slave_IO_Running"] != "Yes" || slaveInfo["Slave_SQL_Running"] != "Yes") {
		serverError(w, errors.New("Slave is not running."))
		return
	}

	if (int(slaveInfo["Seconds_Behind_Master"].(int64))  > replicationLagLimit) {
		error := errors.New(fmt.Sprintf("Replication lag exceeds the limit of %d", replicationLagLimit))
		serverError(w, error)
		return
	}

	w.Write([]byte("OK"))
}

func serverError(w http.ResponseWriter, err error) {
	log.Printf("error: %s", err)
	code := http.StatusInternalServerError
	http.Error(w, fmt.Sprintf("%s", err), code)
}
