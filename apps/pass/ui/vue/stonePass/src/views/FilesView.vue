<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import {
  addNode,
  deleteNode,
  downloadUrl,
  fetchFileList,
  setFileExpire,
  staticFileUrl,
  uploadFile,
} from '@/api/pass'

const currentPath = ref('/')
const items = ref([])
const loading = ref(false)
const error = ref('')
const uploading = ref(false)
const progress = ref(0)
const uploadName = ref('')
const fileInput = ref(null)
const folderInput = ref(null)
const expireAmountInput = ref(null)
const toast = ref('')
const toastError = ref(false)
const dragging = ref(false)
const showNewFolder = ref(false)
const newFolderName = ref('')
const copiedPath = ref('')
const uploadExpireAmount = ref('') // 空或 0 = 不过期
const uploadExpireUnit = ref('h') // m | h | d
/** @type {import('vue').Ref<Array<{ file: File, relativePath: string }>>} */
const pendingUploadFiles = ref([])
/** @type {import('vue').Ref<Set<string>>} */
const selectedPaths = ref(new Set())

const SIZE_WARN_BYTES = 1024 * 1024 * 1024 // 1 GiB

const expireEditor = ref(null) // { path, name, amount, unit } | null

const expireUnits = [
  { value: 'm', label: '分钟' },
  { value: 'h', label: '小时' },
  { value: 'd', label: '天' },
]

const pendingUploadBytes = computed(() =>
  pendingUploadFiles.value.reduce((sum, it) => sum + (Number(it.file?.size) || 0), 0),
)

const pendingOverSizeWarn = computed(() => pendingUploadBytes.value > SIZE_WARN_BYTES)

const pendingUploadSummary = computed(() => {
  const files = pendingUploadFiles.value
  if (!files.length) return ''
  const sizeLabel = formatBytes(pendingUploadBytes.value)
  if (files.length === 1) return `${files[0].relativePath || files[0].file.name}（${sizeLabel}）`
  const names = files
    .slice(0, 3)
    .map((f) => f.relativePath || f.file.name)
    .join('、')
  const head = files.length > 3 ? `${names} 等 ${files.length} 个文件` : `${names}（共 ${files.length} 个）`
  return `${head} · ${sizeLabel}`
})

watch(
  () => pendingUploadFiles.value.length,
  async (n) => {
    if (n > 0) {
      await nextTick()
      expireAmountInput.value?.focus?.()
    }
  },
)

let dragDepth = 0
let toastTimer = null
let copiedTimer = null

function unitSeconds(unit) {
  if (unit === 'm') return 60
  if (unit === 'd') return 86400
  return 3600
}

function expireUnixFromAmount(amount, unit) {
  const n = Number(amount)
  if (!Number.isFinite(n) || n <= 0) return 0
  return Math.floor(Date.now() / 1000) + Math.round(n * unitSeconds(unit))
}

function formatExpireLeft(seconds) {
  const s = Number(seconds) || 0
  if (s <= 0) return '已到期'
  if (s < 3600) return `${Math.ceil(s / 60)} 分钟后过期`
  if (s < 86400) return `${(s / 3600).toFixed(1)} 小时后过期`
  return `${(s / 86400).toFixed(1)} 天后过期`
}

function suggestExpireFields(expireUnix) {
  const left = Math.max(0, Number(expireUnix) - Math.floor(Date.now() / 1000))
  if (left <= 0) return { amount: '', unit: 'h' }
  if (left % 86400 === 0) return { amount: String(left / 86400), unit: 'd' }
  if (left % 3600 === 0) return { amount: String(left / 3600), unit: 'h' }
  if (left >= 86400) return { amount: String(Math.max(1, Math.round(left / 86400))), unit: 'd' }
  if (left >= 3600) return { amount: String(Math.max(1, Math.round(left / 3600))), unit: 'h' }
  return { amount: String(Math.max(1, Math.round(left / 60))), unit: 'm' }
}

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

const sortedItems = computed(() => {
  return [...items.value].sort((a, b) => {
    if (a.type === 'dir' && b.type !== 'dir') return -1
    if (a.type !== 'dir' && b.type === 'dir') return 1
    return String(a.name).localeCompare(String(b.name), 'zh')
  })
})

const selectedCount = computed(() => {
  let n = 0
  for (const it of sortedItems.value) {
    if (selectedPaths.value.has(it.path)) n += 1
  }
  return n
})

const allSelected = computed(
  () => sortedItems.value.length > 0 && selectedCount.value === sortedItems.value.length,
)

const selectAllIndeterminate = computed(
  () => selectedCount.value > 0 && !allSelected.value,
)

const dirSummary = computed(() => {
  let folderCount = 0
  let fileCount = 0
  let totalBytes = 0
  for (const item of items.value) {
    if (item.type === 'dir') {
      folderCount += 1
      continue
    }
    fileCount += 1
    const n = Number(item.size)
    if (Number.isFinite(n) && n > 0) totalBytes += n
  }
  return {
    folderCount,
    fileCount,
    totalLabel: formatBytes(totalBytes),
  }
})

function formatBytes(bytes) {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let n = bytes
  let i = 0
  while (n >= 1024 && i < units.length - 1) {
    n /= 1024
    i += 1
  }
  const digits = i === 0 ? 0 : n >= 10 ? 1 : 2
  return `${n.toFixed(digits)} ${units[i]}`
}

function summaryText() {
  const { folderCount, fileCount, totalLabel } = dirSummary.value
  const parts = []
  if (folderCount) parts.push(`${folderCount} 个文件夹`)
  parts.push(`${fileCount} 个文件`)
  parts.push(`共 ${totalLabel}`)
  return parts.join(' · ')
}

function normalizeDir(path) {
  if (!path || path === '/') return '/'
  return path.endsWith('/') ? path : `${path}/`
}

function showToast(msg, isError = false) {
  toast.value = msg
  toastError.value = isError
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => {
    if (toast.value === msg) {
      toast.value = ''
      toastError.value = false
    }
  }, isError ? 4000 : 2200)
}

async function loadList(path = currentPath.value) {
  loading.value = true
  error.value = ''
  selectedPaths.value = new Set()
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

function isSelected(path) {
  return selectedPaths.value.has(path)
}

function setSelected(path, checked) {
  const next = new Set(selectedPaths.value)
  if (checked) next.add(path)
  else next.delete(path)
  selectedPaths.value = next
}

function toggleSelectAll() {
  if (allSelected.value) {
    selectedPaths.value = new Set()
    return
  }
  selectedPaths.value = new Set(sortedItems.value.map((it) => it.path))
}

function clearSelection() {
  selectedPaths.value = new Set()
}

function openItem(item) {
  if (item.type !== 'dir') return
  loadList(item.path.endsWith('/') ? item.path : `${item.path}/`)
}

async function removeItem(item) {
  if (!confirm(`确定删除「${item.name}」？`)) return
  const targetPath = item.path
  try {
    await deleteNode(targetPath)
    items.value = items.value.filter((it) => it.path !== targetPath)
    const next = new Set(selectedPaths.value)
    next.delete(targetPath)
    selectedPaths.value = next
    showToast('已删除')
    await loadList()
  } catch (e) {
    showToast(e.message || '删除失败', true)
    await loadList()
  }
}

async function removeSelected() {
  const targets = sortedItems.value.filter((it) => selectedPaths.value.has(it.path))
  if (!targets.length) return
  if (!confirm(`确定删除已选的 ${targets.length} 项？`)) return
  let ok = 0
  let fail = 0
  for (const item of targets) {
    try {
      await deleteNode(item.path)
      ok += 1
    } catch {
      fail += 1
    }
  }
  if (fail) showToast(`已删除 ${ok} 项，失败 ${fail} 项`, true)
  else showToast(`已删除 ${ok} 项`)
  await loadList()
}

function openNewFolder() {
  showNewFolder.value = true
  newFolderName.value = ''
}

function cancelNewFolder() {
  showNewFolder.value = false
  newFolderName.value = ''
}

async function createFolder() {
  const name = newFolderName.value.trim()
  if (!name) return
  if (/[\\/]/.test(name)) {
    showToast('文件夹名称不能包含斜杠', true)
    return
  }
  const path = `${normalizeDir(currentPath.value)}${name}/`
  try {
    await addNode(path)
    showToast('已创建')
    cancelNewFolder()
    await loadList()
  } catch (e) {
    showToast(e.message || '创建失败', true)
  }
}

function triggerUpload() {
  const el = fileInput.value
  if (!el) {
    showToast('上传控件未就绪，请刷新页面', true)
    return
  }
  el.click()
}

function triggerFolderUpload() {
  const el = folderInput.value
  if (!el) {
    showToast('上传控件未就绪，请刷新页面', true)
    return
  }
  el.click()
}

/** 规范化相对路径，拒绝 .. */
function normalizeRelPath(raw) {
  const p = String(raw || '')
    .replace(/\\/g, '/')
    .replace(/^\/+/, '')
  if (!p) return ''
  const parts = p.split('/').filter(Boolean)
  if (parts.some((x) => x === '..')) {
    throw new Error(`非法路径：${raw}`)
  }
  return parts.join('/')
}

function toUploadItems(fileList) {
  return Array.from(fileList || []).map((file) => {
    const relativePath = normalizeRelPath(file.webkitRelativePath || file.name)
    return { file, relativePath: relativePath || file.name }
  })
}

function readEntries(reader) {
  return new Promise((resolve, reject) => {
    reader.readEntries(resolve, reject)
  })
}

function entryFile(entry) {
  return new Promise((resolve, reject) => {
    entry.file(resolve, reject)
  })
}

async function walkDirEntry(entry, prefix, out) {
  if (entry.isFile) {
    const file = await entryFile(entry)
    const relativePath = normalizeRelPath(prefix ? `${prefix}/${file.name}` : file.name)
    if (relativePath) out.push({ file, relativePath })
    return
  }
  if (!entry.isDirectory) return
  const nextPrefix = prefix ? `${prefix}/${entry.name}` : entry.name
  const reader = entry.createReader()
  let batch = await readEntries(reader)
  while (batch.length) {
    for (const child of batch) {
      await walkDirEntry(child, nextPrefix, out)
    }
    batch = await readEntries(reader)
  }
}

async function collectDroppedItems(dataTransfer) {
  const items = dataTransfer?.items
  if (items && items.length) {
    const entries = []
    for (let i = 0; i < items.length; i += 1) {
      const entry = items[i].webkitGetAsEntry?.()
      if (entry) entries.push(entry)
    }
    if (entries.length) {
      const out = []
      for (const entry of entries) {
        await walkDirEntry(entry, '', out)
      }
      return out
    }
  }
  return toUploadItems(dataTransfer?.files)
}

function queueUploadItems(items) {
  let list = []
  try {
    list = (items || []).map((it) => ({
      file: it.file,
      relativePath: normalizeRelPath(it.relativePath || it.file?.name || ''),
    })).filter((it) => it.file && it.relativePath)
  } catch (e) {
    showToast(e.message || '路径无效', true)
    return
  }
  if (!list.length) {
    showToast('未选择到文件', true)
    return
  }
  if (uploading.value) {
    showToast('正在上传中，请稍候', true)
    return
  }
  pendingUploadFiles.value = list
  uploadExpireAmount.value = ''
  uploadExpireUnit.value = 'h'
}

function cancelPendingUpload() {
  pendingUploadFiles.value = []
  uploadExpireAmount.value = ''
  uploadExpireUnit.value = 'h'
}

async function confirmPendingUpload() {
  const files = pendingUploadFiles.value
  if (!files.length) return
  if (pendingOverSizeWarn.value) {
    const ok = confirm(
      `所选内容约 ${formatBytes(pendingUploadBytes.value)}，已超过 1 GB，上传可能较久。确定继续？`,
    )
    if (!ok) return
  }
  const expireUnix = expireUnixFromAmount(uploadExpireAmount.value, uploadExpireUnit.value)
  pendingUploadFiles.value = []
  await uploadFiles(files, expireUnix)
}

async function uploadFiles(itemList, expireUnix = 0) {
  const items = Array.isArray(itemList) ? itemList : []
  if (!items.length) {
    showToast('未选择到文件', true)
    return
  }
  if (uploading.value) {
    showToast('正在上传中，请稍候', true)
    return
  }
  uploading.value = true
  progress.value = 0
  const baseDir = normalizeDir(currentPath.value)
  try {
    for (let i = 0; i < items.length; i += 1) {
      const { file, relativePath } = items[i]
      const parts = relativePath.split('/')
      const fileName = parts.pop()
      const sub = parts.length ? `${parts.join('/')}/` : ''
      const targetDir = `${baseDir}${sub}`
      uploadName.value =
        items.length > 1 ? `${relativePath}（${i + 1}/${items.length}）` : relativePath
      progress.value = 0
      await uploadFile(
        targetDir,
        file,
        (p) => {
          progress.value = p
        },
        expireUnix,
        fileName,
      )
    }
    showToast(items.length > 1 ? `已上传 ${items.length} 个文件` : '上传完成')
    await loadList()
  } catch (err) {
    showToast(err.message || '上传失败', true)
  } finally {
    uploading.value = false
    progress.value = 0
    uploadName.value = ''
  }
}

function onFilesSelected(e) {
  const input = e.target
  const files = Array.from(input.files || [])
  input.value = ''
  queueUploadItems(toUploadItems(files))
}

function onFolderSelected(e) {
  const input = e.target
  const files = Array.from(input.files || [])
  input.value = ''
  queueUploadItems(toUploadItems(files))
}

function onDragEnter(e) {
  e.preventDefault()
  dragDepth += 1
  dragging.value = true
}

function onDragOver(e) {
  e.preventDefault()
}

function onDragLeave(e) {
  e.preventDefault()
  dragDepth = Math.max(0, dragDepth - 1)
  if (dragDepth === 0) dragging.value = false
}

async function onDrop(e) {
  e.preventDefault()
  dragDepth = 0
  dragging.value = false
  try {
    const items = await collectDroppedItems(e.dataTransfer)
    queueUploadItems(items)
  } catch (err) {
    showToast(err.message || '读取拖拽内容失败', true)
  }
}

function previewUrl(item) {
  if (item.type === 'img') return staticFileUrl(item.path)
  return ''
}

function typeLabel(item) {
  if (item.type === 'dir') return '文件夹'
  if (item.type === 'img') return '图片'
  if (item.type === 'vod') return '视频'
  if (item.type === 'pdf') return 'PDF'
  return item.ext || '文件'
}

async function copyLink(item) {
  const url = `${location.origin}${downloadUrl(item.path)}`
  try {
    await navigator.clipboard.writeText(url)
  } catch {
    const el = document.createElement('textarea')
    el.value = url
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
  }
  copiedPath.value = item.path
  if (copiedTimer) clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => {
    if (copiedPath.value === item.path) copiedPath.value = ''
  }, 1600)
}

function openExpireEditor(item) {
  if (expireEditor.value?.path === item.path) {
    expireEditor.value = null
    return
  }
  const suggested = suggestExpireFields(item.expire)
  expireEditor.value = {
    path: item.path,
    name: item.name,
    amount: suggested.amount,
    unit: suggested.unit,
  }
}

function cancelExpireEditor() {
  expireEditor.value = null
}

async function saveExpireEditor(clear = false) {
  const ed = expireEditor.value
  if (!ed) return
  const unix = clear ? 0 : expireUnixFromAmount(ed.amount, ed.unit)
  if (!clear) {
    const n = Number(ed.amount)
    if (!Number.isFinite(n) || n <= 0) {
      showToast('请输入大于 0 的数字，或点「不过期」', true)
      return
    }
  }
  const unitLabel = expireUnits.find((u) => u.value === ed.unit)?.label || ''
  try {
    await setFileExpire(ed.path, unix)
    showToast(unix > 0 ? `已设置 ${ed.amount} ${unitLabel}后过期` : '已取消过期')
    expireEditor.value = null
    await loadList()
  } catch (e) {
    showToast(e.message || '设置失败', true)
  }
}

onMounted(() => loadList('/'))
</script>

<template>
  <div
    class="page files-page"
    @dragenter="onDragEnter"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <div v-if="dragging" class="drop-mask" aria-hidden="true">
      <div class="drop-card">松开以上传文件或文件夹到当前目录</div>
    </div>

    <header class="topbar shell">
      <RouterLink class="back-link" :to="{ name: 'home' }">← 首页</RouterLink>
      <div class="brand">传文件</div>
      <span class="muted tip">高速传输</span>
    </header>

    <main class="shell files-layout">
      <section class="panel files-panel">
        <div class="panel-head">
          <nav class="crumbs" aria-label="路径">
            <button
              v-for="(c, i) in crumbs"
              :key="c.path"
              type="button"
              class="crumb"
              @click="loadList(c.path === '/' ? '/' : c.path + '/')"
            >
              <span v-if="i" class="sep">/</span>{{ c.name }}
            </button>
          </nav>
          <div class="actions">
            <button class="btn btn-ghost btn-sm" type="button" @click="openNewFolder">新建文件夹</button>
            <button class="btn btn-ghost btn-sm" type="button" :disabled="loading" @click="loadList()">
              刷新
            </button>
            <label class="btn btn-ghost btn-sm upload-label" :class="{ disabled: uploading || pendingUploadFiles.length }">
              上传文件夹
              <input
                ref="folderInput"
                type="file"
                multiple
                webkitdirectory
                directory
                class="file-input"
                :disabled="uploading || pendingUploadFiles.length > 0"
                @change="onFolderSelected"
              />
            </label>
            <label class="btn btn-primary btn-sm upload-label" :class="{ disabled: uploading || pendingUploadFiles.length }">
              {{ uploading ? `上传中 ${progress}%` : '上传文件' }}
              <input
                ref="fileInput"
                type="file"
                multiple
                class="file-input"
                :disabled="uploading || pendingUploadFiles.length > 0"
                @change="onFilesSelected"
              />
            </label>
          </div>
        </div>

        <div v-if="showNewFolder" class="new-folder">
          <input
            v-model="newFolderName"
            type="text"
            placeholder="输入文件夹名称"
            @keydown.enter.prevent="createFolder"
            @keydown.esc.prevent="cancelNewFolder"
          />
          <button class="btn btn-primary btn-sm" type="button" @click="createFolder">创建</button>
          <button class="btn btn-ghost btn-sm" type="button" @click="cancelNewFolder">取消</button>
        </div>

        <div v-if="uploading" class="upload-bar">
          <div class="upload-meta">
            <strong>正在上传</strong>
            <span class="muted">{{ uploadName }}</span>
          </div>
          <div class="progress-track">
            <div class="progress-fill" :style="{ width: `${progress}%` }"></div>
          </div>
        </div>

        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="!loading" class="dir-summary muted">{{ summaryText() }}</p>

        <div v-if="!loading && sortedItems.length" class="select-bar">
          <label class="select-all">
            <input
              type="checkbox"
              :checked="allSelected"
              :indeterminate="selectAllIndeterminate"
              @change="toggleSelectAll"
            />
            <span>全选</span>
          </label>
          <span v-if="selectedCount" class="muted">已选 {{ selectedCount }} 项</span>
          <div class="select-actions">
            <button
              v-if="selectedCount"
              class="btn btn-ghost btn-sm"
              type="button"
              @click="clearSelection"
            >
              取消选择
            </button>
            <button
              v-if="selectedCount"
              class="btn btn-danger btn-sm"
              type="button"
              @click="removeSelected"
            >
              删除所选
            </button>
          </div>
        </div>

        <div v-if="loading" class="empty">加载中…</div>
        <div v-else-if="!sortedItems.length" class="empty-box">
          <p>当前目录为空</p>
          <p class="muted">点击「上传文件 / 上传文件夹」，或把文件、文件夹拖到此页面</p>
          <div class="empty-actions">
            <button class="btn btn-primary" type="button" :disabled="uploading" @click="triggerUpload">
              选择文件
            </button>
            <button class="btn btn-ghost" type="button" :disabled="uploading" @click="triggerFolderUpload">
              选择文件夹
            </button>
          </div>
        </div>

        <ul v-else class="file-list">
          <li
            v-for="item in sortedItems"
            :key="item.path"
            class="file-card"
            :class="{
              'is-editing-expire': expireEditor?.path === item.path,
              'is-selected': isSelected(item.path),
            }"
          >
            <div class="file-row">
              <label class="row-check" @click.stop>
                <input
                  type="checkbox"
                  :checked="isSelected(item.path)"
                  @change="setSelected(item.path, $event.target.checked)"
                />
              </label>
              <button
                v-if="item.type === 'dir'"
                class="file-main"
                type="button"
                @click="openItem(item)"
              >
                <span class="icon is-dir"></span>
                <span class="meta">
                  <strong>{{ item.name }}</strong>
                  <small class="muted">
                    {{ typeLabel(item) }} · {{ item.date }}
                    <template v-if="item.expire">
                      · <span class="expire-tag">{{ formatExpireLeft(item.expireLeft) }}</span>
                    </template>
                  </small>
                </span>
              </button>
              <div v-else class="file-main is-static">
                <img v-if="previewUrl(item)" class="thumb" :src="previewUrl(item)" alt="" />
                <span v-else class="icon is-file"></span>
                <span class="meta">
                  <strong>{{ item.name }}</strong>
                  <small class="muted">
                    {{ typeLabel(item) }} · {{ item.sizes }} · {{ item.date }}
                    <template v-if="item.expire">
                      · <span class="expire-tag">{{ formatExpireLeft(item.expireLeft) }}</span>
                    </template>
                  </small>
                </span>
              </div>
              <div class="row-actions">
                <a
                  v-if="item.type !== 'dir'"
                  class="btn btn-primary btn-sm"
                  :href="downloadUrl(item.path)"
                  target="_blank"
                  rel="noopener"
                >下载</a>
                <button
                  class="btn btn-ghost btn-sm"
                  type="button"
                  @click="openExpireEditor(item)"
                >
                  过期
                </button>
                <button
                  v-if="item.type !== 'dir'"
                  class="btn btn-ghost btn-sm"
                  type="button"
                  @click="copyLink(item)"
                >
                  {{ copiedPath === item.path ? '已复制' : '复制链接' }}
                </button>
                <button class="btn btn-danger btn-sm" type="button" @click="removeItem(item)">删除</button>
              </div>
            </div>
            <div v-if="expireEditor?.path === item.path" class="expire-inline">
              <input
                v-model="expireEditor.amount"
                class="expire-amount"
                type="number"
                min="0"
                step="1"
                placeholder="数字"
                @keydown.enter.prevent="saveExpireEditor(false)"
                @keydown.esc.prevent="cancelExpireEditor"
              />
              <select v-model="expireEditor.unit">
                <option v-for="u in expireUnits" :key="u.value" :value="u.value">
                  {{ u.label }}
                </option>
              </select>
              <button class="btn btn-primary btn-sm" type="button" @click="saveExpireEditor(false)">确定</button>
              <button class="btn btn-ghost btn-sm" type="button" @click="saveExpireEditor(true)">不过期</button>
              <button class="btn btn-ghost btn-sm" type="button" @click="cancelExpireEditor">取消</button>
            </div>
          </li>
        </ul>
      </section>
    </main>

    <Teleport to="body">
      <div
        v-if="pendingUploadFiles.length"
        class="modal-mask"
        role="dialog"
        aria-modal="true"
        aria-label="确认上传"
        @click.self="cancelPendingUpload"
        @keydown.esc.prevent="cancelPendingUpload"
      >
        <div class="modal-card upload-confirm-modal" tabindex="-1" @keydown.esc.prevent="cancelPendingUpload">
          <h3 class="modal-title">确认上传</h3>
          <p class="modal-desc muted">{{ pendingUploadSummary }}</p>
          <p v-if="pendingOverSizeWarn" class="size-warn">
            总大小超过 1 GB，上传可能较久，请确认网络与电量。
          </p>
          <div class="expire-select modal-expire">
            <span>过期</span>
            <input
              ref="expireAmountInput"
              v-model="uploadExpireAmount"
              class="expire-amount"
              type="number"
              min="0"
              step="1"
              placeholder="空=不过期"
              @keydown.enter.prevent="confirmPendingUpload"
              @keydown.esc.prevent="cancelPendingUpload"
            />
            <select v-model="uploadExpireUnit">
              <option v-for="u in expireUnits" :key="u.value" :value="u.value">
                {{ u.label }}
              </option>
            </select>
          </div>
          <div class="upload-confirm-actions">
            <button class="btn btn-primary" type="button" @click="confirmPendingUpload">开始上传</button>
            <button class="btn btn-ghost" type="button" @click="cancelPendingUpload">取消</button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div
        v-if="toast"
        class="float-toast"
        :class="{ 'is-error': toastError }"
        role="status"
      >
        {{ toast }}
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.files-page {
  position: relative;
}

.tip {
  font-size: 0.85rem;
}

.files-layout {
  padding-bottom: 28px;
}

.files-panel {
  padding: 14px;
}

.panel-head {
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
  gap: 2px;
  align-items: center;
  min-width: 0;
}

.crumb {
  border: none;
  background: transparent;
  color: var(--brand);
  cursor: pointer;
  padding: 4px 2px;
  font-weight: 600;
}

.sep {
  margin-right: 2px;
  color: var(--muted);
  font-weight: 400;
}

.actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.expire-select {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
}

.expire-amount {
  width: 88px;
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 7px 10px;
  background: rgba(255, 255, 255, 0.9);
  color: var(--ink);
}

.expire-select select,
.expire-editor select {
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 7px 10px;
  background: rgba(255, 255, 255, 0.9);
  color: var(--ink);
}

.expire-tag {
  color: #a15c00;
  font-weight: 600;
}

.upload-label {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  margin: 0;
}

.upload-label.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

.file-input {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  border: 0;
}

.btn-sm {
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.88rem;
}

.new-folder {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.new-folder input {
  flex: 1;
  min-width: 160px;
  border: 1px solid var(--line);
  border-radius: 10px;
  padding: 10px 12px;
  background: rgba(251, 252, 253, 0.95);
}

.new-folder input:focus {
  outline: 2px solid rgba(15, 110, 140, 0.25);
  border-color: var(--brand);
}

.upload-bar {
  margin-bottom: 12px;
  padding: 12px;
  border-radius: 12px;
  background: var(--brand-soft);
  border: 1px solid rgba(15, 110, 140, 0.18);
}

.modal-mask {
  position: fixed;
  inset: 0;
  z-index: 1200;
  display: grid;
  place-items: center;
  padding: 20px;
  background: rgba(20, 36, 48, 0.42);
  backdrop-filter: blur(2px);
}

.modal-card {
  width: min(420px, 100%);
  padding: 20px 22px;
  border-radius: 16px;
  background: #fbfcfd;
  border: 1px solid var(--line);
  box-shadow: 0 18px 48px rgba(20, 40, 55, 0.22);
}

.modal-title {
  margin: 0 0 8px;
  font-size: 1.15rem;
  font-weight: 700;
  color: var(--ink);
}

.modal-desc {
  margin: 0 0 16px;
  word-break: break-all;
  line-height: 1.45;
}

.modal-expire {
  margin-bottom: 18px;
}

.upload-confirm-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
}

.upload-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: baseline;
  margin-bottom: 8px;
  font-size: 0.92rem;
}

.progress-track {
  height: 8px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.7);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #0f6e8c, #1a8fb0);
  transition: width 0.15s ease;
}

.error {
  background: #fdecec;
  color: var(--danger);
  margin: 0 0 10px;
  padding: 10px 12px;
  border-radius: 10px;
}

.dir-summary {
  margin: 0 0 12px;
  font-size: 0.88rem;
  min-height: 1.2em;
}

.float-toast {
  position: fixed;
  top: 24px;
  left: 50%;
  z-index: 1000;
  transform: translateX(-50%);
  max-width: min(90vw, 420px);
  padding: 12px 18px;
  border-radius: 12px;
  background: rgba(15, 110, 140, 0.95);
  color: #fff;
  font-size: 0.95rem;
  font-weight: 600;
  box-shadow: 0 12px 32px rgba(26, 35, 50, 0.22);
  pointer-events: none;
  animation: toast-in 0.2s ease;
}

.float-toast.is-error {
  background: rgba(180, 35, 24, 0.95);
}

@keyframes toast-in {
  from {
    opacity: 0;
    transform: translate(-50%, -8px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}

.empty-box {
  display: grid;
  gap: 8px;
  place-items: center;
  padding: 48px 20px;
  text-align: center;
}

.empty-box p {
  margin: 0;
}

.empty-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
  margin-top: 4px;
}

.size-warn {
  margin: 0 0 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(180, 110, 20, 0.12);
  border: 1px solid rgba(180, 110, 20, 0.28);
  color: #8a4b00;
  font-size: 0.9rem;
  line-height: 1.4;
}

.file-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 10px;
}

.select-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px 14px;
  margin: 0 0 12px;
  padding: 8px 10px;
  border-radius: 12px;
  background: rgba(232, 244, 248, 0.65);
  border: 1px solid rgba(15, 110, 140, 0.12);
}

.select-all {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  user-select: none;
}

.select-all input,
.row-check input {
  width: 16px;
  height: 16px;
  accent-color: var(--brand);
  cursor: pointer;
}

.select-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-left: auto;
}

.row-check {
  display: grid;
  place-items: center;
  flex-shrink: 0;
  padding: 4px 2px;
  cursor: pointer;
}

.file-card.is-selected {
  border-color: rgba(15, 110, 140, 0.4);
  background: rgba(232, 244, 248, 0.9);
}

.file-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px 14px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.72);
}

.file-row {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: space-between;
}

.file-card.is-editing-expire {
  border-color: rgba(15, 110, 140, 0.35);
  background: rgba(232, 244, 248, 0.85);
}

.expire-inline {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  padding-top: 2px;
}

.expire-inline .expire-amount {
  width: 88px;
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
  color: inherit;
}

.file-main.is-static {
  cursor: default;
}

.icon,
.thumb {
  width: 44px;
  height: 44px;
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
  flex-wrap: wrap;
  gap: 6px;
  flex-shrink: 0;
  justify-content: flex-end;
}

.drop-mask {
  position: fixed;
  inset: 0;
  z-index: 40;
  display: grid;
  place-items: center;
  background: rgba(15, 110, 140, 0.18);
  backdrop-filter: blur(2px);
  pointer-events: none;
}

.drop-card {
  padding: 28px 36px;
  border-radius: 18px;
  border: 2px dashed rgba(15, 110, 140, 0.55);
  background: rgba(255, 255, 255, 0.92);
  color: var(--brand-deep);
  font-size: 1.15rem;
  font-weight: 700;
  box-shadow: var(--shadow);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 720px) {
  .file-row {
    flex-direction: column;
    align-items: stretch;
  }

  .row-actions {
    justify-content: flex-start;
  }
}
</style>
