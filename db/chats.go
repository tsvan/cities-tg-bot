package db

import (
	"app/types"
	"fmt"

	"github.com/lib/pq"
)

//AddChat - add new user chat model
func AddChat(model types.ChatModel) {
	query := `INSERT INTO public.chats(
		chat_id, cities, started)
	   VALUES ($1, $2, $3);`
	result, err := database.Exec(query, model.ChatID, pq.Array(model.Cities), model.Started)
	if err != nil {
		panic(err)
	}
}

//UpdateChatStatus - update game status, unset used cities
func UpdateChatStatus(model types.ChatModel) {
	query := `UPDATE public.chats SET started = $1 , cities = $2 WHERE chat_id = $3;`
	result, err := database.Exec(query, model.Started, pq.Array([]string{}), model.ChatID)
	if err != nil {
		panic(err)
	}
}

//UpdateChatCities - update cities used in current game
func UpdateChatCities(model types.ChatModel) {
	query := `UPDATE public.chats SET cities = $1  WHERE chat_id = $2;`
	result, err := database.Exec(query, pq.Array(model.Cities), model.ChatID)
	if err != nil {
		panic(err)
	}
}

//GetChatByChatID - get chat model by chat id
func GetChatByChatID(chatID int64) (types.ChatModel, error) {
	query := fmt.Sprintf("select * from public.chats WHERE chat_id=%d", chatID)
	fmt.Println(query)
	row := database.QueryRow(query)
	model := types.ChatModel{}

	err := row.Scan(&model.ID, &model.ChatID, pq.Array(&model.Cities), &model.Started)
	if err != nil {
		fmt.Println("chat not found")
		return model, err
	}

	return model, nil
}
