package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// variabile globale
var DB *sql.DB

// funzione per connettere Go a PostgreSQL
func connectDB() {
	var err error

	connStr := "host=localhost port=5432 user=postgres password=root dbname=gestionale_db sslmode=disable"

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Errore apertura database:", err)
		return
	}

	err = DB.Ping()

	if err != nil {
		fmt.Println("Errore connessione database:", err)
		return
	}

	fmt.Println("Connessione a PostgreSQL riuscita.")
}
