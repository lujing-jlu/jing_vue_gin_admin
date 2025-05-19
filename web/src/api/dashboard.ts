import request from '@/utils/request'

export interface DashboardStats {
  userCount: number
  roleCount: number
  logCount: number
  fileCount: number
  todayLoginCount: number
  storageUsage: number
  userActivities: {
    date: string
    count: number
  }[]
  fileTypeDistribution: {
    type: string
    count: number
  }[]
  logTypeDistribution: {
    module: string
    count: number
  }[]
  recentRegistrations: {
    date: string
    count: number
  }[]
}

export interface SystemInfo {
  version: string
  environment: string
  serverTime: string
  goVersion: string
  gormVersion: string
  ginVersion: string
  uptime: string
  memoryUsage: string
  cpuUsage: string
}

/**
 * 获取仪表盘统计数据
 */
export const getDashboardStats = () => {
  return request<DashboardStats>({
    url: '/dashboard/stats',
    method: 'get'
  })
}

/**
 * 获取系统信息
 */
export const getSystemInfo = () => {
  return request<SystemInfo>({
    url: '/dashboard/system-info',
    method: 'get'
  })
}

/**
 * 获取用户活跃度数据（最近7天）
 */
export const getUserActivities = () => {
  return request<{ date: string, count: number }[]>({
    url: '/dashboard/user-activities',
    method: 'get'
  })
}

/**
 * 获取文件类型分布
 */
export const getFileTypeDistribution = () => {
  return request<{ type: string, count: number }[]>({
    url: '/dashboard/file-types',
    method: 'get'
  })
}

/**
 * 获取操作日志分布
 */
export const getLogTypeDistribution = () => {
  return request<{ module: string, count: number }[]>({
    url: '/dashboard/log-types',
    method: 'get'
  })
}

/**
 * 获取最近注册用户趋势（最近30天）
 */
export const getRecentRegistrations = () => {
  return request<{ date: string, count: number }[]>({
    url: '/dashboard/recent-registrations',
    method: 'get'
  })
} 