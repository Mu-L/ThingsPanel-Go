// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID                  string     `gorm:"column:id;primaryKey" json:"id"`
	Name                *string    `gorm:"column:name" json:"name"`
	PhoneNumber         string     `gorm:"column:phone_number;not null" json:"phone_number"`
	Email               string     `gorm:"column:email;not null" json:"email"`
	Status              *string    `gorm:"column:status;comment:用户状态 F-冻结 N-正常" json:"status"`                                                 // 用户状态 F-冻结 N-正常
	Authority           *string    `gorm:"column:authority;comment:权限类型 TENANT_ADMIN-租户管理员 TENANT_USER-租户用户 SYS_ADMIN-系统管理员" json:"authority"` // 权限类型 TENANT_ADMIN-租户管理员 TENANT_USER-租户用户 SYS_ADMIN-系统管理员
	Password            string     `gorm:"column:password;not null" json:"password"`
	TenantID            *string    `gorm:"column:tenant_id" json:"tenant_id"`
	Remark              *string    `gorm:"column:remark" json:"remark"`
	AdditionalInfo      *string    `gorm:"column:additional_info;default:{}" json:"additional_info"`
	CreatedAt           *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           *time.Time `gorm:"column:updated_at" json:"updated_at"`
	PasswordLastUpdated *time.Time `gorm:"column:password_last_updated" json:"password_last_updated"`
	LastVisitTime       *time.Time `gorm:"column:last_visit_time;comment:上次访问时间" json:"last_visit_time"` // 上次访问时间
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
