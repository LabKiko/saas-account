package utils

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	idGenerator     *Snowflake
	idGeneratorOnce sync.Once
)

// GetIDGenerator 获取全局ID生成器实例
func GetIDGenerator() *Snowflake {
	idGeneratorOnce.Do(func() {
		// 从环境变量获取机器ID和数据中心ID，如果不存在则使用默认值
		workerIDStr := os.Getenv("SNOWFLAKE_WORKER_ID")
		datacenterIDStr := os.Getenv("SNOWFLAKE_DATACENTER_ID")

		workerID := int64(1)
		datacenterID := int64(1)

		if workerIDStr != "" {
			if id, err := strconv.ParseInt(workerIDStr, 10, 64); err == nil {
				workerID = id
			}
		}

		if datacenterIDStr != "" {
			if id, err := strconv.ParseInt(datacenterIDStr, 10, 64); err == nil {
				datacenterID = id
			}
		}

		var err error
		idGenerator, err = NewSnowflake(workerID, datacenterID)
		if err != nil {
			log.Fatalf("初始化ID生成器失败: %v", err)
		}
	})

	return idGenerator
}

// GenerateID 生成一个新的ID
func GenerateID() int64 {
	id, err := GetIDGenerator().NextID()
	if err != nil {
		log.Printf("生成ID失败: %v", err)
		// 如果生成失败，返回当前时间戳作为备用ID
		return time.Now().UnixNano()
	}
	return id
}

// GenerateStringID 生成一个字符串格式的ID
func GenerateStringID() string {
	return strconv.FormatInt(GenerateID(), 10)
}
