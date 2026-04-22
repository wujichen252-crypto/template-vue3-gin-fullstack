import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

export const useAuth = () => {
  const authStore = useAuthStore()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const userInfo = computed(() => authStore.userInfo)

  return {
    isAuthenticated,
    userInfo,
    login: authStore.login,
    register: authStore.register,
    logout: authStore.logout,
    fetchUserInfo: authStore.fetchUserInfo
  }
}
