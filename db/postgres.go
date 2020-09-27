package db

import (
	"app/configs"
	"app/types"
	"database/sql"
	"fmt"

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
func AddMessage(model types.MessageModel) {
	query := fmt.Sprintf(`INSERT INTO public.messages(
		message, chat_id, create_date, send_date)
	   VALUES ('%s', '%d', '%s', '%s');`, model.Message, model.ChatID, model.CreateDate, model.SendDate)
	db := Connect()
	result, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
