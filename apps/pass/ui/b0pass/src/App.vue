<template>
  <div class="page-shell">
    <header class="page-header">
      <div class="logo-mark">TeamTransfer</div>
      <nav class="page-nav">
        <a href="#features">功能</a>
        <a href="#tech">技术</a>
        <a href="#login">登录</a>
      </nav>
    </header>

    <section class="hero-section">
      <p class="hero-subtitle">个人使用的文件传输与实时协作平台</p>
      <h1>支持大文件传输 · 实时消息 · 权限控制 · 稳定可控</h1>
      <p class="hero-description">
        通过切片上传、断点续传与 WebRTC 的混合能力，保障镜像、日志、视频等内部文件能够稳定通过网络传输。
      </p>
      <div class="hero-actions">
        <button class="primary-btn">上传文件</button>
        <button class="ghost-btn">登录系统</button>
      </div>
    </section>

    <section id="features" class="section-grid">
      <h2>核心能力</h2>
      <div class="feature-grid">
        <article>
          <h3>大文件传输</h3>
          <p>支持分片上传与断点续传，适合镜像、日志、视频等多种大型内部资产。</p>
        </article>
        <article>
          <h3>实时消息</h3>
          <p>传输状态实时同步，所有操作节点可立即看到进度与结果。</p>
        </article>
        <article>
          <h3>权限与审计</h3>
          <p>可控的用户体系、可追溯的操作记录，保障文件安全。</p>
        </article>
      </div>
    </section>

    <section id="tech" class="section-grid section-alt">
      <h2>技术实现</h2>
      <div class="tech-grid">
        <div>
          <p>后端：Golang，保障高并发与稳定性</p>
          <p>文件上传：HTTP 分片 + 可扩展 WebRTC</p>
          <p>实时通信：WebSocket 连接池</p>
        </div>
        <div>
          <p>前端：Vue 3 + Composition API</p>
          <p>状态管理：Pinia</p>
          <p>UI：自定义组件，兼容 Tailwind / Element Plus 设计语言</p>
        </div>
      </div>
    </section>

    <section id="login" class="footer-section">
      <p>Internal Use Only · TeamTransfer</p>
    </section>
  </div>
</template>

<style scoped>
.page-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #ffffff 0%, #f4f6fb 100%);
  color: #0f172a;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 3rem;
  background: #ffffff;
  box-shadow: 0 1px 10px rgba(15, 23, 42, 0.08);
  position: sticky;
  top: 0;
  z-index: 10;
}

.logo-mark {
  font-weight: 700;
  font-size: 1.25rem;
  letter-spacing: 0.05em;
}

.page-nav {
  display: flex;
  gap: 1.5rem;
  font-size: 0.95rem;
}

.page-nav a {
  color: inherit;
  text-decoration: none;
  transition: color 0.2s ease;
}

.page-nav a:hover {
  color: #2563eb;
}

.hero-section {
  padding: 4rem 3rem 3rem;
  text-align: center;
}

.hero-subtitle {
  letter-spacing: 0.3em;
  text-transform: uppercase;
  margin-bottom: 0.8rem;
  font-size: 0.85rem;
  color: #475569;
}

.hero-section h1 {
  font-size: 2.75rem;
  margin: 0;
  font-weight: 700;
  margin-bottom: 1rem;
}

.hero-description {
  max-width: 640px;
  margin: 0 auto 2rem;
  color: #475569;
  line-height: 1.8;
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.primary-btn,
.ghost-btn {
  border: none;
  border-radius: 999px;
  padding: 0.85rem 2.5rem;
  font-size: 0.95rem;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.primary-btn {
  background: #2563eb;
  color: #fff;
  box-shadow: 0 10px 25px rgba(37, 99, 235, 0.35);
}

.ghost-btn {
  background: #fff;
  border: 1px solid #cbd5f5;
  color: #1d1f3d;
}

.primary-btn:hover,
.ghost-btn:hover {
  transform: translateY(-2px);
}

.section-grid {
  padding: 3rem 3rem 2rem;
  max-width: 1100px;
  margin: 0 auto;
}

.section-grid h2 {
  font-size: 2rem;
  margin-bottom: 1.5rem;
  text-align: center;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem;
}

.feature-grid article {
  background: #ffffff;
  border-radius: 1rem;
  padding: 1.6rem;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.08);
}

.feature-grid h3 {
  margin-top: 0;
  margin-bottom: 0.6rem;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
  padding: 1rem;
  background: #fff;
  border-radius: 1rem;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.06);
}

.tech-grid p {
  margin: 0.4rem 0;
  color: #334155;
}

.section-alt {
  background: #f9fbfd;
}

.footer-section {
  margin-top: auto;
  padding: 1.5rem;
  text-align: center;
  color: #475569;
  font-size: 0.9rem;
}

@media (max-width: 720px) {
  .page-header {
    flex-direction: column;
    gap: 0.5rem;
  }

  .section-grid {
    padding: 2rem 1.5rem;
  }
}
</style>
