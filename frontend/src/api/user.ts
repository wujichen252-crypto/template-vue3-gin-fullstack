import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type { UserInfo, LoginParams, RegisterParams } from '@/types/user'

export const login = (data: LoginParams): Promise<ApiResponse<{ token: string; user_info: UserInfo }>> => {
  return request.post('/v1/auth/login', data)
}

export const register = (data: RegisterParams): Promise<ApiResponse<{ token: string; user_info: UserInfo }>> => {
  return request.post('/v1/auth/register', data)
}

export const getUserInfo = (): Promise<ApiResponse<UserInfo>> => {
  return request.get('/v1/auth/userinfo')
}

export const refreshToken = (): Promise<ApiResponse<{ token: string }>> => {
  return request.post('/v1/auth/refresh')
}

export const logout = (): Promise<ApiResponse<{ message: string }>> => {
  return request.post('/v1/auth/logout')
}
