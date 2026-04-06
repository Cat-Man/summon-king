<template>
  <section class="home-page">
    <div class="hero-panel">
      <p class="hero-tag">召唤之王</p>
      <h2>首页总览</h2>
      <p class="hero-copy">
        欢迎回来，{{ displayNickname }}。当前首页已经接入真实 overview 接口，作为地图、副本、修行、竞技场、排行榜以及塔的聚合入口。
      </p>
      <div class="hero-actions">
        <span>玩家 ID：{{ sessionStore.playerId ?? "-" }}</span>
        <span>后端接口：/api/v1/home/overview</span>
      </div>
    </div>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div v-if="nextActionCard" class="next-action-card">
      <div>
        <p>下一步推荐</p>
        <h3>{{ nextActionCard.title }}</h3>
        <p class="next-action-copy">{{ nextActionCard.description }}</p>
      </div>
      <RouterLink class="next-action-cta" :to="nextActionCard.route">
        {{ nextActionCard.cta }}
      </RouterLink>
    </div>

    <div class="overview-grid">
      <article v-for="card in moduleCards" :key="card.title" class="overview-card">
        <p>{{ card.eyebrow }}</p>
        <h3>{{ card.title }}</h3>
        <strong>{{ card.value }}</strong>
        <span>{{ card.description }}</span>
        <RouterLink class="entry-link" :to="card.route">{{ card.cta }}</RouterLink>
      </article>
    </div>

    <article v-if="growthSummary" class="growth-summary-card">
      <header>
        <div>
          <p>养成入口</p>
          <h3>成长收益摘要</h3>
        </div>
        <strong>养成总加成 +{{ growthSummary.totalBonus }}</strong>
      </header>
      <p class="growth-summary-copy">
        当前战骨、战灵、魔魂加成已经统一沉淀到共享战斗队，首页可以直接跳回对应养成入口继续推进。
      </p>
      <div class="growth-summary-metrics">
        <span>战骨 +{{ growthSummary.boneBonus }}</span>
        <span>战灵 +{{ growthSummary.spiritBonus }}</span>
        <span>魔魂 +{{ growthSummary.soulBonus }}</span>
      </div>
      <div class="growth-summary-links">
        <RouterLink class="entry-link" to="/growth/bone">前往战骨</RouterLink>
        <RouterLink class="entry-link" to="/growth/spirit">前往战灵</RouterLink>
        <RouterLink class="entry-link" to="/growth/soul">前往魔魂</RouterLink>
      </div>
    </article>

    <div class="tower-grid" v-if="towerCards.length">
      <article v-for="tower in towerCards" :key="tower.title" class="tower-card">
        <header>
          <h3>{{ tower.title }}</h3>
          <span>{{ tower.label }}</span>
        </header>
        <p>{{ tower.description }}</p>
        <div class="tower-meta">
          <span>当前层级：第 {{ tower.currentFloor }} 层</span>
          <span>剩余挑战：{{ tower.remainingChallenges }}/5</span>
        </div>
        <RouterLink class="entry-link" :to="tower.route">{{ tower.cta }}</RouterLink>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import { getHomeOverview, type HomeOverview } from "@/api/modules/home"
import { getPetCollection, type PetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const overview = ref<HomeOverview | null>(null)
const petCollection = ref<PetCollection | null>(null)
const errorMessage = ref("")

const displayNickname = computed(() => overview.value?.nickname || sessionStore.nickname || "未登录玩家")

const nextActionCard = computed(() => {
  if (!overview.value?.next_action) {
    return null
  }
  return overview.value.next_action
})

const moduleCards = computed(() => {
  if (!overview.value) {
    return [
      {
        eyebrow: "今日概览",
        title: "首页状态",
        value: "加载中",
        description: "正在同步地图、副本与修行的最新状态。",
        route: "/home",
        cta: "刷新中",
      },
    ]
  }

  const cultivationModule = overview.value.modules.cultivation
  const cultivationCard = cultivationModule.claimable
    ? {
        eyebrow: "修行链路",
        title: "修行进度",
        value: "可领取",
        description: `已积累 ${cultivationModule.spirit_power} 灵力，当前可以直接领取收益。`,
        route: "/cultivation",
        cta: "立即领取",
      }
    : cultivationModule.state === "idle"
      ? {
          eyebrow: "修行链路",
          title: "修行进度",
          value: "未开始",
          description: `钱包灵力 ${overview.value.wallet.spirit_power}，新一轮修行还未启动。`,
          route: "/cultivation",
          cta: "开始修行",
        }
      : {
          eyebrow: "修行链路",
          title: "修行进度",
          value: cultivationModule.state,
          description: `当前灵力 ${cultivationModule.spirit_power}，钱包灵力 ${overview.value.wallet.spirit_power}。`,
          route: "/cultivation",
          cta: "前往修行",
        }

  return [
    {
      eyebrow: "地图入口",
      title: "世界探索",
      value: overview.value.modules.map_label,
      description: `当前开放 ${overview.value.modules.map_city_count} 座城市，下一步可直接进入世界地图。`,
      route: "/map",
      cta: "前往地图",
    },
    {
      eyebrow: "副本链路",
      title: "地下城状态",
      value: `第 ${overview.value.modules.dungeon.current_floor} 层`,
      description: `状态 ${overview.value.modules.dungeon.status}，剩余骰子 ${overview.value.modules.dungeon.remain_dice}。`,
      route: "/dungeon",
      cta: "进入副本",
    },
    cultivationCard,
    {
      eyebrow: "阵容概览",
      title: "幻兽阵容",
      value: `${overview.value.modules.pet.total_power} 战力`,
      description: `主战 ${overview.value.modules.pet.starter_pet_name || "未上阵幻兽"}，已上阵 ${overview.value.modules.pet.active_count} 只。`,
      route: "/pet",
      cta: "继续养成",
    },
    {
      eyebrow: "竞技场",
      title: "斗法连胜",
      value: `当前连胜 ${overview.value.modules.arena.current_streak} 场`,
      description: overview.value.modules.arena.last_win ? "上场斗法取胜，适合继续冲榜。" : "上场斗法失利，先挑选低战对手回稳。",
      route: "/arena",
      cta: "前往竞技场",
    },
    {
      eyebrow: "排行榜",
      title: "巅峰荣耀榜",
      value: `第 ${overview.value.modules.ranking.self_rank} 名`,
      description: `当前榜单分数 ${overview.value.modules.ranking.self_score}，继续斗法可进一步抬升。`,
      route: "/ranking",
      cta: "查看排行榜",
    },
  ]
})

const growthSummary = computed(() => {
  const activeTeam = petCollection.value?.active_team ?? []
  if (activeTeam.length === 0) {
    return null
  }

  const boneBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.bone ?? 0), 0)
  const spiritBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.spirit ?? 0), 0)
  const soulBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.soul ?? 0), 0)

  return {
    boneBonus,
    spiritBonus,
    soulBonus,
    totalBonus: boneBonus + spiritBonus + soulBonus,
  }
})

type TowerCard = {
  route: string
  title: string
  label: string
  currentFloor: number
  remainingChallenges: number
  description: string
  cta: string
}

const towerCards = computed<TowerCard[]>(() => {
  const tower = overview.value?.modules.tower
  if (!tower) {
    return []
  }
  return [
    {
      title: "通天塔试炼",
      label: "pagoda",
      currentFloor: tower.pagoda.current_floor,
      remainingChallenges: tower.pagoda.remaining_challenges,
      description: `当前奖励预览：${tower.pagoda.reward_preview}，累计挑战带来成长收益。`,
      route: "/tower/pagoda",
      cta: "继续挑战",
    },
    {
      title: "战灵塔试炼",
      label: "spirit",
      currentFloor: tower.spirit.current_floor,
      remainingChallenges: tower.spirit.remaining_challenges,
      description: `当前奖励预览：${tower.spirit.reward_preview}，灵力与魔魂成长同步。`,
      route: "/tower/spirit",
      cta: "前往战灵塔",
    },
  ]
})

async function loadOverview() {
  if (!sessionStore.playerId) {
    errorMessage.value = "当前未登录，无法加载首页总览。"
    return
  }

  try {
    overview.value = await getHomeOverview(sessionStore.playerId)
    errorMessage.value = ""
  } catch (error) {
    if (error instanceof APIError) {
      errorMessage.value = error.message
      return
    }
    errorMessage.value = "首页总览加载失败，请稍后重试。"
  }
}

async function loadPetSummary() {
  if (!sessionStore.playerId) {
    return
  }

  try {
    petCollection.value = await getPetCollection(sessionStore.playerId)
  } catch {}
}

async function loadHomeData() {
  await Promise.all([loadOverview(), loadPetSummary()])
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadHomeData()
  },
  { flush: "sync" },
)

onMounted(async () => {
  await loadHomeData()
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

.next-action-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-radius: 24px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(15, 18, 28, 0.85);
  box-shadow: 0 18px 50px rgba(0, 0, 0, 0.35);
}

.next-action-card p {
  margin: 0;
  color: rgba(247, 189, 120, 0.78);
  letter-spacing: 0.2em;
  font-size: 12px;
  text-transform: uppercase;
}

.next-action-card h3 {
  margin: 6px 0;
  font-size: 32px;
}

.next-action-copy {
  margin: 0;
  color: rgba(247, 239, 225, 0.72);
}

.next-action-cta {
  padding: 12px 24px;
  border-radius: 999px;
  background: rgba(247, 189, 120, 0.12);
  color: #f7d8a8;
  border: 1px solid rgba(247, 189, 120, 0.4);
  text-decoration: none;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}

.growth-summary-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(147, 197, 253, 0.16);
  background:
    radial-gradient(circle at top right, rgba(125, 211, 252, 0.14), transparent 25%),
    rgba(15, 20, 31, 0.92);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.24);
}

.growth-summary-card header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 16px;
}

.growth-summary-card header p,
.growth-summary-card header h3,
.growth-summary-copy {
  margin: 0;
}

.growth-summary-card header p {
  color: rgba(147, 197, 253, 0.82);
  letter-spacing: 0.18em;
  font-size: 12px;
  text-transform: uppercase;
}

.growth-summary-card header h3 {
  margin-top: 6px;
}

.growth-summary-card header strong {
  color: #7dd3fc;
  font-size: 20px;
}

.growth-summary-copy {
  margin-top: 12px;
  color: rgba(247, 239, 225, 0.72);
  line-height: 1.7;
}

.growth-summary-metrics,
.growth-summary-links {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
}

.growth-summary-metrics span {
  padding: 10px 12px;
  border-radius: 14px;
  background: rgba(247, 239, 225, 0.06);
  color: rgba(247, 239, 225, 0.86);
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
  text-decoration: none;
}

.tower-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 18px;
}

.tower-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(255, 147, 97, 0.3);
  background: linear-gradient(180deg, rgba(15, 8, 3, 0.85), rgba(30, 7, 7, 0.95));
  box-shadow: 0 20px 45px rgba(0, 0, 0, 0.3);
}

.tower-card header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.tower-card h3 {
  margin: 0;
  font-size: 24px;
}

.tower-card span {
  font-size: 12px;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.8);
}

.tower-card p {
  margin: 16px 0 12px;
  color: rgba(247, 239, 225, 0.7);
}

.tower-meta {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  margin-bottom: 12px;
}

@media (max-width: 720px) {
  .hero-panel {
    padding: 24px;
  }

  .next-action-card {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
