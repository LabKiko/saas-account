package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// IsValidEmail 检查邮箱格式是否有效
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// IsValidPhone 检查手机号格式是否有效（简单校验）
func IsValidPhone(phone string) bool {
	pattern := `^\+?[0-9]{10,15}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(phone)
}

// IsStrongPassword 检查密码强度
func IsStrongPassword(password string) bool {
	// 密码长度至少8位
	if len(password) < 8 {
		return false
	}

	// 包含至少一个数字
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// 包含至少一个小写字母
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// 包含至少一个大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 包含至少一个特殊字符
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasNumber && hasLower && hasUpper && hasSpecial
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// MD5Hash 计算字符串的MD5哈希值
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// FormatTime 格式化时间为指定格式
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return time.Parse(format, timeStr)
}

// TruncateString 截断字符串到指定长度，并添加省略号
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// SlugifyString 将字符串转换为URL友好的格式
func SlugifyString(s string) string {
	// 转换为小写
	s = strings.ToLower(s)
	// 替换非字母数字字符为连字符
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	// 移除开头和结尾的连字符
	s = strings.Trim(s, "-")
	return s
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unitIndex := 0
	value := float64(size)

	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", value, units[unitIndex])
}
