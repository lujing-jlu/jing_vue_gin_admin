<template>
  <el-container class="h-full w-full">
    <el-aside width="200px" class="bg-gray-800">
      <div class="h-16 flex items-center justify-center border-b border-gray-700">
        <h2 class="text-xl font-bold text-white">Vue Admin</h2>
      </div>
      <el-menu
        class="h-[calc(100vh-4rem)] border-0"
        background-color="#1f2937"
        text-color="#fff"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="bg-white border-b h-16 flex items-center justify-between px-4">
        <div class="flex items-center">
          <el-icon class="text-xl cursor-pointer mr-4"><Fold /></el-icon>
          <el-breadcrumb>
            <el-breadcrumb-item>首页</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <el-dropdown @command="handleCommand">
          <span class="flex items-center cursor-pointer">
            <el-avatar :size="32" class="mr-2">{{ userStore.userInfo.nickname?.[0] || 'A' }}</el-avatar>
            {{ userStore.userInfo.nickname || 'Admin' }}
            <el-icon class="ml-2"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="bg-gray-100 p-4 h-[calc(100vh-4rem)] overflow-auto">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { HomeFilled, ArrowDown, Fold, User } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      ElMessage.info('个人信息功能开发中')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('退出成功')
      break
  }
}
</script> 