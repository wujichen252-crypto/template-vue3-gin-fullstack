# Vue3 前端开发规范

## 代码风格

### 1. 组件规范
- 使用 `<script setup>` + Composition API
- 组件文件使用 PascalCase 命名
- 组件 props 必须定义类型
- 组件必须使用 defineProps 和 defineEmits

### 2. 样式规范
- 使用 Tailwind CSS
- 颜色使用预设的配色方案
- 避免使用内联样式，优先使用 Tailwind 类名

### 3. TypeScript 规范
- 禁止使用 `any` 类型
- 必须为所有函数参数和返回值定义类型
- 接口命名使用 PascalCase

### 4. 命名规范
- 组件: PascalCase (如: UserProfile.vue)
- 变量/函数: camelCase (如: userName, getUserInfo)
- 常量: UPPER_SNAKE_CASE (如: API_BASE_URL)
- CSS 类: 小写加连字符 (如: btn-primary)

## 项目结构

```
src/
├── api/              # API 接口封装
├── components/       # 全局组件
├── composables/     # 组合式函数
├── router/          # 路由配置
├── stores/          # Pinia 状态管理
├── types/          # 类型定义
├── utils/          # 工具函数
└── views/          # 页面组件
```

## API 设计

- 使用 axios 封装请求
- 统一处理响应拦截
- 401 响应自动跳转登录页
- 所有 API 调用通过 api/ 目录下的模块

## 性能优化

- 图片使用懒加载
- 路由使用懒加载
- 组件使用异步导入
- 合理使用 computed 和 watch
