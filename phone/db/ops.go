package db

import (
	"database/sql"
	"fmt"
	"regexp"
)

type PhoneNumber struct {
	ID    int
	Value string
}

func (p *phoneDB) Create(rawNum string) (*PhoneNumber, error) {
	normNum := Normalize(rawNum)
	stmt := fmt.Sprintf("INSERT INTO %s(value) VALUES($1) RETURNING id, value", tableName)
	var res PhoneNumber
	if err := p.QueryRow(stmt, normNum).Scan(&res.ID, &res.Value); err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *phoneDB) FindByID(id int) (*PhoneNumber, error) {
	stmt := fmt.Sprintf(`SELECT id, value FROM %s WHERE id = $1`, tableName)
	var res PhoneNumber
	if err := p.QueryRow(stmt, id).Scan(&res.ID, &res.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (p *phoneDB) FindByValue(rawNum string) (*PhoneNumber, error) {
	norm := Normalize(rawNum)
	stmt := fmt.Sprintf(`SELECT id, value FROM %s WHERE value = $1`, tableName)
	var res PhoneNumber
	if err := p.QueryRow(stmt, norm).Scan(&res.ID, &res.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (p *phoneDB) Update(pNum *PhoneNumber, rawNum string) error {
	normal := Normalize(rawNum)
	stmt := `UPDATE phone_numbers SET value=$1 WHERE id=$2 RETURNING id, value`
	return p.QueryRow(stmt, normal, pNum.ID).Scan(&pNum.ID, &pNum.Value)
}

func (p *phoneDB) DeleteBy(column string, value interface{}) error {
	if column == "value" {
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("string value required for column \"value\"")
		}
		value = Normalize(strVal)
	}
	stmt := fmt.Sprintf(`DELETE FROM phone_numbers WHERE %s = $1`, column)
	_, err := p.Exec(stmt, value)
	return err
}

func (p *phoneDB) All() ([]*PhoneNumber, error) {
	stmt := fmt.Sprintf(`SELECT id, value FROM %s`, tableName)
	rows, err := p.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*PhoneNumber
	for rows.Next() {
		var pNum PhoneNumber
		if err := rows.Scan(&pNum.ID, &pNum.Value); err != nil {
			return nil, err
		}
		res = append(res, &pNum)
	}
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return res, nil
}

func Normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}
