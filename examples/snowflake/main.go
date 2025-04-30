package main

import (
	"fmt"
	"saas-account/utils"
	"time"
)

func main() {
	// 初始化日志
	utils.InitLogger("debug", "console")
	logger := utils.GetLogger()

	// 生成ID
	logger.Info("开始生成ID")

	// 生成10个ID
	for i := 0; i < 10; i++ {
		id := utils.GenerateID()
		stringID := utils.GenerateStringID()
		fmt.Printf("生成的ID: %d (%s)\n", id, stringID)
		time.Sleep(10 * time.Millisecond) // 等待一段时间，确保ID不同
	}

	logger.Info("ID生成完成")
}
