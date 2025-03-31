import request from '@/utils/request'

export interface SystemLog {
  id: number
  user_id: number
  username: string
  module: string
  action: string
  resource: string
  detail: string
  ip: string
  user_agent: string
  status: number
  created_at: string
  updated_at: string
}

export interface LogListRequest {
  module?: string
  action?: string
  username?: string
  start_time?: string
  end_time?: string
  page: number
  page_size: number
}

export interface LogListResponse {
  total: number
  list: SystemLog[]
}

/**
 * 获取日志列表
 */
export const getLogList = (params: LogListRequest) => {
  return request<LogListResponse>({
    url: '/logs',
    method: 'get',
    params
  })
}

/**
 * 批量删除日志
 */
export const deleteLogs = (ids: number[]) => {
  return request<any>({
    url: '/logs',
    method: 'delete',
    data: { ids }
  })
}

/**
 * 清空日志
 */
export const clearLogs = () => {
  return request<any>({
    url: '/logs/clear',
    method: 'delete'
  })
}

/**
 * 获取日志模块列表
 */
export const getLogModules = () => {
  return request<string[]>({
    url: '/logs/modules',
    method: 'get'
  })
}

/**
 * 获取日志操作类型列表
 */
export const getLogActions = () => {
  return request<string[]>({
    url: '/logs/actions',
    method: 'get'
  })
} 