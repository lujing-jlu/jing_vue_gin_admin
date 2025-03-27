import request from '@/utils/request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
  user: {
    id: number
    username: string
    nickname: string
    avatar: string
    role: string
    status: number
  }
}

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
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

export const login = (data: LoginRequest) => {
  return request<LoginResponse>({
    url: '/api/login',
    method: 'post',
    data
  })
}

export const getUserList = () => {
  return request<User[]>({
    url: '/api/users',
    method: 'get'
  })
}

export const createUser = (data: CreateUserRequest) => {
  return request<User>({
    url: '/api/users',
    method: 'post',
    data
  })
}

export const updateUser = (id: number, data: UpdateUserRequest) => {
  return request<User>({
    url: `/api/users/${id}`,
    method: 'put',
    data
  })
}

export const toggleUserStatus = (id: number) => {
  return request<User>({
    url: `/api/users/${id}/status`,
    method: 'put'
  })
} 