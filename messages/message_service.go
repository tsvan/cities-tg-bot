package messages

import (
	"app/configs"
	"app/types"
	"app/db"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"fmt"
)

const (
	HELP_MESSAGE = `Бот для получения информации о городе и игры в города.
	Каждое сообщение воспринимает как название города за исключением команд.  
	Вы пишете название города, бот называет город на последнюю букву, если последняя буква - ьъый, то на предпоследнюю. 
	Проверки на соответствие города названого ботом с городом названным вами нет, каждым сообщением можно начинать игру заново `
	NO_CITY_MESSAGE = "Нет такого города"
	NO_CITIES_FOUND_MESSAGE = "Городов не найдено"
	INFO_NOT_FOUND = "Информация не найдена"
)


func HandleMessage(res *types.WebhookReqBody) {

	if (len([]rune(res.Message.Text)) <= 1) {
		return
	}
	if string([]rune(res.Message.Text)[0]) == "/" {
		handleCommands(res)
		return
	}
	chat, chatErr := db.GetChatByChatID(res.Message.Chat.ID)
	if chatErr != nil {
		fmt.Println("нет такого чата")
		return
	}
	if !chat.Started {
		fmt.Println("режим игры не включен")
		return
	}

	city, err :=db.GetCityByName(res.Message.Text)
	if err!=nil {
		fmt.Println("нету города такого")
		sendMessage(res.Message.Chat.ID, NO_CITY_MESSAGE)
	} else {
		if(!cityCheck(city.City, chat.Cities)) {
			fmt.Println("нужно назвать город на букву -")
			sendMessage(res.Message.Chat.ID, "нужно назвать город на букву -")
			return
		}
		if(stringInSlice(city.City, chat.Cities)) {
			fmt.Println("Такой город уже был")
			sendMessage(res.Message.Chat.ID, "Такой город уже был")
			return
		}
		randomCity, err := db.GetRandomCityByLetter(getLetter(city.City))
		if err == nil  {
			sendMessage(res.Message.Chat.ID, randomCity.City)
			chat.Cities = setCitiesList(chat.Cities, city.City, randomCity.City)
			db.UpdateChatCities(chat)
		} else {
			sendMessage(res.Message.Chat.ID, NO_CITIES_FOUND_MESSAGE)
		}
	}
	return
}

func handleCommands(res *types.WebhookReqBody) {
	tmp := strings.Split(res.Message.Text, " ")
	command:= tmp[0]
	param:= strings.Join(tmp[1:], " ")

	switch command {
		case "/help",("/help@"+configs.BOT_NAME) :
			sendMessage(res.Message.Chat.ID, HELP_MESSAGE)
		case "/start",("/start@"+configs.BOT_NAME) :
			chat, chatErr := db.GetChatByChatID(res.Message.Chat.ID)
			if chatErr != nil {
				var model types.ChatModel
				model.ChatID = res.Message.Chat.ID
				//model.Cities = []string{"Москва", "Анапа", "queues"}
				model.Cities = []string{}
				model.Started = true
				db.AddChat(model)
				sendMessage(res.Message.Chat.ID, "Новый чат добавлен")
			} else {
				chat.Started = true
				db.UpdateChatStatus(chat)
				sendMessage(res.Message.Chat.ID, "режим игры включен")
			}
		case "/stop",("/stop@"+configs.BOT_NAME) :
			chat, err := db.GetChatByChatID(res.Message.Chat.ID)
			if err != nil {
				fmt.Println("нет такого чата")
				return
			}
			chat.Started = false
			db.UpdateChatStatus(chat)
			sendMessage(res.Message.Chat.ID, "stop")
		case "/info",("/info@"+configs.BOT_NAME) :
			city, err :=db.GetCityByName(param)
			if err!=nil {
				sendMessage(res.Message.Chat.ID, INFO_NOT_FOUND)
			} else {
				sendMessage(res.Message.Chat.ID, 
					fmt.Sprintf("Страна: %s.\nРегион: %s.\nНаселение: %d." , city.Country, city.Region, city.Population))
			}
		default:
			sendMessage(res.Message.Chat.ID, "command not found")
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

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
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
	letter:= string(r[len(r)-1:])
	return letter
}

func cityCheck(city string, usedCities []string) bool {
	if(len(usedCities) < 1) {
		return true
	}
	lastCity := usedCities[len(usedCities)-1]
	r := []rune(strings.ToLower(city))
	firstletter:= string(r[0])
	fmt.Println(firstletter)
	fmt.Println(getLetter(lastCity))
	if(getLetter(lastCity) == firstletter) {
		return true
	}
	return false
}

func setCitiesList(cities []string, userCity string, randomCity string) []string {
	if(len(cities) > 6 ){
		cities = cities[2:]
	}
	result := append(cities, userCity, randomCity)
	return result
}