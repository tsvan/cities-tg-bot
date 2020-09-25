package db

import (
	"app/configs"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

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

//CreateTables init tables
func CreateTables() {
	db := Connect()
	result, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS public.messages(
			id SERIAL,
			message text NOT NULL,
			chat_id bigint  NOT NULL,
			create_date timestamp NOT NULL,
			send_date timestamp NOT NULL,
			CONSTRAINT messages_pkey PRIMARY KEY (id)
		)
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

//AddMessage - save user message from tg to database
func AddMessage(chatID int64, message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(`INSERT INTO public.messages(
		message, chat_id, create_date, send_date)
	   VALUES ('%s', '%d', '%s', '%s');`, message, chatID, currentTime, currentTime)
	db := Connect()
	result, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
