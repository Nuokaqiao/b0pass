<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  addNode,
  deleteNode,
  downloadUrl,
  fetchFileList,
  staticFileUrl,
  uploadFile,
} from '@/api/pass'

const currentPath = ref('/')
const items = ref([])
const loading = ref(false)
const error = ref('')
const uploading = ref(false)
const progress = ref(0)
const fileInput = ref(null)
const toast = ref('')

const crumbs = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  const list = [{ name: '根目录', path: '/' }]
  let acc = ''
  for (const part of parts) {
    acc += `/${part}`
    list.push({ name: part, path: acc })
  }
  return list
})

function normalizeDir(path) {
  if (!path || path === '/') return '/'
  return path.endsWith('/') ? path : `${path}/`
}

function showToast(msg) {
  toast.value = msg
  setTimeout(() => {
    if (toast.value === msg) toast.value = ''
  }, 2200)
}

async function loadList(path = currentPath.value) {
  loading.value = true
  error.value = ''
  try {
    const res = await fetchFileList(path || '/')
    items.value = Array.isArray(res.data) ? res.data : []
    currentPath.value = path || '/'
  } catch (e) {
    error.value = e.message || '加载失败'
    items.value = []
  } finally {
    loading.value = false
  }
}

function openItem(item) {
  if (item.type === 'dir') {
    loadList(item.path.endsWith('/') ? item.path : `${item.path}/`)
    return
  }
  window.open(downloadUrl(item.path), '_blank')
}

async function removeItem(item) {
  if (!confirm(`确定删除「${item.name}」？`)) return
  try {
    await deleteNode(item.path)
    showToast('已删除')
    await loadList()
  } catch (e) {
    showToast(e.message || '删除失败')
  }
}

async function createFolder() {
  const name = prompt('新建文件夹名称')
  if (!name || !name.trim()) return
  const path = `${normalizeDir(currentPath.value)}${name.trim()}/`
  try {
    await addNode(path)
    showToast('已创建')
    await loadList()
  } catch (e) {
    showToast(e.message || '创建失败')
  }
}

function triggerUpload() {
  fileInput.value?.click()
}

async function onFilesSelected(e) {
  const files = Array.from(e.target.files || [])
  e.target.value = ''
  if (!files.length) return
  uploading.value = true
  progress.value = 0
  const dir = normalizeDir(currentPath.value)
  try {
    for (const file of files) {
      await uploadFile(dir, file, (p) => {
        progress.value = p
      })
    }
    showToast('上传完成')
    await loadList()
  } catch (err) {
    showToast(err.message || '上传失败')
  } finally {
    uploading.value = false
    progress.value = 0
  }
}

function previewUrl(item) {
  if (item.type === 'img') return staticFileUrl(item.path)
  return ''
}

onMounted(() => loadList('/'))
</script>

<template>
  <div class="page">
    <header class="topbar shell">
      <RouterLink class="back-link" :to="{ name: 'home' }">← 首页</RouterLink>
      <div class="brand">传文件</div>
      <span class="muted tip">高速传输</span>
    </header>

    <main class="shell">
      <section class="panel files-panel">
        <div class="toolbar">
          <nav class="crumbs">
            <button
              v-for="(c, i) in crumbs"
              :key="c.path"
              type="button"
              class="crumb"
              @click="loadList(c.path === '/' ? '/' : c.path + '/')"
            >
              <span v-if="i">/</span>{{ c.name }}
            </button>
          </nav>
          <div class="actions">
            <button class="btn btn-ghost" type="button" @click="createFolder">新建</button>
            <button class="btn btn-ghost" type="button" @click="loadList()">刷新</button>
            <button class="btn btn-primary" type="button" :disabled="uploading" @click="triggerUpload">
              {{ uploading ? `上传中 ${progress}%` : '上传' }}
            </button>
            <input ref="fileInput" type="file" multiple hidden @change="onFilesSelected" />
          </div>
        </div>

        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="toast" class="toast">{{ toast }}</p>

        <div v-if="loading" class="empty">加载中…</div>
        <div v-else-if="!items.length" class="empty">当前目录为空，点击「上传」添加文件</div>
        <ul v-else class="file-list">
          <li v-for="item in items" :key="item.path" class="file-row">
            <button class="file-main" type="button" @click="openItem(item)">
              <img v-if="previewUrl(item)" class="thumb" :src="previewUrl(item)" alt="" />
              <span v-else class="icon" :class="item.type === 'dir' ? 'is-dir' : 'is-file'"></span>
              <span class="meta">
                <strong>{{ item.name }}</strong>
                <small class="muted">{{ item.type === 'dir' ? '文件夹' : item.sizes }} · {{ item.date }}</small>
              </span>
            </button>
            <div class="row-actions">
              <a
                v-if="item.type !== 'dir'"
                class="btn btn-ghost btn-sm"
                :href="downloadUrl(item.path)"
                target="_blank"
                rel="noopener"
              >下载</a>
              <button class="btn btn-danger btn-sm" type="button" @click="removeItem(item)">删除</button>
            </div>
          </li>
        </ul>
      </section>
    </main>
  </div>
</template>

<style scoped>
.tip {
  font-size: 0.85rem;
}

.files-panel {
  padding: 18px;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.crumbs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  align-items: center;
}

.crumb {
  border: none;
  background: transparent;
  color: var(--brand);
  cursor: pointer;
  padding: 4px 2px;
  font-weight: 500;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.btn-sm {
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.88rem;
}

.error,
.toast {
  margin: 0 0 10px;
  padding: 10px 12px;
  border-radius: 10px;
}

.error {
  background: #fdecec;
  color: var(--danger);
}

.toast {
  background: var(--brand-soft);
  color: var(--brand-deep);
}

.file-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.file-row {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: space-between;
  padding: 12px 8px;
  border-top: 1px solid var(--line);
}

.file-row:first-child {
  border-top: none;
}

.file-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
  border: none;
  background: transparent;
  text-align: left;
  cursor: pointer;
  padding: 0;
}

.icon,
.thumb {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  flex-shrink: 0;
  object-fit: cover;
  background: #edf2f7;
  display: grid;
  place-items: center;
  position: relative;
}

.icon.is-dir {
  background: #e8f1f6;
}

.icon.is-dir::before {
  content: '';
  width: 22px;
  height: 16px;
  border-radius: 3px 3px 2px 2px;
  background: #7eb6c9;
  box-shadow: 0 -3px 0 #5a9bb0;
}

.icon.is-file::before {
  content: '';
  width: 16px;
  height: 20px;
  border-radius: 2px;
  background: #c5d0da;
  box-shadow: inset 0 0 0 1px #a9b7c4;
}

.meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.meta strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

@media (max-width: 640px) {
  .file-row {
    flex-direction: column;
    align-items: stretch;
  }

  .row-actions {
    justify-content: flex-end;
  }
}
</style>
