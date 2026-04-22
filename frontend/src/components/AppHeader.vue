<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="bg-primary text-background border-b border-surface">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <h1 class="text-xl font-bold">全栈模板</h1>
      <nav class="flex items-center gap-4">
        <template v-if="authStore.isAuthenticated">
          <span class="text-sm">{{ authStore.userInfo?.username }}</span>
          <button
            @click="handleLogout"
            class="px-4 py-2 text-sm bg-surface text-primary rounded hover:bg-accent hover:text-background transition-colors"
          >
            退出
          </button>
        </template>
        <template v-else>
          <router-link to="/login" class="text-sm hover:text-accent transition-colors">登录</router-link>
        </template>
      </nav>
    </div>
  </header>
</template>
