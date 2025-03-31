import request from '@/utils/request'

export interface User {
  id: number
  username: string
  nickname: string
  role: string
  status: number
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  username: string
  password: string
  nickname: string
  role: string
}

export interface UpdateUserRequest {
  nickname: string
  role: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// 用户登录
export const login = (data: LoginRequest) => {
  return request<LoginResponse>({
    url: '/login',
    method: 'post',
    data
  })
}

// 获取用户列表
export const getUserList = () => {
  return request<User[]>({
    url: '/users',
    method: 'get'
  })
}

// 创建用户
export const createUser = (data: CreateUserRequest) => {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

// 更新用户
export const updateUser = (id: number, data: UpdateUserRequest) => {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export const deleteUser = (id: number) => {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

// 切换用户状态
export const toggleUserStatus = (id: number) => {
  return request({
    url: `/users/${id}/status`,
    method: 'put'
  })
}

/**
 * 更新个人信息
 */
export const updateProfile = (data: { nickname: string, avatar?: string }) => {
  return request<any>({
    url: '/user/profile',
    method: 'put',
    data
  })
}

/**
 * 修改密码
 */
export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export const changePassword = (data: ChangePasswordRequest) => {
  return request<any>({
    url: '/user/password',
    method: 'put',
    data
  })
} 