package db

import (
	"app/configs"
	"app/types"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
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
		);

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

func AddChat(model types.ChatModel) {
	query := `INSERT INTO public.chats(
		chat_id, cities, started)
	   VALUES ($1, $2, $3);`
	db := Connect()
	result, err := db.Exec(query, model.ChatID, pq.Array(model.Cities), model.Started)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	defer db.Close()
}

func UpdateChatStatus(model types.ChatModel) {
	query := `UPDATE public.chats SET started = $1 , cities = $2 WHERE chat_id = $3;`
	db := Connect()
	result, err := db.Exec(query, model.Started, pq.Array([]string{}), model.ChatID)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	defer db.Close()
}

func UpdateChatCities(model types.ChatModel) {
	query := `UPDATE public.chats SET cities = $1  WHERE chat_id = $2;`
	db := Connect()
	result, err := db.Exec(query, pq.Array(model.Cities), model.ChatID)
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

	row := db.QueryRow("select * from public.geo WHERE LOWER(city)=$1", strings.ToLower(name))
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

func GetRandomCitiesByLetter(letter string) (types.CityModel, error) {
	db := Connect()
	fmt.Println("letter", letter)
	query := fmt.Sprintf("select * from public.geo WHERE LOWER(city) LIKE '%s%s' ORDER BY random() LIMIT 1", strings.ToLower(letter), "%")
	fmt.Println(query)

	row := db.QueryRow(query)
	city := types.CityModel{}

	err := row.Scan(&city.ID, &city.CountryEn, &city.RegionEn, &city.CityEn, &city.Country, &city.Region, &city.City, &city.Lat, &city.Lng, &city.Population)
	if err != nil {
		fmt.Println("no random city")
		return city, err
	}
	fmt.Println(city.ID, city.City)
	
	defer db.Close()
	return city, nil

}

func GetChatByChatID(chatID int64) (types.ChatModel, error) {
	db := Connect()
	query := fmt.Sprintf("select * from public.chats WHERE chat_id=%d", chatID)
	fmt.Println(query)
	row := db.QueryRow(query)
	model := types.ChatModel{}

	err := row.Scan(&model.ID, &model.ChatID, pq.Array(&model.Cities), &model.Started)
	if err != nil {
		fmt.Println("chat not found")
		return model, err
	}
	fmt.Println(model.ID, model.ChatID, model.Cities)
	
	defer db.Close()
	return model, nil
}
