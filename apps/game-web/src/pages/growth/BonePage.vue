<template>
  <section class="growth-page growth-page--bone">
    <header>
      <p class="tag">战骨</p>
      <h1>{{ bone.name || "Pyramid Bone" }}</h1>
      <p>铸魂之骨为成长基石，当前等级：{{ bone.level }} 级</p>
    </header>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="metrics">
      <article>
        <strong>当前等级</strong>
        <span>{{ bone.level }}</span>
      </article>
      <article>
        <strong>成长评价</strong>
        <span>{{ bone.level >= 3 ? "稳固" : "初成" }}</span>
      </article>
    </div>
    <div class="panel panel--power">
      <p>阵容战力 {{ teamPower }}</p>
      <span>战骨提升已实时计入共享战斗队。</span>
    </div>
    <button class="primary" type="button" @click="upgradeNow">升级战骨</button>
    <div class="panel">
      <p>当前战骨正在吸收雷火精华，后续再补升级资源与强化动作。</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue"

import { APIError } from "@/api/http"
import { getBoneState, type BoneState, upgradeBone } from "@/api/modules/growth"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const bone = ref<BoneState>({
  name: "",
  level: 0,
})
const teamPower = ref(0)
const errorMessage = ref("")

async function loadBone() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载战骨。"
    return
  }

  try {
    bone.value = await getBoneState(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战骨状态加载失败。"
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

async function upgradeNow() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法升级战骨。"
    return
  }

  try {
    bone.value = await upgradeBone(playerId)
    errorMessage.value = ""
    await loadTeamPower()
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "战骨升级失败。"
  }
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await Promise.all([loadBone(), loadTeamPower()])
  },
)

onMounted(async () => {
  await Promise.all([loadBone(), loadTeamPower()])
})
</script>

<style scoped>
.growth-page {
  min-height: 60vh;
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(160deg, #1f1b4b, #2b2a7b);
  color: #fff;
  box-shadow: 0 20px 40px rgba(22, 15, 70, 0.5);
}
.growth-page header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.tag {
  background: rgba(255, 255, 255, 0.15);
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 12px;
  letter-spacing: 0.2em;
  text-transform: uppercase;
}
.status-text {
  margin-top: 1rem;
  color: #ffd7e4;
}
.primary {
  align-self: flex-start;
  padding: 12px 18px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #8fb3ff, #b8a0ff);
  color: #130b30;
  font-weight: 700;
  cursor: pointer;
}
.metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin: 24px 0;
}
.metrics article {
  padding: 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.08);
}
.metrics strong {
  display: block;
  text-transform: uppercase;
  letter-spacing: 0.15em;
  font-size: 12px;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.06);
  line-height: 1.6;
}
.panel--power p,
.panel--power span {
  margin: 0;
}
.panel--power span {
  display: block;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 14px;
}
</style>
