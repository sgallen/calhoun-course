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

type User struct {
	Id    int
	Name  string
	Email string
}

func main() {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	pingDatabase(psqlInfo)
	runTestQuery(psqlInfo)
	insertUserQuery(psqlInfo)
	selectUsersQuery(psqlInfo)
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

func insertUserQuery(connString string) {
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

func selectUsersQuery(connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(
		context.Background(),
		"SELECT id, name, email FROM users;",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to select data: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Name, &u.Email)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan rows: %v\n", err)
			os.Exit(1)
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Users:", users)
}
