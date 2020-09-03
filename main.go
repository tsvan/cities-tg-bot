package main

import (
	"app/configs"
	"app/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//Handler handle webhook from tg
func Handler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &types.WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}
	fmt.Println("message", body)

	if err := sendTest(body.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}

	fmt.Println("reply sent")
}

func sendTest(chatID int64) error {
	// Create the request body struct
	reqBody := &types.SendMessageReqBody{
		ChatID: chatID,
		Text:   "test",
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

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(Handler))
}
