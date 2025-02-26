package consumer

type CommentMessage struct {
	MsgType      string `json:"type"`
	SendUserName string `json:"send_username"`
	Content      string `json:"content"`
}
