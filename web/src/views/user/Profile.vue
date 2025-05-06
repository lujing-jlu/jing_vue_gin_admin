<template>
  <div class="p-4">
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">个人中心</h2>
      <p class="text-gray-600">查看和管理您的个人信息</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- 左侧个人信息卡片 -->
      <div class="md:col-span-1">
        <el-card class="border-0 shadow-sm rounded-xl overflow-hidden h-full">
          <div class="flex flex-col items-center p-6">
            <el-avatar :size="100" :src="userInfo.avatar || defaultAvatar" class="mb-4 border-4 border-white shadow-md">
              {{ userInfo.nickname?.[0] || 'U' }}
            </el-avatar>
            <h3 class="text-xl font-bold text-gray-800 mb-1">{{ userInfo.nickname }}</h3>
            <p class="text-gray-500 mb-4">{{ userInfo.username }}</p>
            <el-tag class="rounded-full px-4 py-1" :type="userInfo.role === 'admin' ? 'danger' : 'primary'" effect="plain">
              {{ roleLabel }}
            </el-tag>
            <div class="w-full mt-6 pt-6 border-t border-gray-100">
              <div class="flex justify-between items-center mb-4">
                <span class="text-gray-500">状态</span>
                <el-tag 
                  :type="userInfo.status === 1 ? 'success' : 'danger'" 
                  class="rounded-full px-3 py-1"
                  effect="plain"
                >
                  {{ userInfo.status === 1 ? '正常' : '禁用' }}
                </el-tag>
              </div>
              <div class="flex justify-between items-center mb-4">
                <span class="text-gray-500">创建时间</span>
                <span class="text-gray-700">{{ formatDate(userInfo.created_at) }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-gray-500">上次更新</span>
                <span class="text-gray-700">{{ formatDate(userInfo.updated_at) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 右侧信息编辑区 -->
      <div class="md:col-span-2">
        <el-card class="border-0 shadow-sm rounded-xl overflow-hidden">
          <template #header>
            <div class="flex justify-between items-center">
              <h3 class="text-lg font-medium text-gray-900">个人资料设置</h3>
            </div>
          </template>
          <el-form 
            ref="formRef"
            :model="form"
            :rules="rules"
            label-position="top"
            class="p-2"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" disabled class="!h-10" />
              <div class="text-gray-500 text-xs mt-1">用户名不可修改</div>
            </el-form-item>
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" class="!h-10" />
            </el-form-item>
            <el-form-item label="头像" prop="avatar">
              <div class="flex items-center">
                <el-avatar :size="60" :src="form.avatar || defaultAvatar" class="mr-4">
                  {{ form.nickname?.[0] || 'U' }}
                </el-avatar>
                <el-upload
                  class="avatar-uploader"
                  accept="image/*"
                  :show-file-list="false"
                  :before-upload="handleAvatarUpload"
                >
                  <el-button class="text-gray-600">
                    <el-icon class="mr-1"><Upload /></el-icon>
                    上传头像
                  </el-button>
                </el-upload>
              </div>
            </el-form-item>
            <el-divider content-position="left">安全设置</el-divider>
            <el-form-item>
              <el-button 
                type="primary" 
                @click="showPasswordDialog = true"
                class="bg-blue-600 text-white hover:bg-blue-700"
              >
                <el-icon class="mr-1"><Lock /></el-icon>
                修改密码
              </el-button>
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary"
                @click="handleSubmit"
                :loading="submitting"
                class="bg-blue-600 text-white hover:bg-blue-700"
              >
                <el-icon class="mr-1"><Check /></el-icon>
                保存修改
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>
    </div>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="500px"
      class="rounded-lg"
      destroy-on-close
    >
      <el-form
        ref="pwdFormRef"
        :model="pwdForm"
        :rules="pwdRules"
        label-position="top"
      >
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input 
            v-model="pwdForm.oldPassword"
            type="password" 
            show-password
            class="!h-10"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="pwdForm.newPassword"
            type="password" 
            show-password
            class="!h-10"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="pwdForm.confirmPassword"
            type="password" 
            show-password
            class="!h-10"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <el-button @click="showPasswordDialog = false" class="text-gray-600">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleChangePassword"
            :loading="changingPassword"
            class="bg-blue-600 text-white hover:bg-blue-700"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Upload, Lock, Check } from '@element-plus/icons-vue'
import type { FormInstance, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateProfile, changePassword } from '@/api/user'
import { uploadFile } from '@/api/file'
import type { ChangePasswordRequest } from '@/api/user'

// 默认头像
const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

// 用户状态
const userStore = useUserStore()
const userInfo = ref(userStore.userInfo)
const roleLabel = ref(userInfo.value.role === 'admin' ? '管理员' : '普通用户')

// 表单相关
const formRef = ref<FormInstance>()
const form = reactive({
  username: userInfo.value.username || '',
  nickname: userInfo.value.nickname || '',
  avatar: userInfo.value.avatar || ''
})

// 提交状态
const submitting = ref(false)

// 表单验证规则
const rules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

// 密码修改相关
const pwdFormRef = ref<FormInstance>()
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码表单验证规则
const pwdRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: any) => {
        if (value !== pwdForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 头像上传相关
const handleAvatarUpload: UploadProps['beforeUpload'] = async (file) => {
  // 验证文件类型
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }

  // 验证文件大小（2MB）
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB！')
    return false
  }

  try {
    const res = await uploadFile(file, 'avatar')
    form.avatar = res.url
    ElMessage.success('头像上传成功')
  } catch (error) {
    ElMessage.error('头像上传失败')
  }
  return false
}

// 格式化日期
const formatDate = (dateStr?: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 调用API更新个人资料
    await updateProfile({
      nickname: form.nickname,
      avatar: form.avatar
    })
    
    // 更新本地用户信息
    userStore.updateUserInfo({
      nickname: form.nickname,
      avatar: form.avatar
    })
    
    ElMessage.success('个人信息更新成功')
  } catch (error) {
    ElMessage.error('个人信息更新失败')
  } finally {
    submitting.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  if (!pwdFormRef.value) return
  
  try {
    await pwdFormRef.value.validate()
    changingPassword.value = true
    
    // 调用API修改密码
    const passwordData: ChangePasswordRequest = {
      oldPassword: pwdForm.oldPassword,
      newPassword: pwdForm.newPassword
    }
    await changePassword(passwordData)
    
    ElMessage.success('密码修改成功')
    showPasswordDialog.value = false
    pwdForm.oldPassword = ''
    pwdForm.newPassword = ''
    pwdForm.confirmPassword = ''
  } catch (error) {
    ElMessage.error('密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

// 初始化
onMounted(() => {
  // 可以在这里加载最新的用户信息
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

:deep(.el-button) {
  border-radius: 0.5rem;
  transition: all 0.2s;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

:deep(.el-dialog) {
  border-radius: 0.75rem;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  margin: 0;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #f3f4f6;
  font-weight: 600;
}

:deep(.el-dialog__body) {
  padding: 1.5rem;
}

:deep(.el-dialog__footer) {
  padding: 1rem 1.5rem;
  border-top: 1px solid #f3f4f6;
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

:deep(.el-tag--primary.is-plain) {
  background-color: #e0f2fe;
  border-color: #e0f2fe;
  color: #0ea5e9;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #4b5563;
}

:deep(.avatar-uploader) {
  display: inline-block;
}

:deep(.avatar-uploader .el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: all 0.3s;
}

:deep(.avatar-uploader .el-upload:hover) {
  border-color: #409eff;
}

:deep(.avatar-uploader .el-icon.avatar-uploader-icon) {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}
</style> 