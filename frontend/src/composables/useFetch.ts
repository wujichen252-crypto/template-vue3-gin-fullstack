import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export const useFetch = async <T>(url: string, options?: RequestInit): Promise<T> => {
  const authStore = useAuthStore()
  const router = useRouter()

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options?.headers
  }

  if (authStore.token) {
    (headers as Record<string, string>)['Authorization'] = `Bearer ${authStore.token}`
  }

  const response = await fetch(url, {
    ...options,
    headers
  })

  if (response.status === 401) {
    authStore.logout()
    router.push('/login')
    throw new Error('Unauthorized')
  }

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return response.json()
}
