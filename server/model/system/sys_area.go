package system

import (
	"server/global"
)

// SysArea 三级区域表
type SysArea struct {
	global.GVA_MODEL
	I     int    `json:"i" gorm:"column:i;comment:区域编码;not null"`                   // 区域编码
	N     string `json:"n" gorm:"column:n;comment:区域名称;not null;size:50"`           // 区域名称
	P     int    `json:"p" gorm:"column:p;comment:父级ID，0表示顶级;not null;default:0"`   // 父级ID
	Y     string `json:"y" gorm:"column:y;comment:拼音前缀;not null;size:1;default:''"` // 拼音前缀
	Level int    `json:"level" gorm:"comment:层级：1-省/直辖市，2-市/区，3-县/区;not null"`      // 层级
}

func (SysArea) TableName() string {
	return "sys_area"
}
