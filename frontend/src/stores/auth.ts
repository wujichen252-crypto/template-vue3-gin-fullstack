import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo } from '@/types/user'
import storage, { TOKEN_KEY } from '@/utils/storage'
import { login as apiLogin, register as apiRegister, getUserInfo as apiGetUserInfo, logout as apiLogout } from '@/api/user'
import type { LoginParams, RegisterParams } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(storage.get(TOKEN_KEY) || '')
  const userInfo = ref<UserInfo | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const login = async (params: LoginParams): Promise<void> => {
    const res = await apiLogin(params)
    token.value = res.data.token
    userInfo.value = res.data.user_info
    storage.set(TOKEN_KEY, res.data.token)
  }

  const register = async (params: RegisterParams): Promise<void> => {
    const res = await apiRegister(params)
    token.value = res.data.token
    userInfo.value = res.data.user_info
    storage.set(TOKEN_KEY, res.data.token)
  }

  const fetchUserInfo = async (): Promise<void> => {
    if (!token.value) return
    try {
      const res = await apiGetUserInfo()
      userInfo.value = res.data
    } catch {
      logout()
    }
  }

  const logout = async (): Promise<void> => {
    try {
      await apiLogout()
    } catch {
    }
    token.value = ''
    userInfo.value = null
    storage.remove(TOKEN_KEY)
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    login,
    register,
    fetchUserInfo,
    logout
  }
})
