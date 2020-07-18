package db

import (
	"database/sql"
	"fmt"
)

type phoneDB struct {
	*sql.DB
}

func Open(driverName, dataSourceName string) (*phoneDB, error) {
	sqlDB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("opening DB connection: %v", err)
	}
	return &phoneDB{sqlDB}, nil
}

func InitDB(driverName, dataSourceName, dbName string) error {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("opening DB connection: %v", err)
	}
	if err = db.resetDB(dbName); err != nil {
		return fmt.Errorf("creating DB: %v", err)
	}
	return db.Close()
}

func (p *phoneDB) resetDB(dbName string) error {
	_, err := p.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return err
	}
	return p.createDB(dbName)
}

func (p *phoneDB) createDB(dbName string) error {
	_, err := p.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		return err
	}
	return nil
}

var tableName = "phone_numbers"

func TableName() string {
	return tableName
}

func (p *phoneDB) CreateTable(name string) error {
	if name != "" {
		tableName = name
	}
	stmt := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL,
			value VARCHAR(255) UNIQUE
		)`, tableName)
	_, err := p.Exec(stmt)
	return err
}
