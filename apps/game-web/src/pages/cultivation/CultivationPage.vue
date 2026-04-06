<template>
  <section class="cultivate">
    <div class="panel">
      <h1>修行池</h1>
      <p>点燃灵火，汲取天地之力，提升战灵和魔魂底蕴。</p>
      <div class="meter">
        <div class="meter-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="meter-label">神力进度 {{ progress }}%</p>
      <p class="meter-label">当前状态 {{ cultivation.state }}</p>
      <div class="actions">
        <button class="primary" type="button" @click="beginCultivation">立刻修行</button>
        <button class="ghost" type="button" :disabled="!canClaimReward" @click="claimReward">领取灵力</button>
      </div>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="grid">
      <article>
        <h2>灵力储备</h2>
        <p>当前：{{ cultivation.spirit_power }}</p>
        <small>最多可储存 300 点</small>
      </article>
      <article>
        <h2>炼妖记录</h2>
        <p>{{ cultivation.state }}</p>
        <small>预计可领取：{{ cultivation.claimable_at || "未开始修行" }}</small>
        <small v-if="cultivation.pet_growth?.exp">幻兽经验 +{{ cultivation.pet_growth.exp }}</small>
      </article>
      <article>
        <h2>化仙池</h2>
        <p>钱包灵力 {{ walletPower }}</p>
        <small>已解锁 2 个法阵</small>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { claimCultivation, startCultivation, type CultivationStatus } from "@/api/modules/dungeon"
import { getHomeOverview } from "@/api/modules/home"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const cultivation = ref<CultivationStatus>({
  player_id: 0,
  spirit_power: 0,
  state: "idle",
  pet_growth: {
    exp: 0,
    team_total_power: 0,
  },
})
const walletPower = ref(0)
const errorMessage = ref("")
const CULTIVATION_NOT_FOUND_CODE = 4041
const CULTIVATION_NOT_READY_CODE = 4091

const progress = computed(() => Math.min(100, cultivation.value.spirit_power * 10))
const canClaimReward = computed(() => {
  if (cultivation.value.claimable_at) {
    const claimableAt = new Date(cultivation.value.claimable_at).getTime()
    return Number.isFinite(claimableAt) && Date.now() >= claimableAt
  }

  return cultivation.value.state !== "idle" && cultivation.value.spirit_power > 0
})

function normalizeCultivation(nextStatus: Partial<CultivationStatus>): CultivationStatus {
  const defaultPetGrowth = {
    exp: 0,
    team_total_power: 0,
  }
  const baseStatus = {
    player_id: 0,
    spirit_power: 0,
    state: "idle",
  }

  return {
    ...baseStatus,
    ...nextStatus,
    pet_growth: {
      ...defaultPetGrowth,
      ...nextStatus.pet_growth,
    },
  }
}

function resetCultivation(playerId: number) {
  cultivation.value = normalizeCultivation({
    player_id: playerId,
    spirit_power: 0,
    state: "idle",
  })
}

function isCultivationMissingError(error: unknown) {
  return (
    error instanceof APIError &&
    error.status === 404 &&
    error.code === CULTIVATION_NOT_FOUND_CODE
  )
}

function isCultivationNotReadyError(error: unknown) {
  return (
    error instanceof APIError &&
    error.status === 409 &&
    error.code === CULTIVATION_NOT_READY_CODE
  )
}

async function loadOverview() {
  if (!sessionStore.playerId) {
    errorMessage.value = "当前未登录，无法加载修行状态。"
    return
  }

  try {
    const overview = await getHomeOverview(sessionStore.playerId)
    cultivation.value = normalizeCultivation({
      player_id: overview.player_id,
      spirit_power: overview.modules.cultivation.spirit_power,
      state: overview.modules.cultivation.state,
      claimable_at: overview.modules.cultivation.claimable_at,
    })
    walletPower.value = overview.wallet.spirit_power
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "修行状态加载失败。"
  }
}

async function beginCultivation() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法开始修行。"
    return
  }

  try {
    cultivation.value = normalizeCultivation(await startCultivation(playerId))
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "开始修行失败。"
  }
}

async function claimReward() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法领取灵力。"
    return
  }
  if (!canClaimReward.value) {
    return
  }

  try {
    cultivation.value = normalizeCultivation(await claimCultivation(playerId))
    walletPower.value += cultivation.value.spirit_power
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    if (isCultivationMissingError(error)) {
      resetCultivation(playerId)
      errorMessage.value = "尚未开始修行，请先开始修行。"
      return
    }
    if (isCultivationNotReadyError(error)) {
      errorMessage.value = "修行尚未完成，请稍后领取。"
      return
    }
    errorMessage.value = error instanceof APIError ? error.message : "领取灵力失败。"
  }
}

onMounted(async () => {
  await loadOverview()
})
</script>

<style scoped>
.cultivate {
  min-height: 100vh;
  padding: 36px;
  background: radial-gradient(circle at top, #2c284b, #0a0615 80%);
  color: #fbf8ff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.status-text {
  margin: 0;
  color: #ffcbdf;
}
.panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.03);
  box-shadow: 0 20px 50px rgba(5, 5, 20, 0.45);
}
.meter {
  width: 100%;
  height: 6px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  margin: 18px 0;
  overflow: hidden;
}
.meter-fill {
  height: 100%;
  background: linear-gradient(90deg, #f4d35d, #f9576f);
  transition: width 0.3s;
}
.meter-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
}
.actions {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}
.primary {
  border: none;
  border-radius: 14px;
  padding: 12px 18px;
  background: linear-gradient(135deg, #f9576f, #ffb347);
  color: #190d00;
  font-weight: 600;
  cursor: pointer;
}
.ghost {
  border-radius: 14px;
  padding: 12px 18px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.04);
  color: #fff;
  cursor: pointer;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 18px;
}
.grid article {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.03), rgba(255, 255, 255, 0.06));
}
.grid h2 {
  margin-bottom: 8px;
}
.grid small {
  color: rgba(255, 255, 255, 0.6);
}
</style>
