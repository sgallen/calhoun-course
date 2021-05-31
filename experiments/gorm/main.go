package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/google/uuid"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "foo"
	password = "foo"
	dbname   = "usegolang_dev"
)

type User struct {
	// gorm.Model
	// To user uuid_generate_v4(), need to run the following in
	// postgres: CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name string
	// Note:
	// gorm won't alter the column and add "not null" if it was
	// previously created without that constraint.
	Email string `gorm:"not null;uniqueIndex"`
}

func main() {
	verboseSql := flag.Bool("v", false, "verbose sql logging")
	flag.Parse()

	logLevel := logger.Silent
	if *verboseSql {
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	// db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})

	fmt.Println("Success.")
	name, email := getUserInfo()
	user := User{Name: name, Email: email}
	result := db.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	fmt.Println(user)
}

func getUserInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Name?")
	name, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	name = strings.TrimSpace(name)

	fmt.Println("Email?")
	email, err = reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	email = strings.TrimSpace(email)

	return name, email
}
