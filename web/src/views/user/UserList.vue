<template>
  <div class="p-4">
    <div class="mb-4 flex justify-between items-center">
      <h2 class="text-xl font-bold">用户管理</h2>
      <el-button type="primary" @click="handleAdd">
        <el-icon class="mr-2"><Plus /></el-icon>
        添加用户
      </el-button>
    </div>

    <el-table :data="userList" border stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">
            {{ row.role === 'admin' ? '管理员' : '普通用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button 
              :type="row.status === 1 ? 'danger' : 'success'" 
              size="small" 
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加用户' : '编辑用户'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="dialogType === 'add'">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, updateUser, toggleUserStatus, type User, type CreateUserRequest, type UpdateUserRequest } from '@/api/user'

const userList = ref<User[]>([])
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()

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
        type: 'warning'
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