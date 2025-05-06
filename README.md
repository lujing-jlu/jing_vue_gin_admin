# Vue Gin Admin

一个基于Vue 3和Gin框架的现代化后台管理系统，提供完整的用户、角色、权限管理和系统日志功能。

## 项目特点

- 🔐 完善的RBAC权限控制系统
- 📝 详细的操作日志记录
- 🖥️ 美观的仪表板界面
- 👤 用户管理与角色分配
- 🛡️ 支持细粒度的权限控制
- 📊 系统数据可视化展示
- 🌐 响应式设计，支持多种设备
- 📂 文件管理与存储功能

## 技术栈

### 前端
- Vue 3.x + TypeScript
- Vue Router
- Element Plus
- Tailwind CSS
- Vite

### 后端
- Gin框架
- GORM
- SQLite/MySQL
- JWT认证

## 功能模块

1. **用户管理**
   - 用户创建、编辑、删除
   - 密码管理
   - 用户状态控制

2. **角色管理**
   - 角色创建、编辑、删除
   - 权限分配
   - 角色状态控制

3. **权限控制**
   - 基于角色的权限控制
   - API级别的权限验证
   - 前端菜单和按钮级权限管理

4. **系统日志**
   - 操作日志记录
   - 登录日志记录
   - 日志查询和导出
   - 日志分类和筛选

5. **系统监控**
   - 系统运行状态概览
   - 用户活跃度统计
   - 操作趋势分析

6. **文件管理**
   - 文件上传与下载
   - 文件分类管理
   - 文件权限控制
   - 文件统计分析

## 安装与运行

### 环境要求
- Node.js v16+
- Go 1.16+
- SQLite或MySQL

### 前端部署
```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 开发模式运行
npm run dev

# 构建生产版本
npm run build
```

### 后端部署
```bash
# 进入后端目录
cd server

# 运行服务
go run cmd/main.go
```

## 配置说明

### 前端配置
前端配置文件位于`web/.env`和`web/.env.production`，主要配置项包括：

- `VITE_API_URL`: API服务器地址
- `VITE_APP_TITLE`: 应用标题

### 后端配置
后端配置文件位于`server/config/config.go`，主要配置项包括：

- 数据库连接信息
- JWT密钥和过期时间
- 服务监听地址和端口
- 文件上传存储路径

## 项目截图

![仪表板](https://example.com/dashboard.png)
![用户管理](https://example.com/user-management.png)
![角色管理](https://example.com/role-management.png)
![系统日志](https://example.com/system-logs.png)
![文件管理](https://example.com/file-management.png)

## 开发计划

- [ ] 添加更多数据统计图表
- [ ] 支持数据导出Excel/PDF
- [ ] 集成消息通知系统
- [ ] 支持第三方登录
- [ ] 多语言国际化支持
- [ ] 深色模式主题
- [x] 文件管理功能
- [ ] 图片预览与处理

## 贡献指南

欢迎提交Issue和Pull Request，请遵循以下步骤：

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交Pull Request

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解更多详情。
