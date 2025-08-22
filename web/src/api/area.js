import service from '@/utils/request'

// 创建区域信息
export const createArea = (data) => {
  return service({
    url: '/area/createArea',
    method: 'post',
    data
  })
}

// 删除区域信息
export const deleteArea = (data) => {
  return service({
    url: '/area/deleteArea',
    method: 'delete',
    data
  })
}

// 批量删除区域
export const deleteAreasByIds = (data) => {
  return service({
    url: '/area/deleteAreasByIds',
    method: 'delete',
    data
  })
}

// 更新区域信息
export const updateArea = (data) => {
  return service({
    url: '/area/updateArea',
    method: 'put',
    data
  })
}

// 分页获取区域列表
export const getAreaList = (data) => {
  return service({
    url: '/area/getAreaList',
    method: 'post',
    data
  })
}

// 根据ID获取区域信息
export const getAreaById = (data) => {
  return service({
    url: '/area/getAreaById',
    method: 'post',
    data
  })
}

// 根据区域ID获取区域信息
export const getAreaByAreaId = (areaId) => {
  return service({
    url: `/area/getAreaByAreaId/${areaId}`,
    method: 'get'
  })
}

// 获取区域树形结构
export const getAreaTree = (data = {}) => {
  return service({
    url: '/area/getAreaTree',
    method: 'post',
    data
  })
}

// 根据父级ID获取子区域列表
export const getAreasByParentId = (parentId) => {
  return service({
    url: `/area/getAreasByParentId/${parentId}`,
    method: 'get'
  })
}

// 导入区域数据
export const importAreaData = (data) => {
  return service({
    url: '/area/importAreaData',
    method: 'post',
    data
  })
}
