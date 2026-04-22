import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface UserState {
  id: number
  username: string
  email: string
}

export const useUserStore = defineStore('user', () => {
  const users = ref<UserState[]>([])
  const currentUser = ref<UserState | null>(null)

  const setCurrentUser = (user: UserState): void => {
    currentUser.value = user
  }

  const clearCurrentUser = (): void => {
    currentUser.value = null
  }

  return {
    users,
    currentUser,
    setCurrentUser,
    clearCurrentUser
  }
})
