package response

import "server/model/system"

type SysAreaResponse struct {
	Area system.SysArea `json:"area"`
}

type SysAreaTreeNode struct {
	system.SysArea
	Children []SysAreaTreeNode `json:"children,omitempty"`
}

type SysAreaTreeResponse struct {
	Tree []SysAreaTreeNode `json:"tree"`
}

type SysAreaListResponse struct {
	Areas []system.SysArea `json:"areas"`
}

type ImportAreaResponse struct {
	Success int    `json:"success"` // 成功导入数量
	Failed  int    `json:"failed"`  // 失败数量
	Message string `json:"message"` // 导入结果消息
}
