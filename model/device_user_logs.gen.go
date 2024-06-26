// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDeviceUserLog = "device_user_logs"

// DeviceUserLog mapped from table <device_user_logs>
type DeviceUserLog struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	DeviceNum int32     `gorm:"column:device_nums;not null" json:"device_nums"`
	DeviceOn  int32     `gorm:"column:device_on;not null" json:"device_on"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	TenantID  string    `gorm:"column:tenant_id;not null;comment:租户 id" json:"tenant_id"` // 租户 id
}

// TableName DeviceUserLog's table name
func (*DeviceUserLog) TableName() string {
	return TableNameDeviceUserLog
}
