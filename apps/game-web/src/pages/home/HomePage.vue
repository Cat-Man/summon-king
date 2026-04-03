<template>
  <section class="home-page">
    <div class="hero-panel">
      <p class="hero-tag">召唤之王</p>
      <h2>首页总览</h2>
      <p class="hero-copy">
        欢迎回来，{{ displayNickname }}。当前首页已经接入真实 overview 接口，作为地图、副本与修行的聚合入口。
      </p>
      <div class="hero-actions">
        <span>玩家 ID：{{ sessionStore.playerId ?? "-" }}</span>
        <span>后端接口：/api/v1/home/overview</span>
      </div>
    </div>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div class="overview-grid">
      <article v-for="card in cards" :key="card.title" class="overview-card">
        <p>{{ card.eyebrow }}</p>
        <h3>{{ card.title }}</h3>
        <strong>{{ card.value }}</strong>
        <span>{{ card.description }}</span>
        <a class="entry-link" :href="card.href">{{ card.cta }}</a>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getHomeOverview, type HomeOverview } from "@/api/modules/home"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const overview = ref<HomeOverview | null>(null)
const errorMessage = ref("")

const displayNickname = computed(() => overview.value?.nickname || sessionStore.nickname || "未登录玩家")

const cards = computed(() => {
  if (!overview.value) {
    return [
      {
        eyebrow: "今日概览",
        title: "首页状态",
        value: "加载中",
        description: "正在同步地图、副本与修行的最新状态。",
        href: "/home",
        cta: "刷新中",
      },
    ]
  }

  return [
    {
      eyebrow: "地图入口",
      title: "世界探索",
      value: overview.value.modules.map_label,
      description: `当前开放 ${overview.value.modules.map_city_count} 座城市，下一步可直接进入世界地图。`,
      href: "/map",
      cta: "前往地图",
    },
    {
      eyebrow: "副本链路",
      title: "地下城状态",
      value: `第 ${overview.value.modules.dungeon.current_floor} 层`,
      description: `状态 ${overview.value.modules.dungeon.status}，剩余骰子 ${overview.value.modules.dungeon.remain_dice}。`,
      href: "/dungeon",
      cta: "进入副本",
    },
    {
      eyebrow: "修行链路",
      title: "修行进度",
      value: overview.value.modules.cultivation.state,
      description: `当前灵力 ${overview.value.modules.cultivation.spirit_power}，钱包灵力 ${overview.value.wallet.spirit_power}。`,
      href: "/cultivation",
      cta: "前往修行",
    },
  ]
})

onMounted(async () => {
  if (!sessionStore.playerId) {
    errorMessage.value = "当前未登录，无法加载首页总览。"
    return
  }

  try {
    overview.value = await getHomeOverview(sessionStore.playerId)
  } catch (error) {
    if (error instanceof APIError) {
      errorMessage.value = error.message
      return
    }
    errorMessage.value = "首页总览加载失败，请稍后重试。"
  }
})
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.status-text {
  margin: 0;
  color: #ffb6a2;
}

.hero-panel {
  padding: 32px;
  border-radius: 30px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background:
    radial-gradient(circle at top right, rgba(237, 164, 93, 0.16), transparent 24%),
    linear-gradient(180deg, rgba(22, 25, 37, 0.92), rgba(12, 15, 24, 0.9));
  box-shadow: 0 30px 70px rgba(0, 0, 0, 0.35);
}

.hero-tag {
  display: inline-flex;
  margin: 0 0 10px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(247, 189, 120, 0.16);
  color: #f7d8a8;
  letter-spacing: 0.24em;
  font-size: 12px;
  text-transform: uppercase;
}

.hero-panel h2,
.hero-copy {
  margin: 0;
}

.hero-panel h2 {
  font-size: clamp(32px, 5vw, 52px);
  line-height: 1.05;
}

.hero-copy {
  margin-top: 14px;
  max-width: 720px;
  color: rgba(247, 239, 225, 0.72);
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
}

.hero-actions span {
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(247, 239, 225, 0.06);
  color: rgba(247, 239, 225, 0.8);
  font-family: "Avenir Next", "PingFang SC", sans-serif;
  font-size: 14px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}

.overview-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(247, 239, 225, 0.1);
  background: rgba(12, 14, 22, 0.8);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.25);
}

.overview-card p,
.overview-card h3,
.overview-card strong,
.overview-card span {
  display: block;
}

.overview-card p {
  margin: 0 0 10px;
  color: rgba(247, 189, 120, 0.78);
  letter-spacing: 0.16em;
  font-size: 12px;
  text-transform: uppercase;
}

.overview-card h3 {
  margin: 0;
  font-size: 22px;
}

.overview-card strong {
  margin-top: 18px;
  font-size: 28px;
}

.overview-card span {
  margin-top: 12px;
  color: rgba(247, 239, 225, 0.68);
  line-height: 1.6;
}

.entry-link {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  margin-top: 18px;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgba(247, 189, 120, 0.12);
  color: #f7d8a8;
}

@media (max-width: 720px) {
  .hero-panel {
    padding: 24px;
  }
}
</style>
