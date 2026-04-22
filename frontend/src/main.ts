import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { ElButton, ElMessage, ElForm, ElFormItem, ElInput } from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './assets/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElButton)
app.use(ElMessage)
app.use(ElForm)
app.use(ElFormItem)
app.use(ElInput)

app.mount('#app')
