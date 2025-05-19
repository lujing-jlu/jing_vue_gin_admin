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

    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden mb-4">
      <div class="mb-4 flex flex-wrap gap-4">
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索角色名称/描述" 
          class="max-w-xs" 
          clearable
          prefix-icon="Search"
        />
        <el-select v-model="filterStatus" placeholder="状态" class="w-32" clearable>
          <el-option label="正常" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button type="primary" plain class="text-blue-600" @click="handleSearch">
          <el-icon class="mr-1"><Search /></el-icon>
          搜索
        </el-button>
        <el-button class="text-gray-600" @click="handleReset">
          <el-icon class="mr-1"><Refresh /></el-icon>
          重置
        </el-button>
      </div>
    </el-card>

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

      <div class="flex justify-end mt-4">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          background
          class="!p-0"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
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
import { Plus, Edit, Lock, Unlock, Search, Refresh } from '@element-plus/icons-vue'
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
  type PermissionGroup,
  type PermissionInfo,
  type RoleListRequest
} from '@/api/role'

// 角色列表数据
const roles = ref<Role[]>([])
const loading = ref(false)
const submitting = ref(false)
const permissions = ref<PermissionGroup>({})
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

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
    const params: RoleListRequest = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchQuery.value || undefined,
      status: filterStatus.value !== '' ? Number(filterStatus.value) : undefined
    }
    
    const res = await getRoleList(params)
    roles.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取角色列表失败:', error)
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
    console.error('获取权限列表失败:', error)
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

// 模块全选/取消全选
const handleModuleCheckAll = (module: string, val: boolean) => {
  const modulePerms = permissions.value[module]
  if (val) {
    // 全选
    modulePerms.forEach(perm => {
      if (!form.value.permissions.includes(perm.value)) {
        form.value.permissions.push(perm.value)
      }
    })
  } else {
    // 取消全选
    form.value.permissions = form.value.permissions.filter(
      value => !modulePerms.some(perm => perm.value === value)
    )
  }
}

// 添加角色
const handleAdd = () => {
  dialogType.value = 'add'
  form.value = {
    name: '',
    description: '',
    permissions: []
  }
  currentId.value = undefined
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row: Role) => {
  dialogType.value = 'edit'
  form.value = {
    name: row.name,
    description: row.description,
    permissions: [...row.permissions]
  }
  currentId.value = row.id
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
        type: 'warning'
      }
    )
    
    await toggleRoleStatus(row.id)
    ElMessage.success(`${action}成功`)
    fetchRoles()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        submitting.value = true
        
        if (dialogType.value === 'add') {
          await createRole(form.value)
          ElMessage.success('添加角色成功')
        } else if (currentId.value) {
          await updateRole(currentId.value, form.value)
          ElMessage.success('更新角色成功')
        }
        
        dialogVisible.value = false
        fetchRoles()
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 搜索角色
const handleSearch = () => {
  currentPage.value = 1
  fetchRoles()
}

// 重置搜索条件
const handleReset = () => {
  searchQuery.value = ''
  filterStatus.value = ''
  currentPage.value = 1
  fetchRoles()
}

// 页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchRoles()
}

// 每页条数变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchRoles()
}

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
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
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