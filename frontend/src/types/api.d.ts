export interface ApiResponse<T = unknown> {
  code: number
  data: T
  msg: string
  request_id: string
}

export interface PageParams {
  page: number
  page_size: number
}

export interface PageResult<T = unknown> {
  list: T[]
  total: number
  page: number
  page_size: number
}
