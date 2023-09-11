package dto

type RoomRequest struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type RoomResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type RoomPage struct {
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Total int64           `json:"total"`
	Rooms []*RoomResponse `json:"rooms"`
}
