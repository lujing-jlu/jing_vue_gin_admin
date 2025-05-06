<template>
  <el-container class="h-full w-full">
    <el-aside width="240px" class="bg-gradient-to-b from-gray-900 to-gray-800 shadow-xl">
      <div class="h-16 flex items-center justify-center border-b border-gray-700">
        <h2 class="text-xl font-bold text-white">Vue Admin</h2>
      </div>
      <el-scrollbar class="h-[calc(100vh-4rem)]">
        <el-menu
          class="border-0"
          background-color="transparent"
          text-color="#9ca3af"
          active-text-color="#ffffff"
          router
        >
          <div class="px-4 py-3 text-xs uppercase text-gray-500 font-semibold">
            主要功能
          </div>
          <el-menu-item index="/" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/users" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/roles" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><Setting /></el-icon>
            <span>角色管理</span>
          </el-menu-item>
          <el-menu-item index="/files" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><Files /></el-icon>
            <span>文件管理</span>
          </el-menu-item>
          
          <div class="px-4 py-3 text-xs uppercase text-gray-500 font-semibold mt-4">
            系统管理
          </div>
          <el-menu-item index="/logs" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><Document /></el-icon>
            <span>系统日志</span>
          </el-menu-item>
          <el-menu-item index="/settings" class="hover:bg-gray-800 transition-colors" style="color: #ffffff !important">
            <el-icon><Tools /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
      </el-scrollbar>
    </el-aside>
    <el-container class="bg-gray-50">
      <el-header class="bg-white border-b h-16 flex items-center justify-between px-6 shadow-sm">
        <div class="flex items-center">
          <el-icon class="text-xl cursor-pointer mr-4 text-gray-600 hover:text-gray-900 transition-colors"><Fold /></el-icon>
          <el-breadcrumb>
            <el-breadcrumb-item>首页</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="flex items-center space-x-4">
          <el-badge :value="3" class="cursor-pointer">
            <el-icon class="text-xl text-gray-600"><Bell /></el-icon>
          </el-badge>
          <el-dropdown @command="handleCommand" trigger="click">
            <span class="flex items-center cursor-pointer hover:bg-gray-50 px-3 py-2 rounded-lg transition-colors text-gray-700" style="color: #374151 !important">
              <el-avatar :size="32" class="mr-2 bg-gradient-to-r from-blue-500 to-indigo-500">
                {{ userStore.userInfo.nickname?.[0] || 'A' }}
              </el-avatar>
              <span class="text-gray-700" style="color: #374151 !important">{{ userStore.userInfo.nickname || 'Admin' }}</span>
              <el-icon class="ml-2 text-gray-500"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu class="rounded-lg shadow-lg border-0">
                <el-dropdown-item command="profile" class="hover:bg-gray-50 text-gray-700" style="color: #374151 !important">
                  <el-icon class="mr-2"><User /></el-icon>个人信息
                </el-dropdown-item>
                <el-dropdown-item command="settings" class="hover:bg-gray-50 text-gray-700" style="color: #374151 !important">
                  <el-icon class="mr-2"><Setting /></el-icon>账号设置
                </el-dropdown-item>
                <el-divider class="!my-1" />
                <el-dropdown-item command="logout" class="hover:bg-gray-50 text-gray-700" style="color: #374151 !important">
                  <el-icon class="mr-2"><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="p-6 h-[calc(100vh-4rem)] overflow-auto">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { HomeFilled, ArrowDown, Fold, User, SwitchButton, Setting, Bell, Tools, Document, Files } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

const handleCommand = (command: string) => {
  if (command === 'logout') {
    handleLogout()
  } else if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'settings') {
    router.push('/settings')
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    customClass: 'rounded-lg'
  }).then(() => {
    userStore.logout()
    ElMessage.success('退出成功')
  }).catch(() => {})
}
</script>

<style scoped>
:deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
  margin: 4px 0;
  border-radius: 0.5rem;
  margin-left: 8px;
  margin-right: 8px;
}

:deep(.el-menu-item.is-active) {
  background-color: rgba(59, 130, 246, 0.2) !important;
  color: #ffffff !important;
  font-weight: 500;
}

:deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.05) !important;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 16px;
  display: flex;
  align-items: center;
}

:deep(.el-breadcrumb__item) {
  color: #6b7280;
}

:deep(.el-breadcrumb__inner.is-link:hover) {
  color: #3b82f6;
}

:deep(.el-icon) {
  vertical-align: middle;
}

:deep(.el-aside) {
  transition: width 0.3s;
}

:deep(.el-header) {
  padding: 0 1.5rem;
}

:deep(.el-main) {
  padding: 1.5rem;
}
</style> 