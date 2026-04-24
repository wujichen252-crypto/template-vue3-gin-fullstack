# Vue3 前端开发规范

## 代码风格

### 1. 组件规范
- 使用 `<script setup>` + Composition API
- 组件文件使用 PascalCase 命名
- 组件 props 必须定义类型
- 组件必须使用 defineProps 和 defineEmits
- 组件使用异步导入（defineAsyncComponent）实现懒加载

### 2. 样式规范
- 使用 Tailwind CSS
- 颜色使用预设的配色方案（主色 #000000、#FFFFFF、#F5F5F5，强调色 #2563EB）
- 避免使用内联样式，优先使用 Tailwind 类名

### 3. TypeScript 规范
- 禁止使用 `any` 类型
- 必须为所有函数参数和返回值定义类型
- 接口命名使用 PascalCase
- 枚举命名使用 PascalCase

### 4. 命名规范
- 组件: PascalCase (如: UserProfile.vue)
- 变量/函数: camelCase (如: userName, getUserInfo)
- 常量: UPPER_SNAKE_CASE (如: API_BASE_URL)
- CSS 类: 小写加连字符 (如: btn-primary)
- 类型/接口: PascalCase (如: UserInfo, ApiResponse)

## 项目结构

```
src/
├── api/              # API 接口封装
├── assets/          # 静态资源（图片、字体等）
├── components/       # 全局组件
├── composables/     # 组合式函数（可复用的响应式逻辑）
├── router/          # 路由配置
├── stores/          # Pinia 状态管理
├── test/            # 测试文件
│   ├── setup.ts    # 测试环境配置
│   ├── api.test.ts  # API 接口测试
│   └── stores/      # Store 单元测试
├── types/          # 类型定义
├── utils/          # 工具函数
├── views/          # 页面组件
├── App.vue         # 根组件
└── main.ts         # 应用入口
```

## API 设计

- 使用 axios 封装请求
- 统一处理响应拦截（code != 200 视为错误）
- 401 响应自动跳转登录页
- 所有 API 调用通过 api/ 目录下的模块
- 请求/响应类型使用 TypeScript 接口定义

## 性能优化

- 图片使用懒加载（v-lazy 或自定义指令）
- 路由使用懒加载（Vue Router 的 component: () => import(...)）
- 组件使用异步导入（defineAsyncComponent）
- 合理使用 computed 和 watch
- Three.js 场景必须封装为 composable，注意内存释放（dispose）

## 测试规范

### 测试框架
- Vitest（单元测试框架）
- @vue/test-utils（Vue 组件测试）
- Happy DOM（DOM 环境）

### 测试文件命名
- 组件测试: `*.test.ts` 或 `*.spec.ts`
- 放在同目录或 `test/` 目录下

### 测试命令
```bash
# 运行所有测试
npm run test

# 监听模式（开发时）
npm run test:watch

# 类型检查
npm run typecheck

# ESLint 检查
npm run lint
```

### 测试覆盖率要求
- 核心业务逻辑（stores, composables）覆盖率 > 70%
- 工具函数覆盖率 > 80%
