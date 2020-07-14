package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"

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
		if err := initDB(initialInfo); err != nil {
			exit(err.Error())
		}
	}

	psqlInfo := fmt.Sprintf("%s dbname=%s", initialInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		exit(fmt.Sprintf("opening DB connection: %v", err))
	}
	defer db.Close()

	if err = createPhoneNumbersTable(db); err != nil {
		exit(fmt.Sprintf("createPhoneNumbersTable: %v", err))
	}

	fmt.Println("\nSeeding phone_numbers table...")
	for _, num := range phoneNums {
		phone, err := create(db, num)
		if err != nil {
			exit(fmt.Sprintf("create: %v", err))
		}
		fmt.Printf("%+v\n", phone)
	}

	fmt.Println("\ncreate violates UNIQUE constraint demo...")
	_, err = create(db, phoneNums[0])
	if err != nil {
		// Print UNIQUE violation error message for demonstration, then continue
		fmt.Fprintf(os.Stderr, "\tcreate: %v\n", err)
	}

	fmt.Printf("\nDeleting %s...\n", phoneNums[0])
	if err = deleteByNumber(db, phoneNums[0]); err != nil {
		exit(fmt.Sprintf("delete: %v\n", err))
	}

	fmt.Println("\nRecreating record without violating UNIQUE constraint")
	num, err := create(db, phoneNums[0])
	if err != nil {
		exit(fmt.Sprintf("create: %v", err))
	}
	fmt.Printf("%+v\n", num)

	fmt.Println("\nfindAll demo")
	foundNums, err := findAll(db)
	if err != nil {
		exit(fmt.Sprintf("findAllPhoneNumbers: %v", foundNums))
	}
	for _, num := range foundNums {
		fmt.Printf("%+v\n", num)
	}

	fmt.Println("\nfindByNumber demo")
	foundNum, err := findByNumber(db, "123 456 7891")
	if err != nil {
		exit(fmt.Sprintf("findByNumber: %v", err))
	}
	fmt.Printf("found: %+v\n", foundNum)

	fmt.Println("\nupdate demo")
	foundNum, err = update(db, foundNum, "123 456 7897")
	if err != nil {
		exit(fmt.Sprintf("update: %v", err))
	}
	fmt.Printf("updated: %+v\n", foundNum)
}

func initDB(psqlInfo string) error {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("opening DB connection: %v", err)
	}
	if err = resetDB(db, dbname); err != nil {
		return fmt.Errorf("creating DB: %v", err)
	}
	db.Close()
	return nil
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createPhoneNumbersTable(db *sql.DB) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255) UNIQUE
		)`
	_, err := db.Exec(stmt)
	return err
}

type phoneNumber struct {
	id    int
	value string
}

func create(db *sql.DB, number string) (*phoneNumber, error) {
	norm := normalize(number)
	stmt := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id, value`
	var num phoneNumber
	if err := db.QueryRow(stmt, norm).Scan(&num.id, &num.value); err != nil {
		return nil, err
	}
	return &num, nil
}

func find(db *sql.DB, id int) (string, error) {
	stmt := `SELECT value FROM phone_numbers WHERE id = $1`
	var number string
	if err := db.QueryRow(stmt, id).Scan(&number); err != nil {
		return "", err
	}
	return number, nil
}

func findByNumber(db *sql.DB, num string) (*phoneNumber, error) {
	norm := normalize(num)
	stmt := `SELECT id, value FROM phone_numbers WHERE value = $1`
	var result phoneNumber
	if err := db.QueryRow(stmt, norm).Scan(&result.id, &result.value); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func update(db *sql.DB, num *phoneNumber, newValue string) (*phoneNumber, error) {
	normal := normalize(newValue)
	stmt := `UPDATE phone_numbers SET value=$1 WHERE id=$2 RETURNING id, value`
	var updatedNum phoneNumber
	err := db.QueryRow(stmt, normal, num.id).Scan(&updatedNum.id, &updatedNum.value)
	if err != nil {
		return nil, err
	}
	return &updatedNum, nil
}

func deleteByNumber(db *sql.DB, num string) error {
	normal := normalize(num)
	stmt := `DELETE FROM phone_numbers WHERE value = $1`
	_, err := db.Exec(stmt, normal)
	return err
}

func findAll(db *sql.DB) ([]phoneNumber, error) {
	rows, err := db.Query(`SELECT id, value FROM phone_numbers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nums []phoneNumber
	for rows.Next() {
		var num phoneNumber
		if err := rows.Scan(&num.id, &num.value); err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nums, nil
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}
