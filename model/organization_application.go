package model

// OrganizationApplication 组织应用模型，代表租户名下的应用
type OrganizationApplication struct {
	Base
	OrganizationId int64  `gorm:"not null;index" json:"organization_id"`  // 组织ID
	Name           string `gorm:"size:100;not null" json:"name"`          // 应用名称
	Description    string `gorm:"size:500" json:"description"`            // 应用描述
	AppKey         string `gorm:"size:100;uniqueIndex" json:"app_key"`    // 应用唯一标识
	AppSecret      string `gorm:"size:100" json:"-"`                      // 应用密钥，不返回给前端
	Status         string `gorm:"size:20;default:'active'" json:"status"` // 应用状态：active, inactive, suspended
	Type           string `gorm:"size:50;not null" json:"type"`           // 应用类型
	Config         string `gorm:"type:jsonb" json:"config"`               // 应用配置，JSON格式
}
