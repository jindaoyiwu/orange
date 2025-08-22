<template>
  <div class="area-container">
    <div class="gva-search-box">
      <el-form ref="elSearchForm" :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="区域名称" prop="n">
          <el-input v-model="searchInfo.n" placeholder="请输入区域名称" clearable />
        </el-form-item>
        <el-form-item label="区域编码" prop="i">
          <el-input v-model.number="searchInfo.i" placeholder="请输入区域编码" clearable />
        </el-form-item>
        <el-form-item label="父级编码" prop="p">
          <el-input v-model.number="searchInfo.p" placeholder="请输入父级编码" clearable />
        </el-form-item>
        <el-form-item label="级别" prop="level">
          <el-select v-model="searchInfo.level" placeholder="请选择级别" clearable>
            <el-option label="省/直辖市" value="1" />
            <el-option label="市" value="2" />
            <el-option label="区/县" value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable>
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="openDialog">新增</el-button>
        <el-button 
          :disabled="!multipleSelection.length"
          style="margin-left: 10px;" 
          icon="Delete" 
          @click="onDelete"
        >
          删除
        </el-button>
        <el-button type="success" icon="Upload" @click="openImportDialog">导入数据</el-button>
        <el-button type="info" icon="Grid" @click="toggleTreeView">
          {{ showTree ? '列表视图' : '树形视图' }}
        </el-button>
      </div>
      
      <!-- 表格视图 -->
      <el-table
        v-if="!showTree"
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="90" />
        <el-table-column align="left" label="区域名称" prop="n" width="120" />
        <el-table-column align="left" label="区域编码" prop="i" width="120" />
        <el-table-column align="left" label="父级编码" prop="p" width="120" />
        <el-table-column align="left" label="级别" prop="level" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.level === 1" type="success">省/直辖市</el-tag>
            <el-tag v-else-if="scope.row.level === 2" type="warning">市</el-tag>
            <el-tag v-else-if="scope.row.level === 3" type="info">区/县</el-tag>
            <el-tag v-else>{{ scope.row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="拼音前缀" prop="y" width="100" />
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button 
              type="primary" 
              link 
              icon="Edit" 
              class="table-button" 
              @click="updateAreaFunc(scope.row)"
            >变更</el-button>
            <el-button 
              type="primary" 
              link 
              icon="Delete" 
              @click="deleteRow(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 树形视图 -->
      <el-tree
        v-else
        :data="treeData"
        :props="treeProps"
        node-key="areaId"
        :expand-on-click-node="false"
        :default-expand-all="false"
        class="area-tree"
      >
        <template #default="{ node, data }">
          <span class="tree-node">
            <span class="tree-label">
              {{ data.n }} ({{ data.i }})
              <el-tag size="small" :type="getLevelType(data.level)" class="level-tag">
                {{ getLevelName(data.level) }}
              </el-tag>
            </span>
            <span class="tree-actions">
              <el-button 
                type="primary" 
                size="small" 
                link
                @click.stop="updateAreaFunc(data)"
              >
                编辑
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                link
                @click.stop="deleteRow(data)"
              >
                删除
              </el-button>
            </span>
          </span>
        </template>
      </el-tree>

      <div class="gva-pagination">
        <el-pagination
          v-if="!showTree"
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    
    <!-- 新增/修改弹窗 -->
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="区域信息">
      <el-form ref="elForm" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="区域名称" prop="n">
          <el-input v-model="formData.n" clearable placeholder="请输入区域名称" />
        </el-form-item>
        <el-form-item label="区域编码" prop="i">
          <el-input v-model.number="formData.i" clearable placeholder="请输入区域编码" />
        </el-form-item>
        <el-form-item label="父级编码" prop="p">
          <el-input v-model.number="formData.p" clearable placeholder="请输入父级编码，0为顶级" />
        </el-form-item>
        <el-form-item label="拼音前缀" prop="y">
          <el-input v-model="formData.y" clearable placeholder="请输入拼音前缀" maxlength="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 导入弹窗 -->
    <el-dialog v-model="importDialogVisible" title="导入区域数据" width="50%">
      <el-form :model="importForm" label-width="100px">
        <el-form-item label="JSON数据">
          <el-input
            v-model="importForm.jsonData"
            type="textarea"
            :rows="10"
            placeholder="请粘贴JSON数据，格式如：[{&quot;n&quot;:&quot;北京&quot;,&quot;i&quot;:11,&quot;p&quot;:0,&quot;y&quot;:&quot;b&quot;}]"
          />
        </el-form-item>
        <el-form-item label="清空现有数据">
          <el-switch v-model="importForm.clearData" />
          <span class="form-tip">开启后将清空所有现有区域数据</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeImportDialog">取 消</el-button>
          <el-button type="primary" @click="handleImport">确认导入</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import {
  createArea,
  deleteArea,
  deleteAreasByIds,
  updateArea,
  getAreaById,
  getAreaList,
  getAreaTree,
  importAreaData
} from '@/api/area.js'
import { formatTimeToStr } from '@/utils/date'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'Area',
  data() {
    return {
      type: '',
      tableData: [],
      treeData: [],
      searchInfo: {},
      page: 1,
      total: 0,
      pageSize: 10,
      multipleSelection: [],
      formData: {
        n: '',
        i: null,
        p: 0,
        y: ''
      },
      dialogFormVisible: false,
      importDialogVisible: false,
      showTree: false,
      importForm: {
        jsonData: '',
        clearData: false
      },
      treeProps: {
        children: 'children',
        label: 'n'
      },
      rules: {
        n: [
          { required: true, message: '请输入区域名称', trigger: 'blur' }
        ],
        i: [
          { required: true, message: '请输入区域编码', trigger: 'blur' },
          { type: 'number', message: '区域编码必须为数字', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    async getTableData() {
      const table = await getAreaList({
        page: this.page,
        pageSize: this.pageSize,
        ...this.searchInfo
      })
      if (table.code === 0) {
        this.tableData = table.data.list
        this.total = table.data.total
        this.page = table.data.page
        this.pageSize = table.data.pageSize
      }
    },
    async getTreeData() {
      const tree = await getAreaTree({ status: true })
      if (tree.code === 0) {
        this.treeData = tree.data.tree
      }
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    deleteRow(row) {
      this.$confirm('此操作将永久删除该区域, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async() => {
        const res = await deleteArea({ ID: row.ID })
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '删除成功!'
          })
          if (this.tableData.length === 1 && this.page > 1) {
            this.page--
          }
          this.getTableData()
          if (this.showTree) {
            this.getTreeData()
          }
        }
      })
    },
    async onDelete() {
      const ids = this.multipleSelection.map(item => item.ID)
      const res = await deleteAreasByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        if (this.tableData.length === ids.length && this.page > 1) {
          this.page--
        }
        this.getTableData()
      }
    },

    async updateAreaFunc(row) {
      const res = await getAreaById({ ID: row.ID })
      this.type = 'update'
      if (res.code === 0) {
        this.formData = res.data.area
        this.dialogFormVisible = true
      }
    },
    closeDialog() {
      this.dialogFormVisible = false
      this.formData = {
        n: '',
        i: null,
        p: 0,
        y: ''
      }
    },
    async enterDialog() {
      await this.$refs.elForm.validate(async(valid) => {
        if (!valid) return
        let res
        switch (this.type) {
          case 'create':
            res = await createArea(this.formData)
            break
          case 'update':
            res = await updateArea(this.formData)
            break
          default:
            res = await createArea(this.formData)
            break
        }
        if (res.code === 0) {
          ElMessage({
            type: 'success',
            message: '创建/更改成功'
          })
          this.closeDialog()
          this.getTableData()
          if (this.showTree) {
            this.getTreeData()
          }
        }
      })
    },
    openDialog() {
      this.type = 'create'
      this.dialogFormVisible = true
    },
    openImportDialog() {
      this.importDialogVisible = true
    },
    closeImportDialog() {
      this.importDialogVisible = false
      this.importForm = {
        jsonData: '',
        clearData: false
      }
    },
    async handleImport() {
      if (!this.importForm.jsonData.trim()) {
        ElMessage.error('请输入JSON数据')
        return
      }
      
      try {
        const data = JSON.parse(this.importForm.jsonData)
        if (!Array.isArray(data)) {
          ElMessage.error('JSON数据格式错误，应为数组格式')
          return
        }
        
        const res = await importAreaData({
          data: data,
          clearData: this.importForm.clearData
        })
        
        if (res.code === 0) {
          ElMessage.success('导入成功：' + res.data.message)
          this.closeImportDialog()
          this.getTableData()
          if (this.showTree) {
            this.getTreeData()
          }
        }
      } catch (error) {
        ElMessage.error('JSON数据格式错误：' + error.message)
      }
    },
    async toggleTreeView() {
      this.showTree = !this.showTree
      if (this.showTree) {
        await this.getTreeData()
      }
    },
    onSubmit() {
      this.page = 1
      this.pageSize = 10
      this.getTableData()
    },
    onReset() {
      this.searchInfo = {}
      this.page = 1
      this.pageSize = 10
      this.getTableData()
    },
    handleSizeChange(val) {
      this.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.page = val
      this.getTableData()
    },
    formatDate(time) {
      if (time != null && time !== '') {
        return formatTimeToStr(time, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
    getLevelType(level) {
      const types = { 1: 'success', 2: 'warning', 3: 'info' }
      return types[level] || ''
    },
    getLevelName(level) {
      const names = { 1: '省', 2: '市', 3: '区' }
      return names[level] || level
    }
  }
}
</script>

<style lang="scss" scoped>
.area-container {
  :deep(.area-tree) {
    margin-top: 20px;
    .tree-node {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 14px;
      padding-right: 8px;
      
      .tree-label {
        display: flex;
        align-items: center;
        gap: 8px;
        
        .level-tag {
          margin-left: 8px;
        }
      }
      
      .tree-actions {
        display: none;
      }
    }
    
    .el-tree-node__content:hover .tree-actions {
      display: block;
    }
  }
  
  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-left: 10px;
  }
}
</style>
