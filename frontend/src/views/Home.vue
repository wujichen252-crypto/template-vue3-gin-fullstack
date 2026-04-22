<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

onMounted(() => {
  if (!authStore.userInfo) {
    authStore.fetchUserInfo()
  }
})
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <div class="bg-background rounded-lg shadow p-6">
      <h1 class="text-2xl font-bold text-primary mb-4">首页</h1>
      <div class="space-y-2">
        <p class="text-primary">欢迎回来，{{ authStore.userInfo?.username || '用户' }}！</p>
        <p class="text-sm text-gray-500">这是您的个人首页</p>
      </div>
    </div>
  </div>
</template>
