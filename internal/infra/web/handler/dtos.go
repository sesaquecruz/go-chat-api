package handler

type RoomRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type RoomResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}
