package web

type RoomRequestDto struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type RoomResponseDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}