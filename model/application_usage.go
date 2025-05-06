package model

// ApplicationUsage 应用使用记录模型，记录应用的使用情况
type ApplicationUsage struct {
	Base
	ApplicationId  int64  `gorm:"not null;index" json:"application_id"`  // 组织应用ID
	UserId         *int64 `gorm:"index" json:"user_id"`                  // 用户ID，可为空表示系统操作
	UsageType      string `gorm:"size:50;not null" json:"usage_type"`    // 使用类型：api_call, storage, feature
	UsageAmount    int64  `gorm:"not null" json:"usage_amount"`          // 使用量
	UsageDate      int64  `gorm:"not null;index" json:"usage_date"`      // 使用日期
	Details        string `gorm:"type:jsonb" json:"details"`             // 详细信息，JSON格式
	OrganizationId int64  `gorm:"not null;index" json:"organization_id"` // 组织ID
	MemberId       int64  `gorm:"not null;index" json:"member_id"`       // 组织成员ID
}
