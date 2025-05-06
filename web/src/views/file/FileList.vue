<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">文件管理</h2>
        <p class="text-gray-600">上传和管理系统文件</p>
      </div>
      <div class="space-x-2 flex">
        <el-button type="primary" @click="showUploadDialog = true" style="color: white !important; background-color: #3b82f6 !important;">
          <el-icon class="mr-2"><Upload /></el-icon>
          上传文件
        </el-button>
        <el-button 
          type="danger" 
          :disabled="selectedFiles.length === 0"
          @click="handleBatchDelete"
          style="color: white !important; background-color: #ef4444 !important;"
        >
          <el-icon class="mr-2"><Delete /></el-icon>
          批量删除
        </el-button>
      </div>
    </div>

    <!-- 文件统计 -->
    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-blue-500 to-blue-600">
            <el-icon class="text-2xl text-white"><Document /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">文件总数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ stats.total_files || 0 }}</p>
          </div>
        </div>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-green-500 to-green-600">
            <el-icon class="text-2xl text-white"><Files /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">总存储空间</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ formatFileSize(stats.total_size) }}</p>
          </div>
        </div>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-purple-500 to-purple-600">
            <el-icon class="text-2xl text-white"><TopRight /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">今日上传</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ todayUploads }}</p>
          </div>
        </div>
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-gradient-to-br from-orange-500 to-orange-600">
            <el-icon class="text-2xl text-white"><Download /></el-icon>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">总下载次数</p>
            <p class="text-2xl font-bold text-gray-900 mt-1">{{ totalDownloads }}</p>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 文件搜索 -->
    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden mb-6">
      <div class="mb-4 flex flex-wrap gap-4">
        <el-input 
          v-model="searchParams.file_name" 
          placeholder="搜索文件名" 
          class="max-w-xs" 
          clearable
          prefix-icon="Search"
        />
        <el-select v-model="searchParams.category" placeholder="文件分类" clearable class="w-32">
          <el-option 
            v-for="category in categories" 
            :key="category" 
            :label="formatCategory(category)" 
            :value="category" 
          />
        </el-select>
        <el-select v-model="searchParams.file_type" placeholder="文件类型" clearable class="w-32">
          <el-option 
            v-for="type in fileTypes" 
            :key="type" 
            :label="type.toUpperCase()" 
            :value="type" 
          />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
          class="max-w-md"
        />
        <el-button type="primary" @click="handleSearch" class="text-white" style="color: white !important; background-color: #3b82f6 !important;">
          <el-icon class="mr-1"><Search /></el-icon>
          搜索
        </el-button>
        <el-button @click="handleReset" class="text-gray-600">
          <el-icon class="mr-1"><Refresh /></el-icon>
          重置
        </el-button>
      </div>
    </el-card>

    <!-- 文件列表 -->
    <el-card class="border-0 shadow-sm rounded-xl overflow-hidden">
      <el-table 
        :data="files" 
        border 
        stripe 
        v-loading="loading"
        class="rounded-lg overflow-hidden"
        :header-cell-style="{ background: '#f9fafb', color: '#374151', fontWeight: '600' }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="file_name" label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="flex items-center">
              <span class="mr-2">{{ getFileIcon(row.file_type) }}</span>
              <el-tooltip :content="row.file_name" placement="top">
                <span class="truncate">{{ row.file_name }}</span>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="file_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getFileTypeColor(row.file_type)" class="rounded-full px-2 py-0" size="small">
              {{ row.file_type.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag type="info" class="rounded-full px-2 py-0" size="small" effect="plain">
              {{ formatCategory(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="file_size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="downloads" label="下载次数" width="100" align="center" />
        <el-table-column prop="created_at" label="上传时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="handleDownload(row)"
            >
              <el-icon><Download /></el-icon>
              下载
            </el-button>
            <el-button 
              type="primary" 
              link 
              size="small" 
              @click="handleEdit(row)"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button 
              type="danger" 
              link 
              size="small" 
              @click="handleDelete(row)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
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

    <!-- 上传文件对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传文件"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="uploadForm" label-position="top">
        <el-form-item label="选择文件" required>
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">拖拽文件到此处或 <em>点击上传</em></div>
            <template #tip>
              <div class="el-upload__tip">
                请上传任意类型文件
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="文件分类">
          <el-select v-model="uploadForm.category" placeholder="选择分类" class="w-full">
            <el-option label="文档" value="document" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="音频" value="audio" />
            <el-option label="压缩包" value="archive" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件描述">
          <el-input v-model="uploadForm.description" type="textarea" rows="3" placeholder="请输入文件描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleUpload" 
            :loading="uploading"
            :disabled="!uploadForm.file"
            style="color: white !important; background-color: #3b82f6 !important;"
          >
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑文件对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑文件信息"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="editForm" label-position="top">
        <el-form-item label="文件名">
          <el-input v-model="editForm.file_name" placeholder="请输入文件名" />
        </el-form-item>
        <el-form-item label="文件分类">
          <el-select v-model="editForm.category" placeholder="选择分类" class="w-full">
            <el-option label="文档" value="document" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="音频" value="audio" />
            <el-option label="压缩包" value="archive" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件描述">
          <el-input v-model="editForm.description" type="textarea" rows="3" placeholder="请输入文件描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleUpdate" 
            :loading="updating"
            style="color: white !important; background-color: #3b82f6 !important;"
          >
            更新
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { 
  Search, Refresh, Delete, Upload, Download, Document, Edit, 
  Files, TopRight, UploadFilled, DocumentCopy, Picture
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  getFileList, uploadFile, deleteFile, updateFile, batchDeleteFiles,
  getFileStats, getFileCategories, getFileTypes, getFileDownloadUrl,
  type FileInfo, type FileListRequest, type FileStats
} from '@/api/file'

// 文件列表相关
const files = ref<FileInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 文件统计信息
const stats = ref<FileStats>({
  total_files: 0,
  total_size: 0,
  categories: [],
  types: [],
  recent_uploads: [],
  popular_downloads: []
})

// 文件分类和类型选择
const categories = ref<string[]>([])
const fileTypes = ref<string[]>([])

// 选中的文件
const selectedFiles = ref<FileInfo[]>([])

// 上传文件
const showUploadDialog = ref(false)
const uploading = ref(false)
const uploadForm = reactive({
  file: null as File | null,
  category: 'other',
  description: ''
})

// 编辑文件
const showEditDialog = ref(false)
const updating = ref(false)
const editForm = reactive({
  id: 0,
  file_name: '',
  category: '',
  description: ''
})

// 搜索参数
const searchParams = reactive<FileListRequest>({
  file_name: '',
  category: '',
  file_type: '',
  start_date: '',
  end_date: '',
  page: 1,
  page_size: 10
})

// 日期范围
const dateRange = ref<[string, string] | null>(null)

// 计算属性
const todayUploads = computed(() => {
  const today = new Date().toISOString().split('T')[0]
  return stats.value.recent_uploads?.filter(file => 
    file.created_at.startsWith(today)
  ).length || 0
})

const totalDownloads = computed(() => {
  let total = 0
  files.value.forEach(file => {
    total += file.downloads
  })
  return total
})

// 获取文件列表
const fetchFiles = async () => {
  try {
    loading.value = true
    const res = await getFileList({
      ...searchParams,
      page: currentPage.value,
      page_size: pageSize.value
    })
    files.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch files:', error)
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

// 获取文件统计信息
const fetchFileStats = async () => {
  try {
    const res = await getFileStats()
    stats.value = res.data
  } catch (error) {
    console.error('Failed to fetch file stats:', error)
  }
}

// 获取文件分类列表
const fetchCategories = async () => {
  try {
    const res = await getFileCategories()
    categories.value = res.data
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

// 获取文件类型列表
const fetchFileTypes = async () => {
  try {
    const res = await getFileTypes()
    fileTypes.value = res.data
  } catch (error) {
    console.error('Failed to fetch file types:', error)
  }
}

// 处理分页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchFiles()
}

// 处理当前页变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchFiles()
}

// 处理搜索
const handleSearch = () => {
  if (dateRange.value) {
    searchParams.start_date = dateRange.value[0]
    searchParams.end_date = dateRange.value[1]
  } else {
    searchParams.start_date = ''
    searchParams.end_date = ''
  }
  
  currentPage.value = 1
  fetchFiles()
}

// 处理重置
const handleReset = () => {
  searchParams.file_name = ''
  searchParams.category = ''
  searchParams.file_type = ''
  searchParams.start_date = ''
  searchParams.end_date = ''
  dateRange.value = null
  currentPage.value = 1
  fetchFiles()
}

// 处理选择变化
const handleSelectionChange = (selection: FileInfo[]) => {
  selectedFiles.value = selection
}

// 处理文件变更
const handleFileChange = (file: any) => {
  uploadForm.file = file.raw
}

// 处理上传
const handleUpload = async () => {
  if (!uploadForm.file) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  try {
    uploading.value = true
    
    const formData = new FormData()
    formData.append('file', uploadForm.file)
    formData.append('category', uploadForm.category)
    formData.append('description', uploadForm.description)
    
    await uploadFile(formData)
    ElMessage.success('文件上传成功')
    showUploadDialog.value = false
    
    // 重置表单
    uploadForm.file = null
    uploadForm.category = 'other'
    uploadForm.description = ''
    
    // 刷新数据
    fetchFiles()
    fetchFileStats()
    fetchCategories()
    fetchFileTypes()
  } catch (error) {
    console.error('Failed to upload file:', error)
    ElMessage.error('文件上传失败')
  } finally {
    uploading.value = false
  }
}

// 处理下载
const handleDownload = (file: FileInfo) => {
  const downloadUrl = getFileDownloadUrl(file.id)
  window.open(downloadUrl, '_blank')
}

// 处理编辑
const handleEdit = (file: FileInfo) => {
  editForm.id = file.id
  editForm.file_name = file.file_name
  editForm.category = file.category
  editForm.description = file.description
  showEditDialog.value = true
}

// 处理更新
const handleUpdate = async () => {
  try {
    updating.value = true
    await updateFile(editForm)
    ElMessage.success('文件信息更新成功')
    showEditDialog.value = false
    fetchFiles()
  } catch (error) {
    console.error('Failed to update file:', error)
    ElMessage.error('文件信息更新失败')
  } finally {
    updating.value = false
  }
}

// 处理删除
const handleDelete = async (file: FileInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件 "${file.file_name}" 吗？`, 
      '提示', 
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteFile(file.id)
    ElMessage.success('删除成功')
    fetchFiles()
    fetchFileStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete file:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 处理批量删除
const handleBatchDelete = async () => {
  if (selectedFiles.value.length === 0) {
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedFiles.value.length} 个文件吗？`, 
      '提示', 
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedFiles.value.map(file => file.id)
    await batchDeleteFiles(ids)
    ElMessage.success('批量删除成功')
    fetchFiles()
    fetchFileStats()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to batch delete files:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + ' KB'
  } else if (size < 1024 * 1024 * 1024) {
    return (size / (1024 * 1024)).toFixed(2) + ' MB'
  } else {
    return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
  }
}

// 格式化分类名称
const formatCategory = (category: string) => {
  const categoryMap: Record<string, string> = {
    'document': '文档',
    'image': '图片',
    'video': '视频',
    'audio': '音频',
    'archive': '压缩包',
    'other': '其他'
  }
  return categoryMap[category] || category
}

// 获取文件图标
const getFileIcon = (fileType: string) => {
  const iconMap: Record<string, any> = {
    'pdf': '📄',
    'doc': '📝',
    'docx': '📝',
    'xls': '📊',
    'xlsx': '📊',
    'ppt': '📽️',
    'pptx': '📽️',
    'txt': '📃',
    'jpg': '🖼️',
    'jpeg': '🖼️',
    'png': '🖼️',
    'gif': '🖼️',
    'mp4': '🎬',
    'avi': '🎬',
    'mp3': '🎵',
    'wav': '🎵',
    'zip': '📦',
    'rar': '📦',
    '7z': '📦'
  }
  return iconMap[fileType] || '📁'
}

// 获取文件类型颜色
const getFileTypeColor = (fileType: string) => {
  const typeColorMap: Record<string, string> = {
    'pdf': 'danger',
    'doc': 'primary',
    'docx': 'primary',
    'xls': 'success',
    'xlsx': 'success',
    'ppt': 'warning',
    'pptx': 'warning',
    'txt': 'info',
    'jpg': 'success',
    'jpeg': 'success',
    'png': 'success',
    'gif': 'success',
    'mp4': 'danger',
    'avi': 'danger',
    'mp3': 'warning',
    'wav': 'warning',
    'zip': 'info',
    'rar': 'info',
    '7z': 'info'
  }
  return typeColorMap[fileType] || 'info'
}

onMounted(() => {
  fetchFiles()
  fetchFileStats()
  fetchCategories()
  fetchFileTypes()
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

:deep(.el-select) {
  width: 100%;
}

:deep(.el-button) {
  border-radius: 0.5rem;
  transition: all 0.2s;
}

:deep(.el-button:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

:deep(.el-tag) {
  border-radius: 9999px;
}

:deep(.el-tag--primary.is-plain) {
  background-color: #e0f2fe;
  border-color: #e0f2fe;
  color: #0ea5e9;
}

:deep(.el-tag--success.is-plain) {
  background-color: #dcfce7;
  border-color: #dcfce7;
  color: #16a34a;
}

:deep(.el-tag--warning.is-plain) {
  background-color: #fef3c7;
  border-color: #fef3c7;
  color: #d97706;
}

:deep(.el-tag--danger.is-plain) {
  background-color: #fee2e2;
  border-color: #fee2e2;
  color: #dc2626;
}

:deep(.el-tag--info.is-plain) {
  background-color: #f3f4f6;
  border-color: #f3f4f6;
  color: #4b5563;
}

:deep(.el-upload) {
  width: 100%;
}

:deep(.el-upload-dragger) {
  width: 100%;
  height: 180px;
}

:deep(.el-table .cell) {
  white-space: nowrap;
}

:deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #3b82f6;
}
</style> 