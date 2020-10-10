package messages

import (
	"app/configs"
	"app/db"
	"app/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	HELP_MESSAGE = `Бот для получения информации о городе и игры в города.
	Каждое сообщение воспринимает как название города за исключением команд.  
	Вы пишете название города, бот называет город на последнюю букву, если последняя буква - ьъый, то на предпоследнюю. 
	Проверки на соответствие города названого ботом с городом названным вами нет, каждым сообщением можно начинать игру заново `
	NO_CITY_MESSAGE = "Нет такого города"
	NO_CITIES_FOUND_MESSAGE = "Городов не найдено"
)

func HandleMessage(res *types.WebhookReqBody) {
	switch res.Message.Text {
	case "/help":
		sendMessage(res.Message.Chat.ID, HELP_MESSAGE)
	default:
		city, err := db.GetCityByName(res.Message.Text)
		if err != nil {
			fmt.Println("нету города такого")
			sendMessage(res.Message.Chat.ID, NO_CITY_MESSAGE)
		} else {
			randomCity, err := db.GetRandomCitiesByLetter(getLetter(city.City))
			if err == nil {
				sendMessage(res.Message.Chat.ID, randomCity.City)
			} else {
				sendMessage(res.Message.Chat.ID, NO_CITIES_FOUND_MESSAGE)
			}
		}
	}
}

func parseMessage(res *types.WebhookReqBody) types.MessageModel {
	message := strings.Split(res.Message.Text, ".")
	parsedTime := message[len(message)-1]
	parsedMessage := strings.ReplaceAll(res.Message.Text, parsedTime, "")
	var model types.MessageModel
	model.Message = parsedMessage
	model.ChatID = res.Message.Chat.ID
	model.SendDate = parsedTime
	model.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	return model
}

func sendMessage(chatID int64, text string) error {
	// Create the request body struct
	reqBody := &types.SendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot"+configs.TOKEN+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func getLetter(word string) string {
	replacer := strings.NewReplacer("ь", "", "ъ", "", "ы", "", "й", "")
	word = replacer.Replace(word)
	fmt.Println(word)
	r := []rune(word)
	letter := string(r[len(r)-1:])
	return letter
}
