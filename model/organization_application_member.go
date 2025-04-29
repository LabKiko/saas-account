package model

// OrganizationApplicationMember 组织应用成员模型，记录用户和组织应用的关系
type OrganizationApplicationMember struct {
	Base
	ApplicationId int64  `gorm:"not null;index:idx_app_user,unique" json:"application_id"` // 组织应用ID
	MemberId      int64  `gorm:"not null;index:idx_app_user,unique" json:"member_id"`      // 用户ID
	Role          string `gorm:"size:50;not null;default:'user'" json:"role"`              // 角色：admin, user, guest
	Status        string `gorm:"size:20;default:'active'" json:"status"`                   // 状态：active, inactive
	Permissions   string `gorm:"type:jsonb" json:"permissions"`                            // 权限配置，JSON格式
}
