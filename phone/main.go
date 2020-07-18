package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/angusgmorrison/gophercises/phone/db"
	_ "github.com/lib/pq"
)

const (
	dbname = "gophercises_phone"
)

var (
	host      = os.Getenv("PGHOST")
	port      = os.Getenv("PGPORT")
	user      = os.Getenv("PGUSER2")
	password  = os.Getenv("PGPASSWORD2")
	phoneNums = []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
	}
)

func main() {
	reset := flag.Bool("reset", false, "reset the database on program start")
	flag.Parse()

	initialInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		host, port, user, password)
	if *reset {
		if err := db.InitDB("postgres", initialInfo, dbname); err != nil {
			exit(err.Error())
		}
	}

	psqlInfo := fmt.Sprintf("%s dbname=%s", initialInfo, dbname)
	phoneDB, err := db.Open("postgres", psqlInfo)
	if err != nil {
		exit(fmt.Sprintf("opening DB connection: %v", err))
	}
	defer phoneDB.Close()

	if err = phoneDB.CreateTable("phone_numbers"); err != nil {
		exit(fmt.Sprintf("createPhoneNumbersTable: %v", err))
	}

	fmt.Println("\nSeeding phone_numbers table...")
	for _, num := range phoneNums {
		phone, err := phoneDB.Create(num)
		if err != nil {
			exit(fmt.Sprintf("create: %v", err))
		}
		fmt.Printf("%+v\n", phone)
	}

	fmt.Println("\nDeleting by ID...")
	err = phoneDB.DeleteBy("id", 1)
	if err != nil {
		exit(fmt.Sprintf("phoneDB.DeleteBy(\"id\"): %v", err))
	}
	pNum, err := phoneDB.FindByID(1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", pNum)
	}

	fmt.Println("\nDeleting by Value...")
	err = phoneDB.DeleteBy("value", "123 456 7891")
	if err != nil {
		exit(fmt.Sprintf("db.DeleteBy(\"value\"): %v", err))
	}
	pNum, err = phoneDB.FindByValue("123 456 7891")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", pNum)
	}

	fmt.Println("\nRetrieving all records...")
	pNums, err := phoneDB.All()
	if err != nil {
		exit(fmt.Sprintf("db.All(): %v", err))
	}
	for _, pNum := range pNums {
		fmt.Printf("%+v\n", pNum)
	}
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}
