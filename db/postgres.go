package db

import (
	"app/configs"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//Connect postgres, return db object
func Connect() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable",configs.HOST, configs.PORT,configs.DB_USER, configs.DB_PASSWORD, configs.DB_NAME)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	//defer db.Close()

	return db
}

func CreateTables() {
	db := Connect()
	result, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id integer NOT NULL,
		message text NOT NULL,
		chat_id varchar(300) NOT NULL,
		PRIMARY KEY (id)
	  )`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}


func Test() {
	fmt.Println("start test insert")
	db := Connect()
	result, err := db.Exec("INSERT INTO public.messages(id, message, chat_id) VALUES (1, 'test from docker', 11);")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
}
