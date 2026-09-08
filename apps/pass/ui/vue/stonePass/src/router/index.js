import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import FilesView from '@/views/FilesView.vue'
import TextView from '@/views/TextView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/files', name: 'files', component: FilesView },
    { path: '/text', name: 'text', component: TextView },
  ],
})

export default router
