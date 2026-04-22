<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-surface">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-primary mb-4">404</h1>
      <p class="text-xl text-primary mb-6">页面未找到</p>
      <button
        @click="goHome"
        class="px-6 py-2 bg-primary text-background rounded hover:bg-accent transition-colors"
      >
        返回首页
      </button>
    </div>
  </div>
</template>
