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
            <p class="text-sm text-gray-600">总用户数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">1,234</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <el-icon class="mr-1"><ArrowUp /></el-icon> 12% 较上月
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-green-500 to-green-600">
            <el-icon class="text-2xl text-white"><Document /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">文章数量</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">567</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <el-icon class="mr-1"><ArrowUp /></el-icon> 8% 较上月
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-purple-500 to-purple-600">
            <el-icon class="text-2xl text-white"><View /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">总浏览量</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">89,012</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <el-icon class="mr-1"><ArrowUp /></el-icon> 23% 较上月
            </p>
          </div>
        </div>
      </el-card>

      <el-card class="border-0 shadow-sm hover:shadow-md transition-shadow rounded-xl">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-orange-500 to-orange-600">
            <el-icon class="text-2xl text-white"><Star /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">收藏数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">3,456</p>
            <p class="text-xs text-green-600 mt-1 flex items-center">
              <el-icon class="mr-1"><ArrowUp /></el-icon> 15% 较上月
            </p>
          </div>
        </div>
      </el-card>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 左侧卡片 - 最近活动 -->
      <el-card class="border-0 shadow-sm rounded-xl col-span-2">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">最近活动</h3>
            <el-button type="primary" link>查看更多</el-button>
          </div>
        </template>
        <div class="space-y-4">
          <div v-for="i in 5" :key="i" class="flex items-start py-3 border-b border-gray-100 last:border-0">
            <el-avatar :size="40" class="bg-gradient-to-r from-blue-500 to-indigo-500">
              {{ ['A', 'B', 'C', 'D', 'E'][i-1] }}
            </el-avatar>
            <div class="ml-4 flex-1">
              <div class="flex justify-between items-start">
                <p class="text-sm text-gray-900">
                  <span class="font-medium">用户{{ i }}</span>
                  发表了新文章
                  <span class="text-blue-600 cursor-pointer hover:underline">《示例文章标题》</span>
                </p>
                <el-tag size="small" class="ml-2 rounded-full text-xs" type="info">2小时前</el-tag>
              </div>
              <p class="text-xs text-gray-500 mt-1 line-clamp-2">这是一篇示例文章的摘要内容，展示在动态列表中。用户可以通过点击标题查看详情...</p>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 右侧卡片 - 系统公告 -->
      <el-card class="border-0 shadow-sm rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-medium text-gray-900">系统公告</h3>
            <el-button type="primary" link>更多</el-button>
          </div>
        </template>
        <div class="space-y-4">
          <div v-for="i in 3" :key="i" class="mb-4 last:mb-0">
            <div class="flex items-center justify-between mb-2">
              <h4 class="text-base text-gray-900 font-medium flex items-center">
                <el-icon class="mr-2 text-blue-600"><Bell /></el-icon>
                系统{{ i === 1 ? '更新' : i === 2 ? '维护' : '安全' }}通知
              </h4>
              <span class="text-xs text-gray-500">{{ i }}天前</span>
            </div>
            <p class="text-sm text-gray-600 mb-3">
              {{ i === 1 ? '系统已升级到最新版本，新增多项功能和安全修复。' : 
                 i === 2 ? '系统将于今晚22:00进行例行维护，预计持续1小时。' : 
                 '为保障账号安全，建议定期修改密码并启用双因素认证。' }}
            </p>
            <el-divider class="!my-1" />
          </div>
        </div>
      </el-card>
    </div>

    <!-- 底部卡片 -->
    <el-card class="border-0 shadow-sm rounded-xl">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-medium text-gray-900">系统使用情况</h3>
          <div class="flex space-x-2">
            <el-select v-model="timeRange" placeholder="Select" size="small" class="!w-32">
              <el-option label="本周" value="week" />
              <el-option label="本月" value="month" />
              <el-option label="本季度" value="quarter" />
            </el-select>
            <el-button type="primary" size="small" class="bg-blue-600 text-white hover:bg-blue-700">
              <el-icon class="mr-1"><Refresh /></el-icon>刷新
            </el-button>
          </div>
        </div>
      </template>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="flex flex-col">
          <h4 class="text-base font-medium text-gray-900 mb-4">活跃用户趋势</h4>
          <div class="flex-1 h-64 bg-gradient-to-b from-blue-50 to-white rounded-lg flex items-center justify-center">
            <p class="text-gray-400">图表区域</p>
          </div>
        </div>
        <div class="flex flex-col">
          <h4 class="text-base font-medium text-gray-900 mb-4">内容分布情况</h4>
          <div class="flex-1 h-64 bg-gradient-to-b from-blue-50 to-white rounded-lg flex items-center justify-center">
            <p class="text-gray-400">图表区域</p>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { User, Document, View, Star, Bell, ArrowUp, Refresh } from '@element-plus/icons-vue'

const timeRange = ref('month')
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

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;  
  overflow: hidden;
}
</style> 