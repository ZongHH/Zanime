package models

type ServiceLog struct {
	Level     string `json:"level"`
	LevelType string `json:"levelType"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	Detail    string `json:"detail"`
	TimeStamp string `json:"timestamp"`
}
