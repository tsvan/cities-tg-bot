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
	defer db.Close()
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
	defer db.Close()
}

func GetMessagesToSend() {
	db := Connect()
	rows, err := db.Query("select * from public.messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	messages := []types.MessageModel{}

	for rows.Next() {
		message := types.MessageModel{}
		err := rows.Scan(&message.ID, &message.Message, &message.ChatID, &message.CreateDate, &message.SendDate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messages = append(messages, message)
	}
	for _, message := range messages {
		fmt.Println(message.ID, message.Message, message.SendDate)
	}
	defer db.Close()

}

func GetCityByName(name string) (types.CityModel, error) {
	db := Connect()
	query := fmt.Sprintf("select * from public.geo WHERE city='%s'", name)
	fmt.Println(query)
	row := db.QueryRow("select * from public.geo WHERE city=$1", name)
	city := types.CityModel{}

	err := row.Scan(&city.ID, &city.CountryEn, &city.RegionEn, &city.CityEn, &city.Country, &city.Region, &city.City, &city.Lat, &city.Lng, &city.Population)
	if err != nil {
		fmt.Println("city not found")
		return city, err
	}
	fmt.Println(city.ID, city.City)
	
	defer db.Close()
	return city, nil
}

func GetCitiesByLetter(letter string) []types.CityModel {
	db := Connect()
	fmt.Println("letter", letter)
	query := fmt.Sprintf("select * from public.geo WHERE LOWER(city) LIKE '%s%s' LIMIT 2", letter, "%")
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("cities not found")
	}
	defer rows.Close()
	cities := []types.CityModel{}

	for rows.Next() {
		city := types.CityModel{}
		err := rows.Scan(&city.ID, &city.CountryEn, &city.RegionEn, &city.CityEn, &city.Country, &city.Region, &city.City, &city.Lat, &city.Lng, &city.Population)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cities = append(cities, city)
	}
	for _, city := range cities {
		fmt.Println(city.ID, city.City, city.Region)
	}

	defer db.Close()
	return cities

}
