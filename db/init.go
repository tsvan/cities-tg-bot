package db

import (
	"app/configs"
	"database/sql"
	"fmt"
)

var database *sql.DB

//Connect postgres, return db object
func Connect() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", configs.HOST, configs.PORT, configs.DB_USER, configs.DB_PASSWORD, configs.DB_NAME)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}
	return db
}

//CloseDB - close connection with database
func CloseDB() error {
    return database.Close()
}

//CreateTables init tables
func CreateTables() {
	database = Connect()
	result, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS public.chats(
			id SERIAL,
			chat_id bigint NOT NULL,
			cities text[] NOT NULL DEFAULT '{}'::text[],
			started boolean NOT NULL,
			CONSTRAINT chats_pkey PRIMARY KEY (id)
		);
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
