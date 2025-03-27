<template>
  <div class="p-4">
    <div class="mb-4 flex justify-between items-center">
      <h2 class="text-2xl font-bold">角色管理</h2>
      <el-button type="primary" @click="handleAdd">新增角色</el-button>
    </div>

    <!-- 角色列表 -->
    <el-table :data="roles" border stripe>
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="角色名称" />
      <el-table-column prop="description" label="描述" />
      <el-table-column label="权限" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="perm in row.permissions" :key="perm" class="mr-1 mb-1">
            {{ permissionLabels[perm] || perm }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="handleStatusChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 角色表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增角色' : '编辑角色'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
        class="mt-4"
      >
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请输入角色描述"
          />
        </el-form-item>
        <el-form-item label="权限" prop="permissions">
          <el-checkbox-group v-model="form.permissions">
            <div v-for="(perms, module) in groupedPermissions" :key="module" class="mb-4">
              <div class="font-bold mb-2">{{ module }}：</div>
              <el-checkbox
                v-for="perm in perms"
                :key="perm.value"
                :label="perm.value"
              >
                {{ perm.label }}
              </el-checkbox>
            </div>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getRoleList,
  createRole,
  updateRole,
  deleteRole,
  toggleRoleStatus,
  getAllPermissions,
  type Role,
  type CreateRoleRequest,
  type UpdateRoleRequest
} from '@/api/role'

// 角色列表数据
const roles = ref<Role[]>([])

// 权限标签映射
const permissionLabels: Record<string, string> = {
  'user:view': '查看用户',
  'user:create': '创建用户',
  'user:edit': '编辑用户',
  'user:delete': '删除用户',
  'role:view': '查看角色',
  'role:create': '创建角色',
  'role:edit': '编辑角色',
  'role:delete': '删除角色'
}

// 分组后的权限列表
const groupedPermissions = {
  用户管理: [
    { label: '查看用户', value: 'user:view' },
    { label: '创建用户', value: 'user:create' },
    { label: '编辑用户', value: 'user:edit' },
    { label: '删除用户', value: 'user:delete' }
  ],
  角色管理: [
    { label: '查看角色', value: 'role:view' },
    { label: '创建角色', value: 'role:create' },
    { label: '编辑角色', value: 'role:edit' },
    { label: '删除角色', value: 'role:delete' }
  ]
}

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
const rules: FormRules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入角色描述', trigger: 'blur' }],
  permissions: [{ required: true, message: '请选择权限', trigger: 'change' }]
}

// 获取角色列表
const fetchRoles = async () => {
  try {
    roles.value = await getRoleList()
  } catch (error) {
    ElMessage.error('获取角色列表失败')
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
    permissions: row.permissions
  }
  dialogVisible.value = true
}

// 删除角色
const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm('确定要删除该角色吗？', '提示', {
      type: 'warning'
    })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    await fetchRoles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 切换角色状态
const handleStatusChange = async (row: Role) => {
  try {
    await toggleRoleStatus(row.id)
    ElMessage.success('状态更新成功')
    await fetchRoles()
  } catch (error) {
    ElMessage.error('状态更新失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
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
  }
}

// 初始化
onMounted(() => {
  fetchRoles()
})
</script> 