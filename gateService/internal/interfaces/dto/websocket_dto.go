package dto

type EstablishWebSocketRequest struct {
	UserID int
}

type EstablishWebSocketResponse struct {
	Code int `json:"code"`
}
