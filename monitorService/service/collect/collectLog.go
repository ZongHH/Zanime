package collect

import (
	"encoding/json"
	"fmt"
	"monitorService/models"
	"monitorService/service/websocket"
)

func collect(msg []byte) error {
	var logInfo models.ServiceLog
	err := json.Unmarshal(msg, &logInfo)
	if err != nil {
		return fmt.Errorf("json.Unmarshal %v", err)
	}

	// 添加到缓存
	logCache.Add(logInfo)

	// 发送到websocket
	websocket.WSManager.SendToAll(msg)
	return nil
}
