import { describe, it, expect } from 'vitest'
import { login as apiLogin } from '@/api/user'

describe('User API', () => {
  it('should have login function', () => {
    expect(apiLogin).toBeDefined()
    expect(typeof apiLogin).toBe('function')
  })
})