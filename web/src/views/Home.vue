<template>
  <div class="p-4">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="flex items-center">
              <el-icon class="mr-2"><User /></el-icon>
              <span>用户总数</span>
            </div>
          </template>
          <div class="text-2xl font-bold">1,234</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="flex items-center">
              <el-icon class="mr-2"><Document /></el-icon>
              <span>文章总数</span>
            </div>
          </template>
          <div class="text-2xl font-bold">567</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="flex items-center">
              <el-icon class="mr-2"><View /></el-icon>
              <span>访问量</span>
            </div>
          </template>
          <div class="text-2xl font-bold">8,901</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <template #header>
            <div class="flex items-center">
              <el-icon class="mr-2"><Star /></el-icon>
              <span>收藏数</span>
            </div>
          </template>
          <div class="text-2xl font-bold">234</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { User, Document, View, Star } from '@element-plus/icons-vue'
</script> 