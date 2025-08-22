package system

import (
	"server/global"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/system"
	systemReq "server/model/system/request"
	systemRes "server/model/system/response"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AreaApi struct{}

// CreateArea 创建区域信息
// @Tags      SysArea
// @Summary   创建区域信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysArea                true  "区域信息"
// @Success   200   {object}  response.Response{msg=string}  "创建成功"
// @Router    /area/createArea [post]
func (areaApi *AreaApi) CreateArea(c *gin.Context) {
	var area system.SysArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if area.N == "" {
		response.FailWithMessage("区域名称不能为空", c)
		return
	}
	if area.I == 0 {
		response.FailWithMessage("区域编码不能为空", c)
		return
	}

	err = areaService.CreateArea(area)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteArea 删除区域信息
// @Tags      SysArea
// @Summary   删除区域信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysArea                true  "区域ID"
// @Success   200   {object}  response.Response{msg=string}  "删除成功"
// @Router    /area/deleteArea [delete]
func (areaApi *AreaApi) DeleteArea(c *gin.Context) {
	var area system.SysArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(area.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = areaService.DeleteArea(area)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAreasByIds 批量删除区域
// @Tags      SysArea
// @Summary   批量删除区域
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "批量删除区域"
// @Success   200   {object}  response.Response{msg=string}  "批量删除成功"
// @Router    /area/deleteAreasByIds [delete]
func (areaApi *AreaApi) DeleteAreasByIds(c *gin.Context) {
	var ids request.IdsReq
	err := c.ShouldBindJSON(&ids)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = areaService.DeleteAreasByIds(ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateArea 更新区域信息
// @Tags      SysArea
// @Summary   更新区域信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysArea                true  "区域信息"
// @Success   200   {object}  response.Response{msg=string}  "更新成功"
// @Router    /area/updateArea [put]
func (areaApi *AreaApi) UpdateArea(c *gin.Context) {
	var area system.SysArea
	err := c.ShouldBindJSON(&area)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(area.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if area.N == "" {
		response.FailWithMessage("区域名称不能为空", c)
		return
	}
	if area.I == 0 {
		response.FailWithMessage("区域编码不能为空", c)
		return
	}

	err = areaService.UpdateArea(area)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// GetAreaList 分页获取区域列表
// @Tags      SysArea
// @Summary   分页获取区域列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SysAreaSearch                              true  "分页获取区域列表"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "获取成功"
// @Router    /area/getAreaList [post]
func (areaApi *AreaApi) GetAreaList(c *gin.Context) {
	var pageInfo systemReq.SysAreaSearch
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := areaService.GetAreaList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAreaById 根据ID获取区域信息
// @Tags      SysArea
// @Summary   根据ID获取区域信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                                  true  "根据ID获取区域信息"
// @Success   200   {object}  response.Response{data=systemRes.SysAreaResponse}  "获取成功"
// @Router    /area/getAreaById [post]
func (areaApi *AreaApi) GetAreaById(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindJSON(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	area, err := areaService.GetAreaById(uint(idInfo.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAreaResponse{Area: area}, "获取成功", c)
}

// GetAreaByAreaId 根据区域ID获取区域信息
// @Tags      SysArea
// @Summary   根据区域ID获取区域信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     areaId  path      int                                              true  "区域ID"
// @Success   200     {object}  response.Response{data=systemRes.SysAreaResponse}  "获取成功"
// @Router    /area/getAreaByAreaId/{areaId} [get]
func (areaApi *AreaApi) GetAreaByAreaId(c *gin.Context) {
	areaId := c.Param("areaId")
	if areaId == "" {
		response.FailWithMessage("区域ID不能为空", c)
		return
	}

	areaIdInt := 0
	if id, err := strconv.Atoi(areaId); err != nil {
		response.FailWithMessage("区域ID格式错误", c)
		return
	} else {
		areaIdInt = id
	}

	area, err := areaService.GetAreaByAreaId(areaIdInt)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAreaResponse{Area: area}, "获取成功", c)
}

// GetAreaTree 获取区域树形结构
// @Tags      SysArea
// @Summary   获取区域树形结构
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.AreaTree                                  true  "获取区域树形结构"
// @Success   200   {object}  response.Response{data=systemRes.SysAreaTreeResponse}  "获取成功"
// @Router    /area/getAreaTree [post]
func (areaApi *AreaApi) GetAreaTree(c *gin.Context) {
	var req systemReq.AreaTree
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tree, err := areaService.GetAreaTree(req)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAreaTreeResponse{Tree: tree}, "获取成功", c)
}

// GetAreasByParentId 根据父级ID获取子区域列表
// @Tags      SysArea
// @Summary   根据父级ID获取子区域列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     parentId  path      int                                                true  "父级区域ID"
// @Success   200       {object}  response.Response{data=systemRes.SysAreaListResponse}  "获取成功"
// @Router    /area/getAreasByParentId/{parentId} [get]
func (areaApi *AreaApi) GetAreasByParentId(c *gin.Context) {
	parentId := c.Param("parentId")
	if parentId == "" {
		response.FailWithMessage("父级ID不能为空", c)
		return
	}

	parentIdInt := 0
	if id, err := strconv.Atoi(parentId); err != nil {
		response.FailWithMessage("父级ID格式错误", c)
		return
	} else {
		parentIdInt = id
	}

	areas, err := areaService.GetAreasByParentId(parentIdInt)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAreaListResponse{Areas: areas}, "获取成功", c)
}

// ImportAreaData 导入区域数据
// @Tags      SysArea
// @Summary   导入区域数据
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.ImportAreaReq                              true  "导入区域数据"
// @Success   200   {object}  response.Response{data=systemRes.ImportAreaResponse}  "导入成功"
// @Router    /area/importAreaData [post]
func (areaApi *AreaApi) ImportAreaData(c *gin.Context) {
	var req systemReq.ImportAreaReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if len(req.Data) == 0 {
		response.FailWithMessage("导入数据不能为空", c)
		return
	}

	result, err := areaService.ImportAreaData(req)
	if err != nil {
		global.GVA_LOG.Error("导入失败!", zap.Error(err))
		response.FailWithMessage("导入失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(result, "导入完成", c)
}
