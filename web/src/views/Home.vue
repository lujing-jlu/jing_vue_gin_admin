<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900 mb-2">仪表盘</h2>
      <p class="text-gray-600">系统概览和统计数据</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-blue-500 to-blue-600">
            <el-icon class="text-2xl text-white"><User /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">用户总数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.userCount || 0 }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <router-link to="/users" class="hover:underline">查看详情</router-link>
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-green-500 to-green-600">
            <el-icon class="text-2xl text-white"><Setting /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">角色数量</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.roleCount || 0 }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <router-link to="/roles" class="hover:underline">查看详情</router-link>
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-purple-500 to-purple-600">
            <el-icon class="text-2xl text-white"><Document /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">日志总数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.logCount || 0 }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <router-link to="/logs" class="hover:underline">查看详情</router-link>
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-orange-500 to-orange-600">
            <el-icon class="text-2xl text-white"><Bell /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">今日登录</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.todayLoginCount || 0 }}</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <el-icon class="mr-1"><ArrowUp /></el-icon> 近24小时
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 左侧卡片 - 最近登录记录 -->
      <el-card class="border-0 shadow-sm rounded-xl col-span-2">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">最近登录记录</h3>
            <router-link to="/logs">
              <el-button type="primary" link class="text-blue-600 hover:text-blue-700">查看更多</el-button>
            </router-link>
          </div>
        </template>
        <div class="space-y-4" v-loading="loading">
          <div v-for="(log, index) in recentLogins" :key="log.id" class="flex items-start py-3 border-b border-gray-100 last:border-0">
            <el-avatar :size="40" class="bg-gradient-to-r from-blue-500 to-indigo-500">
              {{ log.username?.[0] || 'U' }}
            </el-avatar>
            <div class="ml-4 flex-1">
              <div class="flex justify-between items-start">
                <p class="text-sm text-gray-900">
                  <span class="font-medium">{{ log.username }}</span>
                  <span class="text-gray-500">登录系统</span>
                </p>
                <div>
                  <el-tag 
                    size="small" 
                    class="ml-2 rounded-full text-xs" 
                    :type="log.status === 1 ? 'success' : 'danger'"
                    effect="plain">
                    {{ log.status === 1 ? '成功' : '失败' }}
                  </el-tag>
                  <el-tag size="small" class="ml-2 rounded-full text-xs" type="info">
                    {{ formatTime(log.created_at) }}
                  </el-tag>
                </div>
              </div>
              <p class="text-xs text-gray-500 mt-1 flex items-center">
                <span class="mr-2">IP: {{ log.ip }}</span>
                <span>{{ log.detail }}</span>
              </p>
            </div>
          </div>
          <div v-if="recentLogins.length === 0 && !loading" class="py-8 text-center text-gray-500">
            暂无登录记录
          </div>
        </div>
      </el-card>

      <!-- 右侧卡片 - 系统信息 -->
      <el-card class="border-0 shadow-sm rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">系统信息</h3>
            <el-button 
              type="primary" 
              link 
              class="text-blue-600 hover:text-blue-700"
              @click="refreshSystemInfo">
              <el-icon class="mr-1"><Refresh /></el-icon>刷新
            </el-button>
          </div>
        </template>
        <div class="space-y-4">
          <div class="mb-4 last:mb-0">
            <div class="flex py-2 border-b border-gray-100">
              <div class="w-1/3 text-gray-500">系统名称</div>
              <div class="w-2/3">Vue Gin Admin</div>
            </div>
            <div class="flex py-2 border-b border-gray-100">
              <div class="w-1/3 text-gray-500">系统版本</div>
              <div class="w-2/3">v1.0.0</div>
            </div>
            <div class="flex py-2 border-b border-gray-100">
              <div class="w-1/3 text-gray-500">前端框架</div>
              <div class="w-2/3">Vue 3 + TypeScript</div>
            </div>
            <div class="flex py-2 border-b border-gray-100">
              <div class="w-1/3 text-gray-500">后端框架</div>
              <div class="w-2/3">Gin + GORM</div>
            </div>
            <div class="flex py-2 border-b border-gray-100">
              <div class="w-1/3 text-gray-500">服务器时间</div>
              <div class="w-2/3">{{ systemInfo.currentTime }}</div>
            </div>
            <div class="flex py-2">
              <div class="w-1/3 text-gray-500">运行环境</div>
              <div class="w-2/3">{{ systemInfo.environment }}</div>
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 底部卡片 - 最近操作记录 -->
    <el-card class="border-0 shadow-sm rounded-xl">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-medium text-gray-900">最近操作记录</h3>
          <router-link to="/logs">
            <el-button type="primary" link class="text-blue-600 hover:text-blue-700">查看更多</el-button>
          </router-link>
        </div>
      </template>
      <div v-loading="loading">
        <el-table :data="recentLogs" style="width: 100%" border="true">
          <el-table-column prop="username" label="操作用户" width="120" />
          <el-table-column prop="module" label="操作模块" width="120">
            <template #default="{ row }">
              <el-tag :type="getModuleType(row.module)" class="rounded-full px-2 py-0" size="small" effect="plain">
                {{ getModuleLabel(row.module) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="action" label="操作类型" width="120">
            <template #default="{ row }">
              <el-tag :type="getActionType(row.action)" class="rounded-full px-2 py-0" size="small" effect="plain">
                {{ getActionLabel(row.action) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="resource" label="操作资源" min-width="150" />
          <el-table-column prop="detail" label="操作详情" min-width="200" show-overflow-tooltip />
          <el-table-column prop="created_at" label="操作时间" width="170" />
        </el-table>
        <div v-if="recentLogs.length === 0 && !loading" class="py-8 text-center text-gray-500">
          暂无操作记录
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { User, Document, View, Star, Bell, ArrowUp, Refresh, Setting } from '@element-plus/icons-vue'
import { getLogList, type SystemLog } from '@/api/log'

const loading = ref(false)
const recentLogins = ref<SystemLog[]>([])
const recentLogs = ref<SystemLog[]>([])

const stats = reactive({
  userCount: 5,  // 模拟数据
  roleCount: 3,  // 模拟数据
  logCount: 0,
  todayLoginCount: 0
})

const systemInfo = reactive({
  currentTime: new Date().toLocaleString(),
  environment: 'Production'
})

// 获取最近登录记录
const fetchRecentLogins = async () => {
  try {
    loading.value = true
    const res = await getLogList({
      module: 'auth',
      action: 'login',
      page: 1,
      page_size: 5
    })
    recentLogins.value = res.list
    stats.todayLoginCount = res.list.filter((log: SystemLog) => {
      const today = new Date().toISOString().split('T')[0]
      return log.created_at.startsWith(today) && log.status === 1
    }).length
  } catch (error) {
    console.error('Failed to fetch recent logins:', error)
  } finally {
    loading.value = false
  }
}

// 获取最近操作记录
const fetchRecentLogs = async () => {
  try {
    loading.value = true
    const res = await getLogList({
      page: 1,
      page_size: 5
    })
    recentLogs.value = res.list
    stats.logCount = res.total
  } catch (error) {
    console.error('Failed to fetch recent logs:', error)
  } finally {
    loading.value = false
  }
}

// 格式化时间显示
const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  // 1分钟内
  if (diff < 60 * 1000) {
    return '刚刚'
  }
  // 1小时内
  if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))}分钟前`
  }
  // 24小时内
  if (diff < 24 * 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 60 * 1000))}小时前`
  }
  // 超过24小时
  return timeStr.substr(0, 16).replace('T', ' ')
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

// 刷新系统信息
const refreshSystemInfo = () => {
  systemInfo.currentTime = new Date().toLocaleString()
}

onMounted(() => {
  fetchRecentLogins()
  fetchRecentLogs()
})
</script>

<style scoped>
:deep(.el-card) {
  border-radius: 1rem;
  overflow: hidden;
}

:deep(.el-card__header) {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #f3f4f6;
  font-weight: 500;
}

:deep(.el-card__body) {
  padding: 1.5rem;
}

:deep(.el-button--primary.is-link) {
  color: #3b82f6;
}

:deep(.el-button--primary.is-link:hover) {
  color: #2563eb;
}

:deep(.el-tag--info) {
  background-color: #f3f4f6;
  border-color: #f3f4f6;
  color: #6b7280;
}

:deep(.el-tag--success.is-plain) {
  background-color: #dcfce7;
  border-color: #dcfce7;
  color: #16a34a;
}

:deep(.el-tag--danger.is-plain) {
  background-color: #fee2e2;
  border-color: #fee2e2;
  color: #dc2626;
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

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;  
  overflow: hidden;
}
</style> 