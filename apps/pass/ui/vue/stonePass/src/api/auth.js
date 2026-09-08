const TOKEN_KEY = 'stonepass_token'
const AUTH_FLAG_KEY = 'stonepass_auth_enabled'

let authEnabledCache = null

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token) {
  const t = token || ''
  if (t) {
    localStorage.setItem(TOKEN_KEY, t)
    // 供 /files 静态资源与下载链接使用
    document.cookie = `token=${encodeURIComponent(t)}; path=/; SameSite=Lax`
  } else {
    clearToken()
  }
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
  document.cookie = 'token=; path=/; Max-Age=0; SameSite=Lax'
}

export function getCachedAuthEnabled() {
  if (authEnabledCache !== null) return authEnabledCache
  const raw = sessionStorage.getItem(AUTH_FLAG_KEY)
  if (raw === '1') return true
  if (raw === '0') return false
  return null
}

export function setCachedAuthEnabled(enabled) {
  authEnabledCache = !!enabled
  sessionStorage.setItem(AUTH_FLAG_KEY, enabled ? '1' : '0')
}

export async function fetchAuthStatus() {
  const res = await fetch(`/pass/auth-status?_t=${Date.now()}`)
  const data = await res.json().catch(() => null)
  const enabled = !!(data && data.code === 0 && data.data && data.data.enabled)
  setCachedAuthEnabled(enabled)
  return enabled
}

export async function login(password) {
  const res = await fetch('/pass/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ password }),
  })
  const data = await res.json().catch(() => null)
  if (!data || data.code !== 0) {
    throw new Error((data && data.msg) || '登录失败')
  }
  const token = data.data && data.data.token
  if (!token) {
    // 鉴权未启用
    setCachedAuthEnabled(false)
    clearToken()
    return { enabled: false, token: '' }
  }
  setToken(token)
  setCachedAuthEnabled(true)
  return { enabled: true, token }
}

export function logout() {
  clearToken()
}

export function authHeaders() {
  const token = getToken()
  return token ? { token } : {}
}

export function withTokenQuery(url) {
  const token = getToken()
  if (!token) return url
  const join = url.includes('?') ? '&' : '?'
  return `${url}${join}token=${encodeURIComponent(token)}`
}

export function handleUnauthorized() {
  clearToken()
  const hash = location.hash || '#/'
  if (hash.startsWith('#/login')) return
  const redirect = encodeURIComponent(hash.replace(/^#/, '') || '/')
  location.hash = `#/login?redirect=${redirect}`
}
