<template>
  <section class="growth-page growth-page--soul">
    <header>
      <p class="tag">魔魂</p>
      <h1>{{ soul.name || "Soul Catcher" }}</h1>
      <p>当前魂力 {{ soul.power }}</p>
    </header>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="grid">
      <article>
        <strong>魂力值</strong>
        <span>{{ soul.power }}</span>
      </article>
      <article>
        <strong>状态</strong>
        <span>{{ soul.power > 0 ? "已凝聚" : "待收集" }}</span>
      </article>
    </div>
    <button class="primary" type="button" @click="upgradeNow">升级魔魂</button>
    <div class="panel">
      <p class="note">魔魂线当前先接真实数值，后续再扩猎魂与掉落玩法。</p>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue"

import { APIError } from "@/api/http"
import { getSoulState, type SoulState, upgradeSoul } from "@/api/modules/growth"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const soul = ref<SoulState>({
  name: "",
  power: 0,
})
const errorMessage = ref("")

async function loadSoul() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载魔魂。"
    return
  }

  try {
    soul.value = await getSoulState(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "魔魂状态加载失败。"
  }
}

async function upgradeNow() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法升级魔魂。"
    return
  }

  try {
    await upgradeSoul(playerId)
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "魔魂升级失败。"
  }
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadSoul()
  },
)

onMounted(async () => {
  await loadSoul()
})
</script>

<style scoped>
.growth-page--soul {
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(180deg, #0c2c4a, #07233d);
  color: #e3f2ff;
  min-height: 52vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.grid article {
  padding: 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
}
.primary {
  align-self: flex-start;
  padding: 12px 18px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #7ec9ff, #93a8ff);
  color: #05253c;
  font-weight: 700;
  cursor: pointer;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.06);
}
.tag {
  font-size: 12px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
}
.status-text {
  margin: 0;
  color: #ffdce8;
}
.note {
  margin: 0;
  max-width: 420px;
  line-height: 1.6;
}
</style>
