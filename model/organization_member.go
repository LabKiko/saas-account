package model

// OrganizationMember 组织成员模型，记录用户和组织的关系
type OrganizationMember struct {
	Base
	OrganizationId int64  `gorm:"not null;index:idx_org_user,unique" json:"organization_id"` // 组织ID
	UserId         int64  `gorm:"not null;index:idx_org_user,unique" json:"user_id"`         // 用户ID
	Role           string `gorm:"size:50;not null;default:'member'" json:"role"`             // 角色：owner, admin, member
	Status         string `gorm:"size:20;default:'active'" json:"status"`                    // 状态：active, inactive
}
