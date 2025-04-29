package model

// Organization 组织模型，代表租户的逻辑架构
type Organization struct {
	Base
	Name        string `gorm:"size:100;not null" json:"name"`          // 组织名称
	Description string `gorm:"size:500" json:"description"`            // 组织描述
	Logo        string `gorm:"size:255" json:"logo"`                   // 组织Logo URL
	Website     string `gorm:"size:255" json:"website"`                // 组织网站
	Status      string `gorm:"size:20;default:'active'" json:"status"` // 组织状态：active, inactive, suspended
	OwnerId     int64  `gorm:"not null" json:"owner_id"`               // 组织拥有者ID
}
