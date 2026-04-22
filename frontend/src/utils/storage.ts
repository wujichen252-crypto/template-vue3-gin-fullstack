const TOKEN_KEY = 'token'

const storage = {
  get: (key: string): string | null => {
    return localStorage.getItem(key)
  },
  set: (key: string, value: string): void => {
    localStorage.setItem(key, value)
  },
  remove: (key: string): void => {
    localStorage.removeItem(key)
  },
  clear: (): void => {
    localStorage.clear()
  }
}

export { TOKEN_KEY }
export default storage
