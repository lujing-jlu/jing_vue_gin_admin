<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">用户管理</h2>
        <p class="text-gray-600">管理系统中的所有用户账号</p>
      </div>
      <el-button 
        type="primary" 
        @click="handleAdd"
        class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0"
      >
        <el-icon class="mr-2"><Plus /></el-icon>
        添加用户
      </el-button>
    </div>

    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden">
      <div class="mb-4 flex flex-wrap gap-4">
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索用户名/昵称" 
          class="max-w-xs" 
          clearable
          prefix-icon="Search"
        />
        <el-select v-model="filterRole" placeholder="角色" class="w-32" clearable>
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
        </el-select>
        <el-select v-model="filterStatus" placeholder="状态" class="w-32" clearable>
          <el-option label="正常" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button type="primary" plain>
          <el-icon class="mr-1"><Search /></el-icon>
          搜索
        </el-button>
        <el-button>
          <el-icon class="mr-1"><Refresh /></el-icon>
          重置
        </el-button>
      </div>

      <el-table 
        :data="userList" 
        border 
        stripe 
        v-loading="loading"
        class="rounded-lg overflow-hidden"
        :header-cell-style="{ background: '#f9fafb', color: '#374151', fontWeight: '600' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="role" label="角色" width="120" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="row.role === 'admin' ? 'danger' : 'info'"
              class="rounded-full px-3 py-1"
              effect="plain"
            >
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="row.status === 1 ? 'success' : 'danger'"
              class="rounded-full px-3 py-1"
              effect="plain"
            >
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="180" />
        <el-table-column label="操作" width="240" fixed="right" align="center">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(row)"
                class="!rounded-l-lg"
              >
                <el-icon class="mr-1"><Edit /></el-icon>编辑
              </el-button>
              <el-button 
                :type="row.status === 1 ? 'danger' : 'success'" 
                size="small" 
                @click="handleToggleStatus(row)"
                class="!rounded-r-lg"
              >
                <el-icon class="mr-1">
                  <component :is="row.status === 1 ? 'Lock' : 'Unlock'" />
                </el-icon>
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="flex justify-end mt-4">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="100"
          background
          class="!p-0"
        />
      </div>
    </el-card>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加用户' : '编辑用户'"
      width="500px"
      class="rounded-lg user-dialog"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        class="space-y-4"
      >
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="form.username" 
            :disabled="dialogType === 'edit'"
            class="!h-10"
            placeholder="请输入用户名"
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input 
            v-model="form.nickname"
            class="!h-10"
            placeholder="请输入昵称"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="dialogType === 'add'">
          <el-input 
            v-model="form.password" 
            type="password" 
            show-password
            class="!h-10"
            placeholder="请输入密码"
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select 
            v-model="form.role"
            class="w-full !h-10"
            placeholder="请选择角色"
          >
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit" 
            :loading="submitting"
            class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { Plus, Edit, Lock, Unlock, Search, Refresh } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, updateUser, toggleUserStatus, type User, type CreateUserRequest, type UpdateUserRequest } from '@/api/user'

const userList = ref<User[]>([])
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const searchQuery = ref('')
const filterRole = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const form = reactive({
  id: 0,
  username: '',
  nickname: '',
  password: '',
  role: 'user',
  status: 1
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const fetchUserList = async () => {
  try {
    loading.value = true
    const res = await getUserList()
    userList.value = res
  } catch (error) {
    console.error('Failed to fetch user list:', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  form.id = 0
  form.username = ''
  form.nickname = ''
  form.password = ''
  form.role = 'user'
  form.status = 1
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleToggleStatus = async (row: User) => {
  const action = row.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 ${row.username} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        customClass: 'rounded-lg'
      }
    )
    await toggleUserStatus(row.id)
    ElMessage.success(`${action}成功`)
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to toggle user status:', error)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        submitting.value = true
        if (dialogType.value === 'add') {
          const createData: CreateUserRequest = {
            username: form.username,
            password: form.password,
            nickname: form.nickname,
            role: form.role
          }
          await createUser(createData)
        } else {
          const updateData: UpdateUserRequest = {
            nickname: form.nickname,
            role: form.role
          }
          await updateUser(form.id, updateData)
        }
        ElMessage.success(dialogType.value === 'add' ? '添加成功' : '编辑成功')
        dialogVisible.value = false
        fetchUserList()
      } catch (error) {
        console.error('Failed to submit form:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  fetchUserList()
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

:deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #3b82f6;
}

.user-dialog :deep(.el-overlay-dialog) {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style> 