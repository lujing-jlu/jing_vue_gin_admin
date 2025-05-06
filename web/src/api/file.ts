import request from '@/utils/request'

// 文件信息接口
export interface FileInfo {
  id: number
  file_name: string
  file_size: number
  file_type: string
  file_path: string
  category: string
  description: string
  uploaded_by: number
  downloads: number
  created_at: string
  updated_at: string
}

// 文件列表请求接口
export interface FileListRequest {
  file_name?: string
  category?: string
  file_type?: string
  start_date?: string
  end_date?: string
  uploaded_by?: number
  page: number
  page_size: number
}

// 文件列表响应接口
export interface FileListResponse {
  total: number
  list: FileInfo[]
}

// 文件更新请求接口
export interface FileUpdateRequest {
  id: number
  category?: string
  description?: string
  file_name?: string
}

// 文件分类统计接口
export interface FileCategoryCount {
  category: string
  count: number
}

// 文件类型统计接口
export interface FileTypeCount {
  file_type: string
  count: number
}

// 文件统计接口
export interface FileStats {
  total_files: number
  total_size: number
  categories: FileCategoryCount[]
  types: FileTypeCount[]
  recent_uploads: FileInfo[]
  popular_downloads: FileInfo[]
}

export interface FileUploadResponse {
  id: number
  file_name: string
  file_path: string
  file_size: number
  file_type: string
  url: string
}

// 上传文件
export const uploadFile = (file: File, category?: string) => {
  const formData = new FormData()
  formData.append('file', file)
  if (category) {
    formData.append('category', category)
  }
  
  return request<FileUploadResponse>({
    url: '/api/v1/files',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取文件列表
export function getFileList(params: FileListRequest) {
  return request({
    url: '/files',
    method: 'get',
    params
  })
}

// 获取文件详情
export function getFileDetail(id: number) {
  return request({
    url: `/files/${id}`,
    method: 'get'
  })
}

// 更新文件信息
export function updateFile(data: FileUpdateRequest) {
  return request({
    url: '/files',
    method: 'put',
    data
  })
}

// 删除文件
export function deleteFile(id: number) {
  return request({
    url: `/files/${id}`,
    method: 'delete'
  })
}

// 批量删除文件
export function batchDeleteFiles(ids: number[]) {
  return request({
    url: '/files/batch/delete',
    method: 'post',
    data: { ids }
  })
}

// 获取文件统计信息
export function getFileStats() {
  return request({
    url: '/files/stats',
    method: 'get'
  })
}

// 获取文件分类列表
export function getFileCategories() {
  return request({
    url: '/files/categories',
    method: 'get'
  })
}

// 获取文件类型列表
export function getFileTypes() {
  return request({
    url: '/files/types',
    method: 'get'
  })
}

// 获取文件下载链接
export function getFileDownloadUrl(id: number) {
  return `${request.defaults.baseURL}/files/download/${id}`
} 