package chat

type message struct {
	User string `json:"user"`
	Text string `json:"text"`
	Time string `json:"time"`
}
