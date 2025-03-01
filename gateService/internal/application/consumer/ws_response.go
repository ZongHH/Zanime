package consumer

type NotificationMessage struct {
	MsgType      string `json:"type"`
	SendUserName string `json:"send_username"`
	Title        string `json:"title"`
	Content      string `json:"content"`
}
