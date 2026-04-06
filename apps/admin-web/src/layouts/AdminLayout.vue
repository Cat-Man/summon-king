<script setup lang="ts">
const navigationItems = [
  {
    label: "工作台",
    description: "查看后台建设边界与主入口",
    to: { name: "dashboard" as const },
  },
  {
    label: "配置中心",
    description: "数值与玩法配置域",
    to: { name: "config" as const },
  },
  {
    label: "GM 操作台",
    description: "人工运营与异常处置域",
    to: { name: "gm" as const },
  },
]
</script>

<template>
  <div class="admin-shell">
    <aside class="shell-sidebar">
      <div class="brand-panel">
        <p class="eyebrow">Summon King Console</p>
        <h1>召唤之王</h1>
        <p>运营后台</p>
      </div>
      <nav class="shell-nav" aria-label="后台导航">
        <RouterLink
          v-for="item in navigationItems"
          :key="item.label"
          :to="item.to"
          class="nav-link"
          active-class="nav-link-active"
        >
          <strong>{{ item.label }}</strong>
          <span>{{ item.description }}</span>
        </RouterLink>
      </nav>
    </aside>

    <main class="shell-main">
      <header class="topbar">
        <div>
          <p class="topbar-label">当前阶段</p>
          <h2>后台壳初始化</h2>
        </div>
        <div class="status-chip">仅开放工作台 / 配置 / GM 占位</div>
      </header>

      <section class="content-panel">
        <RouterView />
      </section>
    </main>
  </div>
</template>

<style scoped>
.admin-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  background:
    radial-gradient(circle at top left, rgba(205, 143, 60, 0.18), transparent 28%),
    linear-gradient(135deg, #f6efe5 0%, #ede2d2 42%, #d7c4aa 100%);
  color: #211a14;
  font-family: "Avenir Next", "PingFang SC", "Hiragino Sans GB", sans-serif;
}

.shell-sidebar {
  padding: 28px 20px;
  background: rgba(30, 21, 15, 0.9);
  color: #f3ecdf;
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.brand-panel {
  padding: 20px;
  border-radius: 24px;
  background: linear-gradient(180deg, rgba(196, 148, 88, 0.28), rgba(110, 71, 31, 0.3));
  border: 1px solid rgba(233, 206, 164, 0.24);
}

.brand-panel h1,
.brand-panel p,
.eyebrow,
.topbar h2,
.topbar-label {
  margin: 0;
}

.eyebrow,
.topbar-label {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 12px;
  opacity: 0.72;
}

.brand-panel h1 {
  margin-top: 10px;
  font-size: 34px;
  line-height: 1.1;
}

.brand-panel p:last-child {
  margin-top: 8px;
  color: rgba(243, 236, 223, 0.78);
}

.shell-nav {
  display: grid;
  gap: 12px;
}

.nav-link {
  display: grid;
  gap: 6px;
  padding: 16px 18px;
  border-radius: 18px;
  text-decoration: none;
  color: inherit;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid transparent;
  transition:
    transform 140ms ease,
    border-color 140ms ease,
    background-color 140ms ease;
}

.nav-link:hover,
.nav-link-active {
  transform: translateX(4px);
  border-color: rgba(240, 214, 175, 0.32);
  background: rgba(244, 226, 197, 0.12);
}

.nav-link strong {
  font-size: 17px;
}

.nav-link span {
  font-size: 13px;
  color: rgba(243, 236, 223, 0.72);
}

.shell-main {
  padding: 28px;
  display: grid;
  gap: 20px;
}

.topbar {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 16px;
}

.topbar h2 {
  margin-top: 6px;
  font-size: 34px;
}

.status-chip {
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(50, 37, 24, 0.08);
  color: #5f4a34;
  font-size: 14px;
}

.content-panel {
  min-height: 0;
}

@media (max-width: 960px) {
  .admin-shell {
    grid-template-columns: 1fr;
  }

  .shell-sidebar {
    gap: 18px;
  }

  .shell-main {
    padding: 20px;
  }

  .topbar {
    flex-direction: column;
    align-items: start;
  }
}
</style>
