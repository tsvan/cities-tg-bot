package main

import (
	"app/types"
	"app/db"
	"app/messages"
	"encoding/json"
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
	messages.HandleMessage(body)


	//if err := sendTest(body.Message.Chat.ID); err != nil {
	//	fmt.Println("error in sending reply:", err)
	//	return
	//}
	//db.AddMessage(body.Message.Chat.ID, body.Message.Text)

	fmt.Println("reply sent")
}

func main() {
	db.CreateTables()
	http.ListenAndServe(":8000", http.HandlerFunc(Handler))
}
