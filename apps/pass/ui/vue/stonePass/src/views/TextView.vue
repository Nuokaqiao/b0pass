<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

const STORAGE_KEY = 'txtdata'
const MAX_ITEMS = 100
const URL_RE = /https?:\/\/[^\s<>"{}|\\^`[\]]+/gi

const draft = ref('')
const messages = ref([])
const status = ref('connecting')
const statusText = ref('连接中…')
const logEl = ref(null)
const copiedKey = ref(null)

let socket = null
let reconnectTimer = null
let reconnectAttempt = 0
let manualClose = false
let copiedTimer = null

const canSync = computed(() => status.value === 'open')
const newestFirst = computed(() => [...messages.value].reverse())

function loadCache() {
  try {
    const raw = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]')
    messages.value = Array.isArray(raw) ? raw.slice(-MAX_ITEMS) : []
  } catch {
    messages.value = []
  }
}

function saveCache() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(messages.value.slice(-MAX_ITEMS)))
}

function wsUrl() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${location.host}/ws`
}

function setStatus(next, text) {
  status.value = next
  statusText.value = text
}

async function scrollTop() {
  await nextTick()
  if (logEl.value) logEl.value.scrollTop = 0
}

function extractUrls(text) {
  return [...new Set((text.match(URL_RE) || []).map((u) => u.replace(/[),.;!?]+$/, '')))]
}

function pushMessage(text) {
  messages.value.push({
    key: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    val: text,
    time: new Date().toLocaleString(),
  })
  if (messages.value.length > MAX_ITEMS) {
    messages.value = messages.value.slice(-MAX_ITEMS)
  }
  saveCache()
  scrollTop()
}

function clearReconnect() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
}

function scheduleReconnect() {
  if (manualClose) return
  clearReconnect()
  const delay = Math.min(1000 * 2 ** reconnectAttempt, 15000)
  reconnectAttempt += 1
  setStatus('reconnect', `重连中…（${Math.round(delay / 1000)}s）`)
  reconnectTimer = setTimeout(connect, delay)
}

function connect() {
  if (!window.WebSocket) {
    setStatus('error', '当前浏览器不支持 WebSocket')
    return
  }
  clearReconnect()
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return
  }

  setStatus('connecting', reconnectAttempt ? '重连中…' : '连接中…')
  socket = new WebSocket(wsUrl())

  socket.onopen = () => {
    reconnectAttempt = 0
    setStatus('open', '已连接')
  }

  socket.onclose = () => {
    socket = null
    if (manualClose) {
      setStatus('closed', '已断开')
      return
    }
    scheduleReconnect()
  }

  socket.onerror = () => {
    // onclose 会继续处理重连
  }

  socket.onmessage = (evt) => {
    const chunks = String(evt.data || '').split('\n').filter(Boolean)
    for (const chunk of chunks) {
      pushMessage(chunk)
    }
  }
}

function sync() {
  const text = draft.value.trim()
  if (!text) return
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    setStatus(status.value === 'reconnect' ? 'reconnect' : 'closed', '未连接，无法同步')
    return
  }
  socket.send(text)
  draft.value = ''
}

function clearAll() {
  if (!messages.value.length) return
  if (!confirm('清空本地内容记录？')) return
  messages.value = []
  saveCache()
}

async function copyText(msg) {
  const text = msg.val
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    const el = document.createElement('textarea')
    el.value = text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
  }
  copiedKey.value = msg.key
  if (copiedTimer) clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => {
    if (copiedKey.value === msg.key) copiedKey.value = null
  }, 1600)
}

function removeAt(key) {
  const index = messages.value.findIndex((m) => m.key === key)
  if (index < 0) return
  messages.value.splice(index, 1)
  saveCache()
}

onMounted(() => {
  loadCache()
  connect()
})

onBeforeUnmount(() => {
  manualClose = true
  clearReconnect()
  if (copiedTimer) clearTimeout(copiedTimer)
  if (socket) socket.close()
})
</script>

<template>
  <div class="page text-page">
    <header class="topbar shell">
      <RouterLink class="back-link" :to="{ name: 'home' }">← 首页</RouterLink>
      <div class="brand">传内容</div>
      <span class="status" :data-state="status">{{ statusText }}</span>
    </header>

    <main class="shell text-layout">
      <section class="panel history-panel">
        <div class="panel-head">
          <h2>最近内容</h2>
          <button
            class="btn btn-ghost btn-sm"
            type="button"
            :disabled="!messages.length"
            @click="clearAll"
          >
            清空
          </button>
        </div>

        <div ref="logEl" class="history-scroll">
          <div v-if="!messages.length" class="empty">还没有内容，粘贴一段试试</div>

          <article v-for="msg in newestFirst" :key="msg.key" class="card">
            <div class="card-meta">
              <small class="muted">{{ msg.time }}</small>
              <button class="link-btn danger" type="button" @click="removeAt(msg.key)">删除</button>
            </div>

            <button class="card-body" type="button" :title="'点击复制'" @click="copyText(msg)">
              <pre>{{ msg.val }}</pre>
            </button>

            <div class="card-actions">
              <button
                class="btn btn-primary btn-sm"
                type="button"
                @click="copyText(msg)"
              >
                {{ copiedKey === msg.key ? '已复制' : '复制' }}
              </button>
              <a
                v-for="url in extractUrls(msg.val)"
                :key="url"
                class="btn btn-ghost btn-sm"
                :href="url"
                target="_blank"
                rel="noopener noreferrer"
              >
                打开链接
              </a>
            </div>
          </article>
        </div>
      </section>

      <section class="panel composer">
        <textarea
          v-model="draft"
          rows="4"
          placeholder="粘贴文字或链接，同步到其他设备…"
          @keydown.meta.enter.prevent="sync"
          @keydown.ctrl.enter.prevent="sync"
        />
        <div class="composer-actions">
          <span class="muted hint">Ctrl / ⌘ + Enter</span>
          <button
            class="btn btn-primary"
            type="button"
            :disabled="!canSync || !draft.trim()"
            @click="sync"
          >
            同步
          </button>
        </div>
      </section>
    </main>
  </div>
</template>

<style scoped>
.text-page {
  min-height: 100vh;
}

.text-layout {
  display: grid;
  grid-template-rows: 1fr auto;
  gap: 14px;
  min-height: calc(100vh - 96px);
  padding-bottom: 24px;
}

.status {
  font-size: 0.82rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.55);
  border: 1px solid var(--line);
}

.status[data-state='open'] {
  color: var(--text);
  background: rgba(232, 245, 238, 0.9);
  border-color: rgba(45, 106, 79, 0.25);
}

.status[data-state='connecting'],
.status[data-state='reconnect'] {
  color: #8a5a00;
  background: rgba(255, 244, 214, 0.95);
  border-color: rgba(196, 146, 42, 0.3);
}

.status[data-state='closed'],
.status[data-state='error'] {
  color: var(--danger);
  background: rgba(253, 236, 236, 0.95);
  border-color: rgba(180, 35, 24, 0.22);
}

.history-panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 14px 14px 10px;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.panel-head h2 {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
}

.history-scroll {
  flex: 1;
  min-height: 220px;
  max-height: min(58vh, 620px);
  overflow: auto;
  padding-right: 2px;
}

.card {
  padding: 14px;
  margin-bottom: 10px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.72);
}

.card:last-child {
  margin-bottom: 0;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}

.card-body {
  display: block;
  width: 100%;
  border: none;
  background: transparent;
  padding: 0;
  text-align: left;
  cursor: pointer;
  color: inherit;
}

.card-body pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  font-size: 0.98rem;
  line-height: 1.55;
}

.card-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.btn-sm {
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.88rem;
}

.link-btn {
  border: none;
  background: transparent;
  color: var(--brand);
  cursor: pointer;
  font-weight: 500;
  padding: 0;
}

.link-btn.danger {
  color: var(--danger);
}

.composer {
  padding: 14px;
  position: sticky;
  bottom: 12px;
}

.composer textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 14px;
  padding: 14px;
  resize: vertical;
  min-height: 110px;
  background: rgba(251, 252, 253, 0.95);
}

.composer textarea:focus {
  outline: 2px solid rgba(15, 110, 140, 0.25);
  border-color: var(--brand);
}

.composer-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: 12px;
}

.hint {
  font-size: 0.82rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .text-layout {
    min-height: calc(100vh - 88px);
  }

  .history-scroll {
    max-height: none;
    flex: 1;
  }

  .composer {
    bottom: 8px;
  }
}
</style>
