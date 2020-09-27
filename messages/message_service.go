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
)

const (
	HELP_MESSAGE = `Бот напоминаний о событии. Введите сообщение, затем время напоминания. Пример:
	Текст сообщения. 2020-10-22 10:10`
)


func HandleMessage(res *types.WebhookReqBody) {
	if res.Message.Text == "/help" {
		sendMessage(res.Message.Chat.ID, HELP_MESSAGE)
		return
	} else {
		messageModel := parseMessage(res)
		db.AddMessage(messageModel)
		return
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