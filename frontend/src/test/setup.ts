import { beforeEach, vi } from 'vitest'

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn(),
    warning: vi.fn(),
    info: vi.fn()
  }
}))

beforeEach(() => {
  vi.clearAllMocks()
})

global.localStorage = {
  getItem: vi.fn(() => null),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
} as unknown as Storage

global.window = {
  location: {
    href: '',
    pathname: '/'
  },
  navigator: {
    userAgent: 'test'
  }
} as unknown as Window & typeof globalThis
