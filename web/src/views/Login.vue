<template>
  <div class="min-h-screen flex items-stretch">
    <!-- 左侧登录表单 -->
    <div class="w-full md:w-2/5 flex items-center justify-center px-6 lg:px-12 xl:px-16 bg-gray-50">
      <div class="w-full max-w-md">
        <div class="mb-10">
          <h2 class="text-3xl font-bold text-gray-900 mb-2">欢迎回来</h2>
          <p class="text-gray-600">请登录您的账号继续访问系统</p>
        </div>
        <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin" class="space-y-6">
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="用户名"
              :prefix-icon="User"
              class="!h-12"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码"
              :prefix-icon="Lock"
              show-password
              class="!h-12"
            />
          </el-form-item>
          <div class="flex justify-between items-center text-sm mb-4">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="#" class="text-blue-600 hover:text-blue-800">忘记密码?</a>
          </div>
          <el-form-item>
            <el-button 
              type="primary" 
              class="w-full !h-12 text-base font-medium bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0"
              :loading="loading"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        <div class="mt-8 text-center text-sm text-gray-600">
          <p>© {{ new Date().getFullYear() }} Vue Admin. 保留所有权利。</p>
        </div>
      </div>
    </div>
    
    <!-- 右侧图片 -->
    <div class="hidden md:block md:w-3/5 bg-gradient-to-br from-blue-500 to-indigo-700 relative overflow-hidden">
      <div class="absolute inset-0 flex items-center justify-center p-10">
        <div class="max-w-lg text-white">
          <h2 class="text-4xl font-bold mb-4">现代化管理系统</h2>
          <p class="mb-8 text-blue-100 text-lg">提供强大的用户管理、内容管理和数据分析功能，助力业务高效运转。</p>
          <div class="grid grid-cols-2 gap-6">
            <div class="bg-white/10 backdrop-blur-sm p-6 rounded-xl">
              <div class="text-3xl font-bold mb-1">10k+</div>
              <div class="text-blue-100 text-sm">活跃用户</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm p-6 rounded-xl">
              <div class="text-3xl font-bold mb-1">5M+</div>
              <div class="text-blue-100 text-sm">数据记录</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm p-6 rounded-xl">
              <div class="text-3xl font-bold mb-1">98%</div>
              <div class="text-blue-100 text-sm">用户满意度</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm p-6 rounded-xl">
              <div class="text-3xl font-bold mb-1">24/7</div>
              <div class="text-blue-100 text-sm">技术支持</div>
            </div>
          </div>
        </div>
      </div>
      <div class="absolute bottom-0 left-0 right-0 h-1/3 bg-gradient-to-t from-indigo-900/50 to-transparent"></div>
      <div class="absolute -bottom-40 -right-40 w-80 h-80 bg-blue-400 rounded-full opacity-20"></div>
      <div class="absolute -top-20 -left-20 w-60 h-60 bg-indigo-300 rounded-full opacity-20"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        await userStore.loginAction(form)
        ElMessage.success('登录成功')
      } catch (error) {
        console.error('Login failed:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
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
  height: 3rem;
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

:deep(.el-checkbox__label) {
  color: #6b7280;
}

:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: #3b82f6;
  border-color: #3b82f6;
}

:deep(.el-checkbox__input.is-focus .el-checkbox__inner) {
  border-color: #3b82f6;
}
</style> 