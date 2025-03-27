<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">系统设置</h2>
        <p class="text-gray-600">配置系统运行参数和功能选项</p>
      </div>
      <el-button 
        type="primary" 
        @click="saveSettings"
        :loading="saving"
        class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0"
      >
        <el-icon class="mr-2"><Check /></el-icon>
        保存设置
      </el-button>
    </div>

    <el-tabs type="border-card" class="settings-tabs">
      <el-tab-pane>
        <template #label>
          <div class="flex items-center">
            <el-icon class="mr-2"><Setting /></el-icon>
            基本设置
          </div>
        </template>
        <el-form :model="basicForm" label-position="top" class="p-4">
          <el-form-item label="系统名称">
            <el-input v-model="basicForm.siteTitle" />
          </el-form-item>
          <el-form-item label="系统描述">
            <el-input v-model="basicForm.siteDescription" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="Logo URL">
            <el-input v-model="basicForm.logoUrl" />
          </el-form-item>
          <el-form-item label="系统语言">
            <el-select v-model="basicForm.language" class="w-full">
              <el-option label="简体中文" value="zh-CN" />
              <el-option label="English" value="en-US" />
            </el-select>
          </el-form-item>
          <el-form-item label="时区设置">
            <el-select v-model="basicForm.timezone" class="w-full">
              <el-option label="北京时间 (UTC+8)" value="Asia/Shanghai" />
              <el-option label="UTC" value="UTC" />
              <el-option label="美国东部时间 (UTC-5)" value="America/New_York" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane>
        <template #label>
          <div class="flex items-center">
            <el-icon class="mr-2"><Lock /></el-icon>
            安全设置
          </div>
        </template>
        <el-form :model="securityForm" label-position="top" class="p-4">
          <el-form-item label="登录尝试次数限制">
            <el-input-number v-model="securityForm.loginAttempts" :min="1" :max="10" class="w-full" />
          </el-form-item>
          <el-form-item label="会话超时时间（分钟）">
            <el-input-number v-model="securityForm.sessionTimeout" :min="5" :max="1440" class="w-full" />
          </el-form-item>
          <el-form-item label="密码安全策略">
            <el-checkbox-group v-model="securityForm.passwordPolicy" class="flex flex-col space-y-2">
              <el-checkbox label="minLength">要求密码长度至少8位</el-checkbox>
              <el-checkbox label="requireNumber">要求包含数字</el-checkbox>
              <el-checkbox label="requireSpecial">要求包含特殊字符</el-checkbox>
              <el-checkbox label="requireUppercase">要求包含大写字母</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item>
            <el-switch
              v-model="securityForm.enableTwoFactor"
              active-text="启用双因素认证"
              class="w-full"
            />
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane>
        <template #label>
          <div class="flex items-center">
            <el-icon class="mr-2"><Message /></el-icon>
            通知设置
          </div>
        </template>
        <el-form :model="notificationForm" label-position="top" class="p-4">
          <el-form-item label="通知类型">
            <el-checkbox-group v-model="notificationForm.enabledNotifications" class="flex flex-col space-y-2">
              <el-checkbox label="email">邮件通知</el-checkbox>
              <el-checkbox label="browser">浏览器通知</el-checkbox>
              <el-checkbox label="sms">短信通知</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item label="系统维护提前通知时间（小时）">
            <el-input-number v-model="notificationForm.maintenanceNoticeHours" :min="1" :max="72" class="w-full" />
          </el-form-item>
          <el-form-item label="邮件服务器设置">
            <el-card class="border border-gray-200 shadow-none">
              <el-form-item label="SMTP 服务器">
                <el-input v-model="notificationForm.smtp.server" />
              </el-form-item>
              <el-form-item label="SMTP 端口">
                <el-input v-model="notificationForm.smtp.port" />
              </el-form-item>
              <el-form-item label="发件人邮箱">
                <el-input v-model="notificationForm.smtp.username" />
              </el-form-item>
              <el-form-item label="密码/授权码">
                <el-input v-model="notificationForm.smtp.password" type="password" show-password />
              </el-form-item>
              <el-form-item>
                <el-switch
                  v-model="notificationForm.smtp.ssl"
                  active-text="使用 SSL 加密"
                  class="w-full"
                />
              </el-form-item>
              <div class="flex justify-end">
                <el-button size="small" @click="testEmailConfig">测试邮件连接</el-button>
              </div>
            </el-card>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      
      <el-tab-pane>
        <template #label>
          <div class="flex items-center">
            <el-icon class="mr-2"><DataLine /></el-icon>
            数据设置
          </div>
        </template>
        <el-form :model="dataForm" label-position="top" class="p-4">
          <el-form-item label="数据库类型">
            <el-radio-group v-model="dataForm.dbType">
              <el-radio-button label="mysql">MySQL</el-radio-button>
              <el-radio-button label="postgres">PostgreSQL</el-radio-button>
              <el-radio-button label="sqlite">SQLite</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="数据库连接字符串" v-if="dataForm.dbType !== 'sqlite'">
            <el-input v-model="dataForm.connectionString" />
          </el-form-item>
          <el-form-item label="SQLite 文件路径" v-else>
            <el-input v-model="dataForm.sqlitePath" />
          </el-form-item>
          <el-divider content-position="left">数据备份</el-divider>
          <el-form-item label="自动备份">
            <el-switch
              v-model="dataForm.autoBackup"
              active-text="启用自动备份"
              class="w-full mb-2"
            />
          </el-form-item>
          <el-form-item label="备份频率" v-if="dataForm.autoBackup">
            <el-select v-model="dataForm.backupFrequency" class="w-full">
              <el-option label="每天" value="daily" />
              <el-option label="每周" value="weekly" />
              <el-option label="每月" value="monthly" />
            </el-select>
          </el-form-item>
          <el-form-item label="保留备份数量" v-if="dataForm.autoBackup">
            <el-input-number v-model="dataForm.backupRetention" :min="1" :max="30" class="w-full" />
          </el-form-item>
          <div class="flex justify-between">
            <el-button @click="createBackup" :loading="backingUp">
              <el-icon class="mr-1"><Download /></el-icon>
              创建备份
            </el-button>
            <el-button @click="restoreBackup">
              <el-icon class="mr-1"><Upload /></el-icon>
              恢复备份
            </el-button>
          </div>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { 
  Setting, 
  Lock, 
  Message, 
  DataLine, 
  Check, 
  Download, 
  Upload 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const saving = ref(false)
const backingUp = ref(false)

// 基本设置表单
const basicForm = reactive({
  siteTitle: 'Vue Admin 系统',
  siteDescription: '一个现代化的管理系统框架，提供用户管理、数据分析等功能',
  logoUrl: '/logo.png',
  language: 'zh-CN',
  timezone: 'Asia/Shanghai'
})

// 安全设置表单
const securityForm = reactive({
  loginAttempts: 5,
  sessionTimeout: 60,
  passwordPolicy: ['minLength', 'requireNumber'],
  enableTwoFactor: false
})

// 通知设置表单
const notificationForm = reactive({
  enabledNotifications: ['email', 'browser'],
  maintenanceNoticeHours: 24,
  smtp: {
    server: 'smtp.example.com',
    port: '587',
    username: 'admin@example.com',
    password: '',
    ssl: true
  }
})

// 数据设置表单
const dataForm = reactive({
  dbType: 'sqlite',
  connectionString: '',
  sqlitePath: './data.sqlite',
  autoBackup: true,
  backupFrequency: 'daily',
  backupRetention: 7
})

// 保存设置
const saveSettings = async () => {
  try {
    saving.value = true
    // 这里添加保存配置的 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('设置保存成功')
  } catch (error) {
    ElMessage.error('保存设置失败')
    console.error('保存设置失败:', error)
  } finally {
    saving.value = false
  }
}

// 测试邮件配置
const testEmailConfig = async () => {
  try {
    // 这里添加测试邮件配置的 API 调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('邮件服务器连接测试成功')
  } catch (error) {
    ElMessage.error('邮件服务器连接测试失败')
    console.error('邮件配置测试失败:', error)
  }
}

// 创建备份
const createBackup = async () => {
  try {
    backingUp.value = true
    // 这里添加创建备份的 API 调用
    await new Promise(resolve => setTimeout(resolve, 1500))
    ElMessage.success('数据备份创建成功')
  } catch (error) {
    ElMessage.error('数据备份创建失败')
    console.error('数据备份失败:', error)
  } finally {
    backingUp.value = false
  }
}

// 恢复备份
const restoreBackup = async () => {
  try {
    // 这里添加恢复备份的 API 调用
    ElMessage.warning('该功能尚未实现')
  } catch (error) {
    ElMessage.error('恢复备份失败')
    console.error('恢复备份失败:', error)
  }
}
</script>

<style scoped>
.settings-tabs {
  border-radius: 0.75rem;
  overflow: hidden;
  border: none;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

:deep(.el-tabs__header) {
  background-color: #f9fafb;
  margin: 0;
}

:deep(.el-tabs__item) {
  height: 50px;
  line-height: 50px;
  font-weight: 500;
}

:deep(.el-tabs__item.is-active) {
  color: #3b82f6;
}

:deep(.el-tabs__active-bar) {
  background-color: #3b82f6;
}

:deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #f3f4f6;
}

:deep(.el-tabs__content) {
  background-color: #ffffff;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
}

:deep(.el-input__wrapper) {
  border-radius: 0.5rem;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
}

:deep(.el-select__wrapper) {
  border-radius: 0.5rem;
}

:deep(.el-input-number__decrease),
:deep(.el-input-number__increase) {
  border-radius: 0.5rem;
}

:deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-radius: 0.5rem 0 0 0.5rem;
}

:deep(.el-radio-button:last-child .el-radio-button__inner) {
  border-radius: 0 0.5rem 0.5rem 0;
}

:deep(.el-checkbox__label) {
  color: #4b5563;
}

:deep(.el-switch__core) {
  border-radius: 999px;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #4b5563;
}

:deep(.el-button) {
  border-radius: 0.5rem;
  transition: all 0.2s;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style> 