package driver

import "database/sql"

func DBConnection(driver, connectionString string) (*sql.DB, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
