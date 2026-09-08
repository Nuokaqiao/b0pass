import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import FilesView from '@/views/FilesView.vue'
import TextView from '@/views/TextView.vue'
import LoginView from '@/views/LoginView.vue'
import {
  fetchAuthStatus,
  getCachedAuthEnabled,
  getToken,
} from '@/api/auth'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', name: 'login', component: LoginView, meta: { public: true } },
    { path: '/', name: 'home', component: HomeView },
    { path: '/files', name: 'files', component: FilesView },
    { path: '/text', name: 'text', component: TextView },
  ],
})

router.beforeEach(async (to) => {
  let enabled = getCachedAuthEnabled()
  if (enabled === null) {
    try {
      enabled = await fetchAuthStatus()
    } catch {
      // 状态接口失败时：有 token 先放行，无 token 进登录
      enabled = true
    }
  }

  if (!enabled) {
    if (to.name === 'login') return { name: 'home' }
    return true
  }

  if (to.meta.public) {
    if (getToken() && to.name === 'login') {
      const redirect = typeof to.query.redirect === 'string' ? to.query.redirect : '/'
      return redirect || '/'
    }
    return true
  }

  if (!getToken()) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }
  return true
})

export default router
