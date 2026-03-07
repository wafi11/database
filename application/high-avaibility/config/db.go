package config

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig struct {
	connectionString string
}

func NewDBConfig(connectionString string) *DBConfig {
	return &DBConfig{
		connectionString: connectionString,
	}
}

func (config *DBConfig) Connect() (*sql.DB, error) {

	db, err := sql.Open("pgx", config.connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		return nil, err
	}
	fmt.Println("Connected to database successfully")
	return db, nil
}
