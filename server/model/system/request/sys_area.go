package request

import (
	"server/model/common/request"
	"server/model/system"
)

type SysAreaSearch struct {
	system.SysArea
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序字段
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// AreaTree 区域树结构请求
type AreaTree struct {
	ParentId *int  `json:"parentId"` // 父级ID，为空则获取所有
	Level    *int  `json:"level"`    // 级别过滤
	Status   *bool `json:"status"`   // 状态过滤
}

// ImportAreaData 导入区域数据请求 - 直接使用用户JSON格式
type ImportAreaData struct {
	N string `json:"n"` // 名称
	I int    `json:"i"` // 区域编码
	P int    `json:"p"` // 父级ID
	Y string `json:"y"` // 拼音前缀
}

type ImportAreaReq struct {
	Data      []ImportAreaData `json:"data" binding:"required"` // 区域数据
	ClearData bool             `json:"clearData"`               // 是否清空现有数据
}
