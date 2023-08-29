package dto

type MessageRequestDto struct {
	RoomId string `json:"room_id"`
	Text   string `json:"text"`
}
