package main

type Update struct {
	Update  int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}
type Chat struct {
	chatId int `json:"chat_id"`
}
type RestResponse struct {
	Result []Update `json:"result"`
}
