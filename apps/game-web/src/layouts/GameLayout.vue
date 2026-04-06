<template>
  <div class="game-shell">
    <header class="shell-header">
      <div class="brand-block">
        <p class="brand-mark">ZHZW</p>
        <div>
          <h1>召唤之王</h1>
          <p>正式应用壳</p>
        </div>
      </div>

      <div class="shell-nav-wrap">
        <div class="session-panel">
          <div class="session-copy">
            <strong>{{ sessionStore.nickname || "未登录访客" }}</strong>
            <span>玩家 ID：{{ sessionStore.playerId ?? "-" }}</span>
          </div>
          <button data-testid="logout-button" class="logout-button" type="button" @click="logout">
            退出登录
          </button>
        </div>
        <nav class="shell-nav" aria-label="主导航">
          <RouterLink v-for="item in primaryNav" :key="item.to" :to="item.to">{{ item.label }}</RouterLink>
        </nav>
        <nav class="shell-subnav" aria-label="功能导航">
          <RouterLink v-for="item in featureNav" :key="item.to" :to="item.to">{{ item.label }}</RouterLink>
        </nav>
      </div>
    </header>

    <main class="shell-main">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router"

import { useSessionStore } from "@/stores/session"

const router = useRouter()
const sessionStore = useSessionStore()

const primaryNav = [
  { to: "/home", label: "首页" },
  { to: "/map", label: "地图" },
  { to: "/dungeon", label: "副本" },
  { to: "/cultivation", label: "修行" },
]

const featureNav = [
  { to: "/pet", label: "幻兽" },
  { to: "/growth/spirit", label: "战灵" },
  { to: "/growth/bone", label: "战骨" },
  { to: "/growth/soul", label: "魔魂" },
  { to: "/growth/manor", label: "庄园" },
  { to: "/tower/pagoda", label: "通天塔" },
  { to: "/tower/spirit", label: "战灵塔" },
  { to: "/arena", label: "竞技场" },
  { to: "/ranking", label: "排行榜" },
]

async function logout() {
  sessionStore.clearSession()
  await router.push("/login")
}
</script>

<style scoped>
:global(:root) {
  color-scheme: dark;
  font-family: "Noto Serif SC", "Hiragino Mincho ProN", "Songti SC", serif;
  background:
    radial-gradient(circle at top left, rgba(247, 189, 120, 0.18), transparent 28%),
    radial-gradient(circle at top right, rgba(65, 122, 220, 0.22), transparent 34%),
    linear-gradient(180deg, #120e1d, #07090f 60%);
}

:global(body) {
  margin: 0;
  min-width: 320px;
  background: transparent;
  color: #f7efe1;
}

:global(*) {
  box-sizing: border-box;
}

:global(a) {
  color: inherit;
  text-decoration: none;
}

.game-shell {
  min-height: 100vh;
  padding: 24px;
}

.shell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin: 0 auto 24px;
  max-width: 1180px;
  padding: 18px 22px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  border-radius: 24px;
  background: rgba(11, 13, 20, 0.68);
  backdrop-filter: blur(16px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
}

.brand-block {
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand-mark {
  margin: 0;
  padding: 10px 12px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(244, 164, 96, 0.24), rgba(217, 119, 6, 0.48));
  font-family: "Avenir Next Condensed", "DIN Alternate", sans-serif;
  font-size: 13px;
  letter-spacing: 0.32em;
}

.brand-block h1,
.brand-block p {
  margin: 0;
}

.brand-block h1 {
  font-size: clamp(24px, 4vw, 34px);
}

.brand-block p {
  margin-top: 4px;
  color: rgba(247, 239, 225, 0.68);
  font-size: 14px;
}

.shell-nav {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.session-panel {
  display: flex;
  align-items: center;
  gap: 14px;
}

.session-copy {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.session-copy strong,
.session-copy span {
  margin: 0;
}

.session-copy strong {
  font-size: 14px;
  color: #fff8f0;
}

.session-copy span {
  color: rgba(247, 239, 225, 0.64);
  font-size: 12px;
}

.logout-button {
  padding: 9px 14px;
  border: 1px solid rgba(247, 239, 225, 0.14);
  border-radius: 12px;
  background: rgba(247, 239, 225, 0.04);
  color: rgba(247, 239, 225, 0.86);
  font-size: 13px;
  cursor: pointer;
  transition: background 160ms ease, border-color 160ms ease, transform 160ms ease;
}

.logout-button:hover {
  transform: translateY(-1px);
  border-color: rgba(247, 189, 120, 0.4);
  background: rgba(247, 189, 120, 0.08);
}

.shell-nav a {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid rgba(247, 239, 225, 0.14);
  color: rgba(247, 239, 225, 0.78);
  transition: transform 160ms ease, border-color 160ms ease, color 160ms ease;
}

.shell-nav a.router-link-active {
  color: #fff8f0;
  border-color: rgba(247, 189, 120, 0.55);
  background: rgba(247, 189, 120, 0.12);
}

.shell-nav a:hover {
  transform: translateY(-1px);
  border-color: rgba(247, 239, 225, 0.35);
}

.shell-nav-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 10px;
}

.shell-subnav {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
}

.shell-subnav a {
  padding: 8px 12px;
  border-radius: 12px;
  background: rgba(247, 239, 225, 0.04);
  color: rgba(247, 239, 225, 0.68);
  font-size: 13px;
  transition: background 160ms ease, color 160ms ease;
}

.shell-subnav a.router-link-active {
  background: rgba(104, 171, 255, 0.16);
  color: #e8f2ff;
}

.shell-subnav a:hover {
  background: rgba(247, 239, 225, 0.1);
  color: #fff8f0;
}

.shell-main {
  margin: 0 auto;
  max-width: 1180px;
}

@media (max-width: 720px) {
  .game-shell {
    padding: 16px;
  }

  .shell-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .shell-nav-wrap {
    width: 100%;
    align-items: flex-start;
  }

  .session-panel,
  .session-copy {
    align-items: flex-start;
  }

  .shell-subnav {
    justify-content: flex-start;
  }
}
</style>
