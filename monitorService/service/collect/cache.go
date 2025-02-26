package collect

import (
	"fmt"
	"log"
	"monitorService/dao/mysql"
	"monitorService/models"
	"strings"
	"sync"
	"time"
)

const (
	maxCacheSize  = 100              // 最大缓存条数
	flushInterval = 10 * time.Second // 定时刷新间隔
)

type LogCache struct {
	logs   []models.ServiceLog
	mutex  sync.Mutex
	ticker *time.Ticker
	stopCh chan struct{}
}

var (
	cache *LogCache
	once  sync.Once
)

func InitCache() *LogCache {
	once.Do(func() {
		cache = &LogCache{
			logs:   make([]models.ServiceLog, 0, maxCacheSize),
			ticker: time.NewTicker(flushInterval),
			stopCh: make(chan struct{}),
		}
		go cache.startTimer()
	})
	return cache
}

func (c *LogCache) startTimer() {
	for {
		select {
		case <-c.ticker.C:
			c.Flush()
		case <-c.stopCh:
			c.ticker.Stop()
			return
		}
	}
}

func (c *LogCache) Add(log models.ServiceLog) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.logs = append(c.logs, log)

	// 达到最大缓存条数时刷新
	if len(c.logs) >= maxCacheSize {
		c.flushLocked()
	}
}

func (c *LogCache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.flushLocked()
}

func (c *LogCache) flushLocked() {
	if len(c.logs) == 0 {
		return
	}

	// 存储到数据库
	if err := saveToDB(c.logs); err != nil {
		log.Printf("保存日志到数据库失败: %v", err)
		return
	}

	// 清空缓存
	c.logs = make([]models.ServiceLog, 0, maxCacheSize)
}

func (c *LogCache) Stop() {
	close(c.stopCh)
	c.Flush() // 最后一次刷新
}

func saveToDB(logs []models.ServiceLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 构建批量插入SQL
	valueStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*6) // 假设每条日志有4个字段

	for _, log := range logs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			log.Level,
			log.LevelType,
			log.Service,
			log.Message,
			log.Detail,
			log.TimeStamp,
		)
	}

	stmt := fmt.Sprintf("INSERT INTO service_logs (level, level_type, service, message, detail, timestamp) VALUES %s",
		strings.Join(valueStrings, ","))

	// 执行批量插入
	result, err := mysql.DB.Exec(stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("执行插入失败: %v", err)
	}

	affected, _ := result.RowsAffected()
	log.Printf("成功保存 %d 条日志到数据库", affected)
	return nil
}
