<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">系统日志</h2>
        <p class="text-gray-600">记录系统中的所有操作日志</p>
      </div>
      <div class="space-x-2">
        <el-popconfirm
          confirm-button-text="确定"
          cancel-button-text="取消"
          title="确定要清空所有日志吗？此操作不可恢复！"
          confirm-button-type="danger"
          @confirm="handleClearLogs"
        >
          <template #reference>
            <el-button type="danger" :disabled="selectedLogs.length === 0" style="color: white !important; background-color: #ef4444 !important;">
              <el-icon class="mr-2"><Delete /></el-icon>
              清空日志
            </el-button>
          </template>
        </el-popconfirm>
        <el-button 
          type="danger" 
          :disabled="selectedLogs.length === 0"
          @click="handleDeleteSelected"
          style="color: white !important; background-color: #ef4444 !important;"
        >
          <el-icon class="mr-2"><Delete /></el-icon>
          删除选中
        </el-button>
      </div>
    </div>

    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden">
      <div class="mb-4 flex flex-wrap gap-4">
        <el-input 
          v-model="searchUsername" 
          placeholder="搜索用户名" 
          class="max-w-xs" 
          clearable
          prefix-icon="Search"
        />
        <el-select v-model="searchModule" placeholder="操作模块" class="w-32" clearable>
          <el-option 
            v-for="module in modules" 
            :key="module" 
            :label="getModuleLabel(module)" 
            :value="module" 
          />
        </el-select>
        <el-select v-model="searchAction" placeholder="操作类型" class="w-32" clearable>
          <el-option 
            v-for="action in actions" 
            :key="action" 
            :label="getActionLabel(action)" 
            :value="action" 
          />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          class="w-72"
        />
        <el-button type="primary" @click="handleSearch" class="text-white" style="color: white !important; background-color: #3b82f6 !important;">
          <el-icon class="mr-1"><Search /></el-icon>
          搜索
        </el-button>
        <el-button @click="handleReset" class="text-gray-600">
          <el-icon class="mr-1"><Refresh /></el-icon>
          重置
        </el-button>
      </div>

      <el-table 
        :data="logs" 
        border 
        stripe 
        v-loading="loading"
        class="rounded-lg overflow-hidden"
        :header-cell-style="{ background: '#f9fafb', color: '#374151', fontWeight: '600' }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="操作用户" min-width="120" />
        <el-table-column prop="module" label="操作模块" width="120">
          <template #default="{ row }">
            <el-tag :type="getModuleType(row.module)" class="rounded-full px-3 py-1" effect="plain">
              {{ getModuleLabel(row.module) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)" class="rounded-full px-3 py-1" effect="plain">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="操作资源" min-width="180" />
        <el-table-column prop="detail" label="操作详情" min-width="250" show-overflow-tooltip />
        <el-table-column prop="ip" label="操作IP" width="140" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="row.status === 1 ? 'success' : 'danger'"
              class="rounded-full px-3 py-1"
              effect="plain"
            >
              {{ row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" min-width="180" />
      </el-table>

      <div class="flex justify-end mt-4">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          background
          class="!p-0"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getLogList, deleteLogs, clearLogs, getLogModules, getLogActions, type SystemLog } from '@/api/log'

const logs = ref<SystemLog[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchUsername = ref('')
const searchModule = ref('')
const searchAction = ref('')
const dateRange = ref<[string, string] | null>(null)
const modules = ref<string[]>([])
const actions = ref<string[]>([])
const selectedLogs = ref<SystemLog[]>([])

// 获取日志列表
const fetchLogs = async () => {
  try {
    loading.value = true
    
    const params = {
      username: searchUsername.value,
      module: searchModule.value,
      action: searchAction.value,
      start_time: dateRange.value ? dateRange.value[0] + ' 00:00:00' : '',
      end_time: dateRange.value ? dateRange.value[1] + ' 23:59:59' : '',
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    const res = await getLogList(params)
    logs.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('Failed to fetch logs:', error)
    ElMessage.error('获取日志列表失败')
  } finally {
    loading.value = false
  }
}

// 获取日志模块列表
const fetchModules = async () => {
  try {
    modules.value = await getLogModules()
  } catch (error) {
    console.error('Failed to fetch modules:', error)
  }
}

// 获取日志操作类型列表
const fetchActions = async () => {
  try {
    actions.value = await getLogActions()
  } catch (error) {
    console.error('Failed to fetch actions:', error)
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchLogs()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchLogs()
}

// 处理搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchLogs()
}

// 处理重置
const handleReset = () => {
  searchUsername.value = ''
  searchModule.value = ''
  searchAction.value = ''
  dateRange.value = null
  currentPage.value = 1
  fetchLogs()
}

// 处理选择变化
const handleSelectionChange = (selection: SystemLog[]) => {
  selectedLogs.value = selection
}

// 处理删除选中
const handleDeleteSelected = async () => {
  if (selectedLogs.value.length === 0) {
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedLogs.value.length} 条日志吗？`, 
      '提示', 
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedLogs.value.map(log => log.id)
    await deleteLogs(ids)
    ElMessage.success('删除成功')
    fetchLogs()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理清空日志
const handleClearLogs = async () => {
  try {
    await clearLogs()
    ElMessage.success('清空成功')
    fetchLogs()
  } catch (error) {
    ElMessage.error('清空失败')
  }
}

// 获取模块标签
const getModuleLabel = (module: string) => {
  const moduleLabels: Record<string, string> = {
    user: '用户模块',
    role: '角色模块',
    system: '系统模块',
    auth: '认证模块'
  }
  return moduleLabels[module] || module
}

// 获取操作类型标签
const getActionLabel = (action: string) => {
  const actionLabels: Record<string, string> = {
    login: '登录',
    logout: '登出',
    create: '创建',
    update: '更新',
    delete: '删除',
    view: '查看',
    export: '导出',
    import: '导入',
    enable: '启用',
    disable: '禁用'
  }
  return actionLabels[action] || action
}

// 获取模块标签类型
const getModuleType = (module: string) => {
  const moduleTypes: Record<string, string> = {
    user: 'primary',
    role: 'success',
    system: 'warning',
    auth: 'info'
  }
  return moduleTypes[module] || 'info'
}

// 获取操作类型标签类型
const getActionType = (action: string) => {
  const actionTypes: Record<string, string> = {
    login: 'info',
    logout: 'info',
    create: 'success',
    update: 'primary',
    delete: 'danger',
    view: 'info',
    export: 'warning',
    import: 'warning',
    enable: 'success',
    disable: 'danger'
  }
  return actionTypes[action] || 'info'
}

onMounted(() => {
  fetchLogs()
  fetchModules()
  fetchActions()
})
</script>

<style scoped>
:deep(.el-input__wrapper) {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  padding: 0 1rem;
}

:deep(.el-input__wrapper:hover) {
  border-color: #d1d5db;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
}

:deep(.el-input__inner) {
  height: 2.5rem;
}

:deep(.el-input__prefix) {
  color: #6b7280;
}

:deep(.el-button) {
  border-radius: 0.5rem;
  transition: all 0.2s;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-select__wrapper) {
  border-radius: 0.5rem;
}

:deep(.el-table) {
  border-radius: 0.5rem;
  overflow: hidden;
}

:deep(.el-table th) {
  font-weight: 600;
}

:deep(.el-table td) {
  padding: 0.75rem;
}

:deep(.el-tag) {
  border-radius: 9999px;
  padding: 0.25rem 0.75rem;
}

:deep(.el-tag--danger.is-plain) {
  background-color: #fee2e2;
  border-color: #fee2e2;
  color: #dc2626;
}

:deep(.el-tag--success.is-plain) {
  background-color: #dcfce7;
  border-color: #dcfce7;
  color: #16a34a;
}

:deep(.el-tag--info.is-plain) {
  background-color: #f3f4f6;
  border-color: #f3f4f6;
  color: #4b5563;
}

:deep(.el-tag--primary.is-plain) {
  background-color: #e0f2fe;
  border-color: #e0f2fe;
  color: #0ea5e9;
}

:deep(.el-tag--warning.is-plain) {
  background-color: #fef3c7;
  border-color: #fef3c7;
  color: #d97706;
}

:deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #3b82f6;
}

:deep(.el-date-editor) {
  --el-date-editor-width: 100%;
}

:deep(.el-popconfirm__action) {
  margin-top: 8px;
}

:deep(.el-table-column--selection .cell) {
  padding-left: 14px;
}
</style> 