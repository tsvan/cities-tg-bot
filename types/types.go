package types

//WebhookReqBody responce received from tg
//https://core.telegram.org/bots/api#update
type WebhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

//SendMessageReqBody request sending to chat
//https://core.telegram.org/bots/api#sendmessage
type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	ParseMode  string `json:"parse_mode"`
}

type MessageModel struct {
	ID int64
	Message string
	ChatID int64
	CreateDate string
	SendDate string
}

type CityModel struct {
	ID int64
	CountryEn string
	RegionEn string
	CityEn string
	Country string
	Region string
	City string
	Lat string
	Lng string
	Population int64
}

type ChatModel struct {
	ID int64
	ChatID int64
	Cities []string
	Started bool
}