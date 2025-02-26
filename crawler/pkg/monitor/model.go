package monitor

type log_info struct {
	Level     string `json:"level"`
	LevelType string `json:"levelType"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	Detail    string `json:"detail"`
	TimeStamp string `json:"timestamp"`
}

type logConfig struct {
	PoolSize    int
	ServiceName string
	TopicName   string
}
