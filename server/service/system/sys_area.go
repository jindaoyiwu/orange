package system

import (
	"errors"
	"fmt"

	"server/global"
	"server/model/common/request"
	"server/model/system"
	systemReq "server/model/system/request"
	systemRes "server/model/system/response"

	"gorm.io/gorm"
)

type AreaService struct{}

var AreaServiceApp = new(AreaService)

// CreateArea 创建区域信息
func (areaService *AreaService) CreateArea(area system.SysArea) (err error) {
	// 检查区域编码是否已存在
	if !errors.Is(global.GVA_DB.Where("i = ?", area.I).First(&system.SysArea{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("区域编码已存在")
	}

	// 检查同一父级下是否存在相同名称
	if !errors.Is(global.GVA_DB.Where("n = ? AND p = ?", area.N, area.P).First(&system.SysArea{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("同一父级下区域名称已存在")
	}

	// 计算级别
	if area.P == 0 {
		area.Level = 1 // 省级
	} else {
		var parent system.SysArea
		if err = global.GVA_DB.Where("i = ?", area.P).First(&parent).Error; err != nil {
			return errors.New("父级区域不存在")
		}
		area.Level = parent.Level + 1
		if area.Level > 3 {
			return errors.New("不支持超过三级的区域层级")
		}
	}

	return global.GVA_DB.Create(&area).Error
}

// DeleteArea 删除区域信息
func (areaService *AreaService) DeleteArea(area system.SysArea) (err error) {
	// 检查是否有子区域
	var count int64
	err = global.GVA_DB.Model(&system.SysArea{}).Where("p = ?", area.I).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("存在子区域，不能删除")
	}

	return global.GVA_DB.Where("id = ?", area.ID).Delete(&area).Error
}

// DeleteAreasByIds 批量删除区域
func (areaService *AreaService) DeleteAreasByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var areas []system.SysArea
		if err := tx.Find(&areas, "id in ?", ids.Ids).Error; err != nil {
			return err
		}

		// 检查每个区域是否有子区域
		for _, area := range areas {
			var count int64
			if err := tx.Model(&system.SysArea{}).Where("p = ?", area.I).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return fmt.Errorf("区域 %s 存在子区域，不能删除", area.N)
			}
		}

		return tx.Delete(&[]system.SysArea{}, "id in ?", ids.Ids).Error
	})
}

// UpdateArea 更新区域信息
func (areaService *AreaService) UpdateArea(area system.SysArea) (err error) {
	var oldArea system.SysArea
	if err = global.GVA_DB.First(&oldArea, "id = ?", area.ID).Error; err != nil {
		return err
	}

	// 如果修改了区域编码，检查新编码是否已存在
	if oldArea.I != area.I {
		if !errors.Is(global.GVA_DB.Where("i = ? AND id != ?", area.I, area.ID).First(&system.SysArea{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("区域编码已存在")
		}
	}

	// 如果修改了名称或父级ID，检查同一父级下名称是否重复
	if oldArea.N != area.N || oldArea.P != area.P {
		if !errors.Is(global.GVA_DB.Where("n = ? AND p = ? AND id != ?", area.N, area.P, area.ID).First(&system.SysArea{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("同一父级下区域名称已存在")
		}
	}

	// 如果修改了父级ID，需要重新计算级别
	if oldArea.P != area.P {
		if area.P == 0 {
			area.Level = 1
		} else {
			var parent system.SysArea
			if err = global.GVA_DB.Where("i = ?", area.P).First(&parent).Error; err != nil {
				return errors.New("父级区域不存在")
			}
			area.Level = parent.Level + 1
			if area.Level > 3 {
				return errors.New("不支持超过三级的区域层级")
			}
		}
	}

	return global.GVA_DB.Save(&area).Error
}

// GetAreaList 获取区域分页列表
func (areaService *AreaService) GetAreaList(info systemReq.SysAreaSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysArea{})
	var areaList []system.SysArea

	// 搜索条件
	if info.N != "" {
		db = db.Where("n LIKE ?", "%"+info.N+"%")
	}
	if info.I != 0 {
		db = db.Where("i = ?", info.I)
	}
	if info.P != 0 {
		db = db.Where("p = ?", info.P)
	}
	if info.Y != "" {
		db = db.Where("y LIKE ?", "%"+info.Y+"%")
	}
	if info.Level != 0 {
		db = db.Where("level = ?", info.Level)
	}

	err = db.Count(&total).Error
	if err != nil {
		return areaList, total, err
	}

	// 排序
	OrderStr := "i ASC, id DESC"
	if info.OrderKey != "" {
		orderMap := make(map[string]bool, 8)
		orderMap["id"] = true
		orderMap["n"] = true
		orderMap["i"] = true
		orderMap["p"] = true
		orderMap["level"] = true
		orderMap["created_at"] = true
		orderMap["updated_at"] = true

		if !orderMap[info.OrderKey] {
			err = fmt.Errorf("非法的排序字段: %v", info.OrderKey)
			return areaList, total, err
		}
		OrderStr = info.OrderKey
		if info.Desc {
			OrderStr = info.OrderKey + " DESC"
		}
	}

	err = db.Limit(limit).Offset(offset).Order(OrderStr).Find(&areaList).Error
	return areaList, total, err
}

// GetAreaById 根据ID获取区域信息
func (areaService *AreaService) GetAreaById(id uint) (area system.SysArea, err error) {
	err = global.GVA_DB.First(&area, "id = ?", id).Error
	return
}

// GetAreaByAreaId 根据区域编码获取区域信息
func (areaService *AreaService) GetAreaByAreaId(areaId int) (area system.SysArea, err error) {
	err = global.GVA_DB.First(&area, "i = ?", areaId).Error
	return
}

// GetAreaTree 获取区域树形结构
func (areaService *AreaService) GetAreaTree(req systemReq.AreaTree) (tree []systemRes.SysAreaTreeNode, err error) {
	var areas []system.SysArea
	db := global.GVA_DB.Model(&system.SysArea{})

	// 过滤条件
	if req.Level != nil {
		db = db.Where("level = ?", *req.Level)
	}

	err = db.Order("i ASC").Find(&areas).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	areaMap := make(map[int][]systemRes.SysAreaTreeNode)
	nodeMap := make(map[int]*systemRes.SysAreaTreeNode)

	// 创建所有节点
	for _, area := range areas {
		node := systemRes.SysAreaTreeNode{
			SysArea:  area,
			Children: []systemRes.SysAreaTreeNode{},
		}
		nodeMap[area.I] = &node
		areaMap[area.P] = append(areaMap[area.P], node)
	}

	// 构建父子关系
	for _, area := range areas {
		if parentNode, exists := nodeMap[area.P]; exists && area.P != 0 {
			parentNode.Children = append(parentNode.Children, *nodeMap[area.I])
		}
	}

	// 获取根节点
	parentId := 0
	if req.ParentId != nil {
		parentId = *req.ParentId
	}
	tree = areaMap[parentId]

	return tree, nil
}

// GetAreasByParentId 根据父级ID获取子区域列表
func (areaService *AreaService) GetAreasByParentId(parentId int) (areas []system.SysArea, err error) {
	err = global.GVA_DB.Where("p = ?", parentId).Order("i ASC").Find(&areas).Error
	return
}

// ImportAreaData 导入区域数据
func (areaService *AreaService) ImportAreaData(req systemReq.ImportAreaReq) (result systemRes.ImportAreaResponse, err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 清空现有数据
		if req.ClearData {
			if err := tx.Unscoped().Delete(&system.SysArea{}, "1 = 1").Error; err != nil {
				return err
			}
		}

		successCount := 0
		failedCount := 0
		messages := []string{}

		for _, item := range req.Data {
			area := system.SysArea{
				N: item.N,
				I: item.I,
				P: item.P,
				Y: item.Y,
			}

			// 计算级别
			if area.P == 0 {
				area.Level = 1
			} else {
				var parent system.SysArea
				if err := tx.Where("i = ?", area.P).First(&parent).Error; err == nil {
					area.Level = parent.Level + 1
				} else {
					area.Level = 1 // 如果找不到父级，默认为1级
				}
			}

			// 检查是否已存在
			var existArea system.SysArea
			if !errors.Is(tx.Where("i = ?", area.I).First(&existArea).Error, gorm.ErrRecordNotFound) {
				// 存在则更新
				area.ID = existArea.ID
				if err := tx.Save(&area).Error; err != nil {
					failedCount++
					messages = append(messages, fmt.Sprintf("更新区域 %s 失败: %v", area.N, err))
				} else {
					successCount++
				}
			} else {
				// 不存在则创建
				if err := tx.Create(&area).Error; err != nil {
					failedCount++
					messages = append(messages, fmt.Sprintf("创建区域 %s 失败: %v", area.N, err))
				} else {
					successCount++
				}
			}
		}

		result.Success = successCount
		result.Failed = failedCount
		if len(messages) > 0 {
			result.Message = fmt.Sprintf("导入完成：成功 %d 条，失败 %d 条。", successCount, failedCount)
			if failedCount > 0 && len(messages) <= 10 {
				result.Message += " 失败原因：" + fmt.Sprintf("%v", messages)
			}
		} else {
			result.Message = fmt.Sprintf("导入完成：成功 %d 条", successCount)
		}

		return nil
	})
	return
}
