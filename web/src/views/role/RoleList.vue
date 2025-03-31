<template>
  <div class="p-4">
    <div class="mb-4 flex justify-between items-center">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">角色管理</h2>
        <p class="text-gray-600">管理系统中的所有角色和权限</p>
      </div>
      <el-button 
        type="primary" 
        @click="handleAdd"
        class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0 text-white"
      >
        <el-icon class="mr-2"><Plus /></el-icon>
        添加角色
      </el-button>
    </div>

    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden">
      <!-- 角色列表 -->
      <el-table 
        :data="roles" 
        border 
        stripe 
        v-loading="loading"
        class="rounded-lg overflow-hidden"
        :header-cell-style="{ background: '#f9fafb', color: '#374151', fontWeight: '600' }"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="角色名称" min-width="120" />
        <el-table-column prop="description" label="描述" min-width="180" />
        <el-table-column label="权限" min-width="300">
          <template #default="{ row }">
            <el-tag
              v-for="perm in row.permissions"
              :key="perm"
              class="mr-1 mb-1 rounded-full px-3 py-1"
              effect="plain"
            >
              {{ getPermissionLabel(perm) }}
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
        <el-table-column label="操作" width="240" fixed="right" align="center">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                type="primary" 
                size="small" 
                @click="handleEdit(row)"
                class="!rounded-l-lg bg-blue-600 text-white hover:bg-blue-700"
              >
                <el-icon class="mr-1"><Edit /></el-icon>编辑
              </el-button>
              <el-button 
                :type="row.status === 1 ? 'danger' : 'success'" 
                size="small" 
                @click="handleToggleStatus(row)"
                class="!rounded-r-lg text-white"
                :class="row.status === 1 ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'"
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
    </el-card>

    <!-- 角色表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加角色' : '编辑角色'"
      width="600px"
      class="rounded-lg role-dialog"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        class="space-y-4"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input 
            v-model="form.name" 
            class="!h-10"
            placeholder="请输入角色名称"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description"
            type="textarea"
            :rows="2"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        <el-form-item label="权限" prop="permissions">
          <div class="border rounded-lg p-4 space-y-6">
            <div v-for="(perms, module) in permissions" :key="module" class="space-y-2">
              <div class="flex items-center justify-between mb-2">
                <div class="font-bold text-gray-700">{{ module }}</div>
                <el-checkbox
                  :model-value="isModuleSelected(module)"
                  :indeterminate="isModuleIndeterminate(module)"
                  @change="handleModuleCheckAll(module, $event)"
                >
                  全选
                </el-checkbox>
              </div>
              <div class="flex flex-wrap gap-4">
                <el-checkbox
                  v-for="perm in perms"
                  :key="perm.value"
                  v-model="form.permissions"
                  :label="perm.value"
                >
                  {{ perm.label }}
                </el-checkbox>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <el-button @click="dialogVisible = false" class="text-gray-600">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit" 
            :loading="submitting"
            class="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 border-0 text-white hover:shadow-lg"
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
import { Plus, Edit, Lock, Unlock } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  toggleRoleStatus,
  getAllPermissions,
  type Role,
  type CreateRoleRequest,
  type UpdateRoleRequest,
  type PermissionGroup
} from '@/api/role'

// 角色列表数据
const roles = ref<Role[]>([])
const loading = ref(false)
const submitting = ref(false)
const permissions = ref<PermissionGroup>({})

// 对话框相关
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const currentId = ref<number>()

// 表单相关
const formRef = ref<FormInstance>()
const form = ref<CreateRoleRequest>({
  name: '',
  description: '',
  permissions: []
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入角色描述', trigger: 'blur' }
  ],
  permissions: [
    { required: true, message: '请选择权限', trigger: 'change' }
  ]
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    loading.value = true
    roles.value = await getRoleList()
  } catch (error) {
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 获取所有权限
const fetchPermissions = async () => {
  try {
    permissions.value = await getAllPermissions()
  } catch (error) {
    ElMessage.error('获取权限列表失败')
  }
}

// 获取权限标签
const getPermissionLabel = (permValue: string) => {
  for (const module in permissions.value) {
    const perm = permissions.value[module].find(p => p.value === permValue)
    if (perm) {
      return perm.label
    }
  }
  return permValue
}

// 检查模块是否全选
const isModuleSelected = (module: string) => {
  const modulePerms = permissions.value[module]
  return modulePerms.every(perm => form.value.permissions.includes(perm.value))
}

// 检查模块是否部分选中
const isModuleIndeterminate = (module: string) => {
  const modulePerms = permissions.value[module]
  const selectedCount = modulePerms.filter(perm => form.value.permissions.includes(perm.value)).length
  return selectedCount > 0 && selectedCount < modulePerms.length
}

// 处理模块全选/取消全选
const handleModuleCheckAll = (module: string, checked: boolean) => {
  const modulePerms = permissions.value[module]
  if (checked) {
    // 添加模块所有权限
    modulePerms.forEach(perm => {
      if (!form.value.permissions.includes(perm.value)) {
        form.value.permissions.push(perm.value)
      }
    })
  } else {
    // 移除模块所有权限
    form.value.permissions = form.value.permissions.filter(
      perm => !modulePerms.some(p => p.value === perm)
    )
  }
}

// 新增角色
const handleAdd = () => {
  dialogType.value = 'add'
  currentId.value = undefined
  form.value = {
    name: '',
    description: '',
    permissions: []
  }
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row: Role) => {
  dialogType.value = 'edit'
  currentId.value = row.id
  form.value = {
    name: row.name,
    description: row.description,
    permissions: [...row.permissions]
  }
  dialogVisible.value = true
}

// 切换角色状态
const handleToggleStatus = async (row: Role) => {
  const action = row.status === 1 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(
      `确定要${action}角色 ${row.name} 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        customClass: 'rounded-lg'
      }
    )
    await toggleRoleStatus(row.id)
    ElMessage.success(`${action}成功`)
    await fetchRoles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action}失败`)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    if (dialogType.value === 'add') {
      await createRole(form.value)
      ElMessage.success('创建成功')
    } else {
      if (!currentId.value) return
      await updateRole(currentId.value, form.value)
      ElMessage.success('更新成功')
    }
    
    dialogVisible.value = false
    await fetchRoles()
  } catch (error) {
    ElMessage.error(dialogType.value === 'add' ? '创建失败' : '更新失败')
  } finally {
    submitting.value = false
  }
}

// 初始化
onMounted(() => {
  fetchRoles()
  fetchPermissions()
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

.role-dialog :deep(.el-overlay-dialog) {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style> 