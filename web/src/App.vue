<script setup lang="ts">
</script>

<template>
  <router-view></router-view>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>

<style>
/* 全局样式 */
.el-button--primary {
  --el-button-text-color: #ffffff;
  --el-button-bg-color: #3b82f6;
  --el-button-border-color: #3b82f6;
  --el-button-hover-text-color: #ffffff;
  --el-button-hover-bg-color: #2563eb;
  --el-button-hover-border-color: #2563eb;
}

.el-button--danger {
  --el-button-text-color: #ffffff;
  --el-button-bg-color: #ef4444;
  --el-button-border-color: #ef4444;
  --el-button-hover-text-color: #ffffff;
  --el-button-hover-bg-color: #dc2626;
  --el-button-hover-border-color: #dc2626;
}

.el-button--success {
  --el-button-text-color: #ffffff;
  --el-button-bg-color: #10b981;
  --el-button-border-color: #10b981;
  --el-button-hover-text-color: #ffffff;
  --el-button-hover-bg-color: #059669;
  --el-button-hover-border-color: #059669;
}

.el-button--default {
  --el-button-text-color: #374151;
  --el-button-hover-text-color: #2563eb;
}

.el-dropdown-menu__item {
  --el-dropdown-menuItem-hover-color: #374151;
}
</style>
