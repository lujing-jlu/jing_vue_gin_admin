import request from '@/utils/request'

export interface Role {
  id: number
  name: string
  description: string
  permissions: string[]
  status: number
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description: string
  permissions: string[]
}

export interface UpdateRoleRequest {
  name: string
  description: string
  permissions: string[]
}

// 获取角色列表
export const getRoleList = () => {
  return request<Role[]>({
    url: '/roles',
    method: 'get'
  })
}

// 创建角色
export const createRole = (data: CreateRoleRequest) => {
  return request({
    url: '/roles',
    method: 'post',
    data
  })
}

// 更新角色
export const updateRole = (id: number, data: UpdateRoleRequest) => {
  return request({
    url: `/roles/${id}`,
    method: 'put',
    data
  })
}

// 删除角色
export const deleteRole = (id: number) => {
  return request({
    url: `/roles/${id}`,
    method: 'delete'
  })
}

// 切换角色状态
export const toggleRoleStatus = (id: number) => {
  return request({
    url: `/roles/${id}/status`,
    method: 'put'
  })
}

// 获取所有权限
export const getAllPermissions = () => {
  return request<string[]>({
    url: '/permissions',
    method: 'get'
  })
} 