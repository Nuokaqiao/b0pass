<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

const STORAGE_KEY = 'txtdata'
const draft = ref('')
const messages = ref([])
const status = ref('连接中…')
const logEl = ref(null)
let socket = null

function loadCache() {
  try {
    messages.value = JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]')
  } catch {
    messages.value = []
  }
}

function saveCache() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(messages.value))
}

function wsUrl() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${location.host}/ws`
}

async function scrollBottom() {
  await nextTick()
  if (logEl.value) {
    logEl.value.scrollTop = logEl.value.scrollHeight
  }
}

function pushMessage(text) {
  messages.value.push({
    key: Date.now() + Math.random(),
    val: text,
    time: new Date().toLocaleString(),
  })
  saveCache()
  scrollBottom()
}

function connect() {
  if (!window.WebSocket) {
    status.value = '当前浏览器不支持 WebSocket'
    return
  }
  socket = new WebSocket(wsUrl())
  socket.onopen = () => {
    status.value = '已连接'
  }
  socket.onclose = () => {
    status.value = '连接已断开'
  }
  socket.onerror = () => {
    status.value = '连接异常'
  }
  socket.onmessage = (evt) => {
    const chunks = String(evt.data || '').split('\n').filter(Boolean)
    for (const chunk of chunks) {
      pushMessage(chunk)
    }
  }
}

function send() {
  const text = draft.value.trim()
  if (!text) return
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    status.value = '未连接，无法发送'
    return
  }
  socket.send(text)
  draft.value = ''
}

function clearAll() {
  if (!confirm('清空本地内容记录？')) return
  messages.value = []
  saveCache()
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
    status.value = '已复制'
  } catch {
    const el = document.createElement('textarea')
    el.value = text
    document.body.appendChild(el)
    el.select()
    document.execCommand('copy')
    document.body.removeChild(el)
    status.value = '已复制'
  }
}

function removeAt(index) {
  messages.value.splice(index, 1)
  saveCache()
}

onMounted(() => {
  loadCache()
  connect()
  scrollBottom()
})

onBeforeUnmount(() => {
  if (socket) socket.close()
})
</script>

<template>
  <div class="page">
    <header class="topbar shell">
      <RouterLink class="back-link" :to="{ name: 'home' }">← 首页</RouterLink>
      <div class="brand">传内容</div>
      <span class="muted tip">{{ status }}</span>
    </header>

    <main class="shell text-layout">
      <section class="panel composer">
        <textarea
          v-model="draft"
          rows="5"
          placeholder="粘贴文字或链接，同步到其他设备…"
          @keydown.meta.enter.prevent="send"
          @keydown.ctrl.enter.prevent="send"
        />
        <div class="composer-actions">
          <button class="btn btn-ghost" type="button" @click="clearAll">清空记录</button>
          <button class="btn btn-primary" type="button" @click="send">发送</button>
        </div>
      </section>

      <section ref="logEl" class="panel message-panel">
        <div v-if="!messages.length" class="empty">还没有消息，发送一条试试</div>
        <article v-for="(msg, index) in messages" :key="msg.key" class="msg">
          <div class="msg-bar">
            <small class="muted">{{ msg.time || ('#' + msg.key) }}</small>
            <div class="msg-actions">
              <button type="button" @click="copyText(msg.val)">复制</button>
              <button type="button" @click="removeAt(index)">删除</button>
            </div>
          </div>
          <pre class="msg-body">{{ msg.val }}</pre>
        </article>
      </section>
    </main>
  </div>
</template>

<style scoped>
.tip {
  font-size: 0.85rem;
}

.text-layout {
  display: grid;
  gap: 16px;
}

.composer {
  padding: 16px;
}

.composer textarea {
  width: 100%;
  border: 1px solid var(--line);
  border-radius: 14px;
  padding: 14px;
  resize: vertical;
  min-height: 120px;
  background: #fbfcfd;
}

.composer textarea:focus {
  outline: 2px solid rgba(15, 110, 140, 0.25);
  border-color: var(--brand);
}

.composer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.message-panel {
  padding: 8px 12px 16px;
  max-height: min(55vh, 560px);
  overflow: auto;
}

.msg {
  padding: 12px 8px;
  border-bottom: 1px solid var(--line);
}

.msg:last-child {
  border-bottom: none;
}

.msg-bar {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.msg-actions {
  display: flex;
  gap: 8px;
}

.msg-actions button {
  border: none;
  background: transparent;
  color: var(--brand);
  cursor: pointer;
  font-weight: 500;
  padding: 0;
}

.msg-body {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: inherit;
  font-size: 0.98rem;
}
</style>
