package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Config 配置结构体
type Config struct {
	// 数据库配置
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// 服务器配置
	ServerHost string
	ServerPort int

	// JWT配置
	JWTSecret     string
	JWTExpiration int // 过期时间（分钟）

	// 日志配置
	LogLevel  string
	LogOutput string

	// Snowflake配置
	SnowflakeWorkerID     int64
	SnowflakeDatacenterID int64

	// 其他配置
	Environment string
	Debug       bool
}

var (
	config     *Config
	configOnce sync.Once
)

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{
			// 默认数据库配置
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnvAsInt("DB_PORT", 5432),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "postgres"),
			DBName:     getEnv("DB_NAME", "saas_account"),
			DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

			// 默认服务器配置
			ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
			ServerPort: getEnvAsInt("SERVER_PORT", 8080),

			// 默认JWT配置
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
			JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 60*24), // 默认24小时

			// 默认日志配置
			LogLevel:  getEnv("LOG_LEVEL", "info"),
			LogOutput: getEnv("LOG_OUTPUT", "console"),

			// 默认Snowflake配置
			SnowflakeWorkerID:     getEnvAsInt64("SNOWFLAKE_WORKER_ID", 1),
			SnowflakeDatacenterID: getEnvAsInt64("SNOWFLAKE_DATACENTER_ID", 1),

			// 默认其他配置
			Environment: getEnv("ENVIRONMENT", "development"),
			Debug:       getEnvAsBool("DEBUG", true),
		}

		// 尝试从配置文件加载
		configFile := getEnv("CONFIG_FILE", "")
		if configFile != "" {
			loadConfigFromFile(configFile)
		}
	})

	return config
}

// 从配置文件加载配置
func loadConfigFromFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return
	}
}

// 获取环境变量
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// 获取环境变量并转换为int64
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// 获取环境变量并转换为字符串切片
func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}
