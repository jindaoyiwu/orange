package system

import (
	"server/middleware"

	"github.com/gin-gonic/gin"
)

type AreaRouter struct{}

func (s *AreaRouter) InitAreaRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	areaRouter := Router.Group("area").Use(middleware.OperationRecord())
	areaRouterWithoutRecord := Router.Group("area")
	areaPublicRouterWithoutRecord := RouterPub.Group("area")
	{
		areaRouter.POST("createArea", areaApi.CreateArea)               // 创建区域信息
		areaRouter.DELETE("deleteArea", areaApi.DeleteArea)             // 删除区域信息
		areaRouter.DELETE("deleteAreasByIds", areaApi.DeleteAreasByIds) // 批量删除区域
		areaRouter.PUT("updateArea", areaApi.UpdateArea)                // 更新区域信息
		areaRouter.POST("importAreaData", areaApi.ImportAreaData)       // 导入区域数据
	}
	{
		areaRouterWithoutRecord.POST("getAreaById", areaApi.GetAreaById)                        // 根据ID获取区域信息
		areaRouterWithoutRecord.POST("getAreaList", areaApi.GetAreaList)                        // 分页获取区域列表
		areaRouterWithoutRecord.POST("getAreaTree", areaApi.GetAreaTree)                        // 获取区域树形结构
		areaRouterWithoutRecord.GET("getAreaByAreaId/:areaId", areaApi.GetAreaByAreaId)         // 根据区域ID获取区域信息
		areaRouterWithoutRecord.GET("getAreasByParentId/:parentId", areaApi.GetAreasByParentId) // 根据父级ID获取子区域列表
	}
	{
		// 公共接口，无需权限验证
		areaPublicRouterWithoutRecord.GET("getAreaByAreaId/:areaId", areaApi.GetAreaByAreaId)         // 根据区域ID获取区域信息
		areaPublicRouterWithoutRecord.GET("getAreasByParentId/:parentId", areaApi.GetAreasByParentId) // 根据父级ID获取子区域列表
		areaPublicRouterWithoutRecord.POST("getAreaTree", areaApi.GetAreaTree)                        // 获取区域树形结构
	}
}
