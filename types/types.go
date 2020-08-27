package types

type webhookReqBody struct {
	Message struct {
		Text string
		Chat struct {
			ID int64
		}
	}
}

type sendMessageReqBody struct {
	ChatID int64
	Text   string
}
