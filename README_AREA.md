# 三级区域管理功能使用说明

## 📊 功能概述

基于您提供的SQL表结构创建的三级区域管理系统，支持省市区三级层级管理和JSON数据导入。

## 🗄️ 数据库表结构

```sql
CREATE TABLE `sys_area` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '区域ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  `i` int NOT NULL COMMENT '区域编码',
  `n` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '区域名称',
  `p` int NOT NULL DEFAULT '0' COMMENT '父级ID，0表示顶级',
  `y` char(1) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '拼音前缀',
  `level` tinyint NOT NULL COMMENT '层级：1-省/直辖市，2-市/区，3-县/区',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`p`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='三级区域表';
```

## 📥 JSON数据导入格式

支持您提供的JSON格式直接导入：

```json
[
  {"n":"北京","i":11,"p":0,"y":"b"},
  {"n":"北京","i":1101,"p":11,"y":"b"},
  {"n":"东城","i":110101,"p":1101,"y":"d"},
  {"n":"西城","i":110102,"p":1101,"y":"x"},
  {"n":"朝阳","i":110105,"p":1101,"y":"c"},
  {"n":"丰台","i":110106,"p":1101,"y":"f"}
]
```

字段说明：
- `n`: 区域名称 (name)
- `i`: 区域编码 (id) 
- `p`: 父级编码 (parent_id)，0为顶级
- `y`: 拼音前缀 (first letter)

## 🚀 API接口

### 基础CRUD操作
- `POST /area/createArea` - 创建区域
- `PUT /area/updateArea` - 更新区域 
- `DELETE /area/deleteArea` - 删除区域
- `DELETE /area/deleteAreasByIds` - 批量删除
- `POST /area/getAreaList` - 分页查询
- `POST /area/getAreaById` - 根据ID查询

### 树形查询
- `POST /area/getAreaTree` - 获取区域树
- `GET /area/getAreasByParentId/:parentId` - 获取子区域
- `GET /area/getAreaByAreaId/:areaId` - 根据编码查询

### 数据导入
- `POST /area/importAreaData` - 导入JSON数据

## 🎨 前端功能

### 管理界面
- ✅ 列表视图 / 树形视图切换
- ✅ 多条件搜索筛选
- ✅ 新增/编辑/删除操作
- ✅ 批量操作支持

### 数据导入
- ✅ JSON格式导入
- ✅ 可选清空现有数据
- ✅ 导入结果反馈

## 📝 使用步骤

1. **启动服务**
   ```bash
   # 后端
   cd server && go run main.go
   
   # 前端  
   cd web && npm run serve
   ```

2. **访问管理页面**
   - 登录系统后找到"系统管理" → "区域管理"

3. **导入数据**
   - 点击"导入数据"按钮
   - 粘贴您的JSON数据到文本框
   - 选择是否清空现有数据
   - 点击"确认导入"

## 📁 文件清单

### 后端文件
- `server/model/system/sys_area.go` - 数据模型
- `server/model/system/request/sys_area.go` - 请求模型
- `server/model/system/response/sys_area.go` - 响应模型
- `server/api/v1/system/sys_area.go` - API控制器
- `server/service/system/sys_area.go` - 业务服务
- `server/router/system/sys_area.go` - 路由配置
- `server/source/system/sys_area.go` - 数据初始化

### 前端文件
- `web/src/api/area.js` - API封装
- `web/src/view/system/area/area.vue` - 管理页面

### 测试数据
- `test_area_data.json` - 示例JSON数据

## ⚠️ 注意事项

1. **数据约束**：
   - 区域编码(`i`)必须唯一
   - 同一父级下区域名称(`n`)不能重复
   - 最多支持3级层级结构
   - 删除时检查是否有子区域

2. **导入规则**：
   - 如果区域编码已存在，会更新现有记录
   - 如果父级区域不存在，该区域的层级会设为1
   - 选择"清空现有数据"会删除所有区域记录

3. **层级计算**：
   - `p=0`: level=1 (省/直辖市)
   - 有父级: level=父级level+1
   - 自动计算，无需手动设置

现在您可以直接使用您提供的JSON格式导入区域数据了！🎯
