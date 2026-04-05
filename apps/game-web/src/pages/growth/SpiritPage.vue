<template>
  <section class="growth-page growth-page--spirit">
    <header>
      <p class="tag">战灵</p>
      <h1>{{ wallet.player_id ? "Spirit Forge" : "战灵熔炉" }}</h1>
      <p>灵力 {{ wallet.spirit_power }} | 洗炼剩余次数 {{ wallet.spirit_free_wash }}</p>
    </header>
    <div class="banner">
      <span>免费洗炼优先消耗</span>
      <span>战骨 {{ wallet.bone_level }} · 魔魂 {{ wallet.soul_pieces }}</span>
    </div>
    <div class="panel panel--power">
      <p>阵容战力 {{ teamPower }}</p>
      <span>当前战灵成长已计入共享战斗队。</span>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <button class="primary" type="button" @click="washNow">立刻洗炼</button>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue"

import { APIError } from "@/api/http"
import { getGrowthWallet, washSpirit, type GrowthWallet } from "@/api/modules/growth"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const wallet = ref<GrowthWallet>({
  player_id: 0,
  spirit_power: 0,
  spirit_free_wash: 0,
  bone_level: 0,
  soul_pieces: 0,
  manor_plots: 0,
})
const teamPower = ref(0)
const errorMessage = ref("")

async function loadWallet() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载战灵。"
    return
  }

  try {
    wallet.value = await getGrowthWallet(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战灵状态加载失败。"
  }
}

async function loadTeamPower() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    return
  }

  try {
    const collection = await getPetCollection(playerId)
    teamPower.value = collection.total_power
  } catch {}
}

async function washNow() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法洗炼战灵。"
    return
  }

  try {
    await washSpirit(playerId)
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战灵洗炼失败。"
  }
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await Promise.all([loadWallet(), loadTeamPower()])
  },
)

onMounted(async () => {
  await Promise.all([loadWallet(), loadTeamPower()])
})
</script>

<style scoped>
.growth-page--spirit {
  padding: 32px;
  border-radius: 24px;
  background: radial-gradient(circle at top right, #f7ba38, #dd2e6d 60%);
  color: #171717;
  min-height: 52vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.tag {
  font-size: 12px;
  letter-spacing: 0.25em;
  text-transform: uppercase;
}
.banner {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.4);
  display: flex;
  justify-content: space-between;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.18);
}
.panel--power p,
.panel--power span {
  margin: 0;
}
.panel--power span {
  display: block;
  margin-top: 6px;
  color: rgba(23, 23, 23, 0.72);
  font-size: 14px;
}
.status-text {
  margin: 0;
  color: #651327;
}
.primary {
  padding: 14px 24px;
  border: none;
  border-radius: 12px;
  background: #171717;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
}
</style>
