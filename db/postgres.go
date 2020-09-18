package db

import (
	"app/configs"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

//Connect postgres, return db object
func Connect() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s "+
		"sslmode=disable", configs.DB_USER, configs.DB_PASSWORD, configs.DB_NAME)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}

func Test() {
	db := Connect()
	result, err := db.Exec("insert into messages (id, text) values (2, 'test2')",
		"Apple", 72000)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LastInsertId()) // не поддерживается
	fmt.Println(result.RowsAffected()) // количество добавленных строк
}
