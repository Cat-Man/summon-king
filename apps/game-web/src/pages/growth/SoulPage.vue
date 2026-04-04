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
    <p class="note">魔魂线当前先接真实数值，后续再扩猎魂与掉落玩法。</p>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { getSoulState, type SoulState } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const soul = ref<SoulState>({
  name: "",
  power: 0,
})
const errorMessage = ref("")

onMounted(async () => {
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
  max-width: 420px;
  line-height: 1.6;
}
</style>
