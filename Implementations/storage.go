package implementations

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	ConnectionString string
	Db               *sql.DB
}

func (s *Storage) Open() bool {
	db, err := sql.Open("postgres", s.ConnectionString)
	if err != nil {

		fmt.Println(err)
		return false
	}

	s.Db = db
	return true
}

func (s *Storage) Single(sql string, params []interface{}) *sql.Row {
	row := s.Db.QueryRow(sql, params...)

	return row
}

func (s *Storage) Where(sql string, params []interface{}) *sql.Rows {

	var err error
	rows, err := s.Db.Query(sql, params...)

	if err != nil {
		fmt.Println("Failed to get accounts", err)
	}
	return rows
}

func (s *Storage) Add(sql *string, params *[]interface{}) *sql.Row {

	// Prepare the SQL statement with placeholders
	stmt, err := s.Db.Prepare(*sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(*params...)

	return result
}

func (s *Storage) Close() bool {
	err := s.Db.Close()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return true
}
