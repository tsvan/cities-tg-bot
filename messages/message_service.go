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
)

//HandleMessage - handles every user msg from tg
func HandleMessage(res *types.WebhookReqBody) {

	if len([]rune(res.Message.Text)) <= 1 {
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
	if !cityCheck(res.Message.Text, chat.Cities) {
		fmt.Println(types.WrongFirstLetter)
		sendMessage(res.Message.Chat.ID, types.WrongFirstLetter)
		return
	}

	city, err := db.GetCityByName(res.Message.Text)
	if err != nil {
		fmt.Println(types.CityNotFound)
		sendMessage(res.Message.Chat.ID, types.CityNotFound)
	} else {
		if stringInSlice(city.City, chat.Cities) {
			fmt.Println(types.CityAlreadyUsed)
			sendMessage(res.Message.Chat.ID, types.CityAlreadyUsed)
			return
		}
		randomCity, err := db.GetRandomCityByLetter(getLastLetter(city.City))
		if err == nil {
			sendMessage(res.Message.Chat.ID, ("<b>" + randomCity.City + "</b> " + "<i>(Страна - " + randomCity.Country + ")</i>"))
			chat.Cities = setCitiesList(chat.Cities, city.City, randomCity.City)
			db.UpdateChatCities(chat)
		} else {
			sendMessage(res.Message.Chat.ID, types.NoCitiesFound)
		}
	}
	return
}

//handleCommands - handle user commands
func handleCommands(res *types.WebhookReqBody) {
	tmp := strings.Split(res.Message.Text, " ")
	command := tmp[0]
	param := strings.Join(tmp[1:], " ")

	switch command {
	case "/help", ("/help@" + configs.BOT_NAME):
		sendMessage(res.Message.Chat.ID, types.HelpMessage)
	case "/start", ("/start@" + configs.BOT_NAME):
		chat, chatErr := db.GetChatByChatID(res.Message.Chat.ID)
		if chatErr != nil {
			var model types.ChatModel
			model.ChatID = res.Message.Chat.ID
			model.Cities = []string{}
			model.Started = true
			db.AddChat(model)
		} else {
			chat.Started = true
			db.UpdateChatStatus(chat)
		}
		sendMessage(res.Message.Chat.ID, types.GameStarted)
	case "/stop", ("/stop@" + configs.BOT_NAME):
		chat, err := db.GetChatByChatID(res.Message.Chat.ID)
		if err != nil {
			fmt.Println("нет такого чата")
			return
		}
		chat.Started = false
		db.UpdateChatStatus(chat)
		sendMessage(res.Message.Chat.ID, types.GameStoped)
	case "/info", ("/info@" + configs.BOT_NAME):
		city, err := db.GetCityByName(param)
		if err != nil {
			sendMessage(res.Message.Chat.ID, types.InfoNotFound)
		} else {
			link := GetWikiLink(city.City)
			sendMessage(res.Message.Chat.ID,
				fmt.Sprintf("<b>Страна:</b> %s.\n<b>Регион:</b> %s.\n<b>Население:</b> %d.\n%s", city.Country, city.Region, city.Population,link))
		}
	default:
		sendMessage(res.Message.Chat.ID, types.CommandNotFound)
	}

}

//sendMessage - send message to user
func sendMessage(chatID int64, text string) error {
	// Create the request body struct
	reqBody := &types.SendMessageReqBody{
		ChatID: chatID,
		Text:   text,
		ParseMode: "HTML",
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
