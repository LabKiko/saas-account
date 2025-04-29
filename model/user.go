package model

// User 用户模型，记录用户的基础信息
type User struct {
	Base
	Name     string `gorm:"size:100;not null" json:"name"`         // 用户姓名
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`     // 电子邮件
	Phone    string `gorm:"size:20;uniqueIndex" json:"phone"`      // 手机号
	Password string `gorm:"size:100;not null" json:"-"`            // 密码，不返回给前端
	Avatar   string `gorm:"size:255" json:"avatar"`                // 头像URL
	Status   string `gorm:"size:20;default:'active'" json:"status"` // 用户状态：active, inactive, suspended
}
