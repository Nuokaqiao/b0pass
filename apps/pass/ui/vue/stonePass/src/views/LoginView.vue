<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login as doLogin } from '@/api/auth'

const route = useRoute()
const router = useRouter()

const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!password.value) {
    error.value = '请输入访问口令'
    return
  }
  loading.value = true
  try {
    const result = await doLogin(password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    if (result.enabled === false) {
      await router.replace(redirect || '/')
      return
    }
    await router.replace(redirect || '/')
  } catch (e) {
    error.value = e.message || '登录失败'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  password.value = ''
})
</script>

<template>
  <div class="page login-page">
    <main class="shell login-main">
      <section class="login-card panel">
        <p class="kicker">STONEPASS</p>
        <h1 class="title">石头记</h1>
        <p class="muted lead">输入访问口令后继续</p>

        <form class="login-form" @submit.prevent="submit">
          <label class="field">
            <span class="muted">口令</span>
            <input
              v-model="password"
              type="password"
              autocomplete="current-password"
              placeholder="请输入口令"
              :disabled="loading"
            />
          </label>
          <p v-if="error" class="error">{{ error }}</p>
          <button class="btn btn-primary login-btn" type="submit" :disabled="loading">
            {{ loading ? '登录中…' : '进入' }}
          </button>
        </form>
      </section>
    </main>
  </div>
</template>

<style scoped>
.login-main {
  flex: 1;
  display: grid;
  place-items: center;
  padding: 32px 16px 48px;
}

.login-card {
  width: min(420px, 100%);
  padding: 28px 26px 24px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid var(--line);
  box-shadow: var(--shadow);
}

.kicker {
  margin: 0 0 8px;
  font-size: 0.75rem;
  font-weight: 600;
  letter-spacing: 0.18em;
  color: var(--brand);
}

.title {
  margin: 0 0 6px;
  font-size: 1.8rem;
  font-weight: 700;
}

.lead {
  margin: 0 0 22px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.9rem;
}

.field input {
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 12px 14px;
  background: rgba(251, 252, 253, 0.95);
  color: var(--ink);
}

.field input:focus {
  outline: 2px solid rgba(15, 110, 140, 0.25);
  border-color: var(--brand);
}

.login-btn {
  width: 100%;
  margin-top: 4px;
  padding: 12px 16px;
  border-radius: 12px;
}

.error {
  margin: 0;
  color: var(--danger);
  font-size: 0.9rem;
}
</style>
