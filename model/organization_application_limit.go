package model

// OrganizationApplicationLimit 组织应用限制模型，记录租户下的应用权益
type OrganizationApplicationLimit struct {
	Base
	OrganizationApplicationId int64  `gorm:"not null;uniqueIndex" json:"organization_application_id"` // 组织应用ID
	PlanName                  string `gorm:"size:100;not null" json:"plan_name"`                      // 套餐名称
	MaxUsers                  int    `gorm:"default:5" json:"max_users"`                              // 最大用户数
	MaxStorage                int64  `gorm:"default:1073741824" json:"max_storage"`                   // 最大存储空间（字节）
	MaxRequests               int    `gorm:"default:10000" json:"max_requests"`                       // 最大请求数/天
	Features                  string `gorm:"type:jsonb" json:"features"`                              // 功能特性，JSON格式
	ExpiresAt                 int64  `json:"expires_at"`                                              // 过期时间
	AutoRenew                 bool   `gorm:"default:false" json:"auto_renew"`                         // 是否自动续费
}
