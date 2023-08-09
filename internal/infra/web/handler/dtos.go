package handler

type CreateRoomRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}
