export interface UserInfo {
  id: number
  username: string
  email: string
  avatar_url: string
  status: number
  created_at: string
  updated_at: string
}

export interface LoginParams {
  username: string
  password: string
}

export interface RegisterParams {
  username: string
  password: string
  email: string
}
