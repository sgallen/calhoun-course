package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	// Had to: go get "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "foo"
	password = "foo"
	dbname   = "usegolang_dev"
)

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	pingDatabase(psqlInfo)
	runTestQuery(psqlInfo)
	runInsertQuery(psqlInfo)
}

func pingDatabase(connString string) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open failed: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successful ping!")
}

func runTestQuery(connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func runInsertQuery(connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var userId int
	row := conn.QueryRow(
		context.Background(),
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		"John Doe",
		"john@example.com",
	)
	err = row.Scan(&userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Inserted user and their id is:", userId)
}
