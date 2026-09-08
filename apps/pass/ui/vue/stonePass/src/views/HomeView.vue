<script setup>
import { computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { getCachedAuthEnabled, getToken, logout } from '@/api/auth'

const router = useRouter()
const showLogout = computed(() => getCachedAuthEnabled() !== false && !!getToken())

function onLogout() {
  logout()
  router.replace({ name: 'login' })
}
</script>

<template>
  <div class="page home">
    <header class="topbar shell home-top">
      <div class="brand">石头记 <span>StonePass</span></div>
      <button v-if="showLogout" class="btn btn-ghost btn-sm logout-btn" type="button" @click="onLogout">
        退出
      </button>
    </header>

    <main class="shell home-main">
      <div class="home-intro">
        <p class="home-kicker">CROSS TRANSFER</p>
        <p class="home-lead muted">跨端文件与内容传递</p>
      </div>

      <div class="hero-actions">
        <RouterLink class="hero-btn hero-file" :to="{ name: 'files' }">
          <i class="hero-mark" aria-hidden="true"></i>
          <strong>传文件</strong>
          <span>高速传输</span>
        </RouterLink>

        <RouterLink class="hero-btn hero-text" :to="{ name: 'text' }">
          <i class="hero-mark" aria-hidden="true"></i>
          <strong>传内容</strong>
          <span>文字、链接一键同步</span>
        </RouterLink>
      </div>
    </main>
  </div>
</template>

<style scoped>
.home-top {
  width: min(720px, calc(100% - 32px));
  margin: 0 auto;
  padding-top: 28px;
}

.logout-btn {
  margin-left: auto;
}

.home-main {
  width: min(720px, calc(100% - 32px));
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding-bottom: 64px;
}

.home-intro {
  margin-bottom: 28px;
  animation: rise 0.55s ease both;
}

.home-kicker {
  margin: 0 0 8px;
  font-size: 0.78rem;
  font-weight: 600;
  letter-spacing: 0.18em;
  color: var(--brand);
}

.home-lead {
  margin: 0;
  font-size: 1.05rem;
}

.hero-actions {
  display: grid;
  gap: 18px;
}

.hero-btn {
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  min-height: 148px;
  padding: 28px 32px;
  border-radius: 22px;
  border: 1px solid rgba(255, 255, 255, 0.28);
  box-shadow: var(--shadow);
  transition: transform 0.18s ease, box-shadow 0.18s ease;
  animation: rise 0.65s ease both;
}

.hero-btn:nth-child(2) {
  animation-delay: 0.08s;
}

.hero-btn::after {
  content: '';
  position: absolute;
  inset: auto -20% -40% auto;
  width: 180px;
  height: 180px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.14);
}

.hero-mark {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.28);
  margin-bottom: 4px;
}

.hero-file .hero-mark {
  box-shadow: inset 8px 0 0 rgba(255, 255, 255, 0.25);
}

.hero-text .hero-mark {
  background:
    linear-gradient(rgba(255, 255, 255, 0.35), rgba(255, 255, 255, 0.35)) center/18px 2px no-repeat,
    linear-gradient(rgba(255, 255, 255, 0.35), rgba(255, 255, 255, 0.35)) center 60%/14px 2px no-repeat,
    rgba(255, 255, 255, 0.2);
}

.hero-btn strong {
  font-size: clamp(1.8rem, 4vw, 2.4rem);
  font-weight: 700;
  letter-spacing: 0.04em;
}

.hero-btn span {
  font-size: 1rem;
  opacity: 0.86;
}

.hero-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 22px 48px rgba(26, 35, 50, 0.14);
}

.hero-file {
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.12), transparent 42%),
    linear-gradient(135deg, #0f6e8c 0%, #1a8fb0 55%, #1485a8 100%);
  color: #fff;
}

.hero-text {
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.12), transparent 42%),
    linear-gradient(135deg, #2d6a4f 0%, #3d8f68 55%, #40916c 100%);
  color: #fff;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (min-width: 720px) {
  .hero-actions {
    grid-template-columns: 1fr 1fr;
  }

  .hero-btn {
    min-height: 220px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .home-intro,
  .hero-btn {
    animation: none;
  }
}
</style>
