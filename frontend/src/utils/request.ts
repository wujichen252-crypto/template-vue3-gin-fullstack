import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import storage from './storage'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

const pendingRequests = new Map<string, AbortController>()

const request: AxiosInstance = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

const generateRequestKey = (config: InternalAxiosRequestConfig): string => {
  return `${config.method}:${config.url}:${JSON.stringify(config.params)}:${JSON.stringify(config.data)}`
}

request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = storage.get('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }

    const requestKey = generateRequestKey(config)
    const controller = new AbortController()
    config.signal = controller.signal
    pendingRequests.set(requestKey, controller)

    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response: AxiosResponse) => {
    const requestKey = generateRequestKey(response.config as InternalAxiosRequestConfig)
    pendingRequests.delete(requestKey)

    const res = response.data
    if (res.code !== 200) {
      ElMessage.error(res.msg || '请求失败')
      return Promise.reject(new Error(res.msg || '请求失败'))
    }
    return res
  },
  (error: AxiosError) => {
    if (error.config) {
      const requestKey = generateRequestKey(error.config as InternalAxiosRequestConfig)
      pendingRequests.delete(requestKey)
    }

    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        ElMessage.error('登录已过期，请重新登录')
        storage.remove('token')
        window.location.href = '/login'
      } else {
        ElMessage.error((error.response.data as { msg?: string })?.msg || '请求失败')
      }
    } else {
      ElMessage.error('网络错误')
    }
    return Promise.reject(error)
  }
)

export const cancelPendingRequest = (config: InternalAxiosRequestConfig): void => {
  const requestKey = generateRequestKey(config)
  const controller = pendingRequests.get(requestKey)
  if (controller) {
    controller.abort()
    pendingRequests.delete(requestKey)
  }
}

export const cancelAllPendingRequests = (): void => {
  pendingRequests.forEach(controller => controller.abort())
  pendingRequests.clear()
}

export default request