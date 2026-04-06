<template>
  <section class="vip-page">
    <header class="hero-panel">
      <p class="hero-tag">VIP Daily</p>
      <h1>VIP 日常宝箱</h1>
      <p class="hero-copy">
        这一轮只实现最小 VIP 收益面板：查看当前等级、下一档权益与每日宝箱领取，不做充值和商城跳转。
      </p>
    </header>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div v-if="index" class="vip-grid">
      <article class="vip-summary">
        <p class="label">VIP{{ index.vip_level }}</p>
        <h2>{{ index.daily_claimed ? "今日已领取" : "领取每日宝箱" }}</h2>
        <p class="reward-copy">
          灵力 +{{ index.daily_reward.spirit_power }}
          <span v-if="index.daily_reward.soul_pieces"> · 魔魂碎片 +{{ index.daily_reward.soul_pieces }}</span>
        </p>
        <div class="actions">
          <button
            data-testid="vip-daily-claim"
            type="button"
            class="claim-button"
            :disabled="isSubmitting || index.daily_claimed"
            @click="handleClaim"
          >
            {{ index.daily_claimed ? "今日已领取" : "领取每日宝箱" }}
          </button>
          <span class="wallet-text">钱包灵力 {{ index.wallet_snapshot.spirit_power }}</span>
        </div>
      </article>

      <article class="vip-card">
        <h3>当前权益</h3>
        <ul class="benefit-list">
          <li v-for="benefit in index.current_benefits" :key="benefit">{{ benefit }}</li>
        </ul>
      </article>

      <article class="vip-card">
        <h3>下一档 VIP{{ index.next_level.level }}</h3>
        <p class="next-meta">累计消费 {{ index.next_level.required_gem_spent }} 元宝解锁</p>
        <ul class="benefit-list">
          <li v-for="benefit in index.next_level.benefits" :key="benefit">{{ benefit }}</li>
        </ul>
      </article>

      <article class="vip-card" v-if="lastReward">
        <h3>本次领取</h3>
        <p>灵力 +{{ lastReward.spirit_power }}</p>
        <p v-if="lastReward.soul_pieces">魔魂碎片 +{{ lastReward.soul_pieces }}</p>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { claimVIPDaily, getVIPIndex, type VIPDailyClaimResult, type VIPIndex, type VIPReward } from "@/api/modules/vip"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const index = ref<VIPIndex | null>(null)
const lastReward = ref<VIPReward | null>(null)
const errorMessage = ref("")
const isSubmitting = ref(false)

function normalizeReward(next?: Partial<VIPReward> | null): VIPReward {
  return {
    spirit_power: next?.spirit_power ?? 0,
    bone_level: next?.bone_level ?? 0,
    soul_pieces: next?.soul_pieces ?? 0,
  }
}

function normalizeIndex(next?: Partial<VIPIndex> | null, playerId = 0): VIPIndex {
  return {
    player_id: next?.player_id ?? playerId,
    vip_level: next?.vip_level ?? 0,
    total_gem_spent: next?.total_gem_spent ?? 0,
    daily_claimed: next?.daily_claimed ?? false,
    daily_reward: normalizeReward(next?.daily_reward),
    current_benefits: next?.current_benefits ?? [],
    next_level: {
      level: next?.next_level?.level ?? 1,
      required_gem_spent: next?.next_level?.required_gem_spent ?? 100,
      benefits: next?.next_level?.benefits ?? [],
    },
    wallet_snapshot: {
      spirit_power: next?.wallet_snapshot?.spirit_power ?? 0,
      bone_level: next?.wallet_snapshot?.bone_level ?? 0,
      soul_pieces: next?.wallet_snapshot?.soul_pieces ?? 0,
    },
  }
}

async function loadVIP() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载 VIP 状态。"
    return
  }

  try {
    index.value = normalizeIndex(await getVIPIndex(playerId), playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "VIP 状态加载失败。"
  }
}

function applyClaimResult(result: VIPDailyClaimResult) {
  index.value = normalizeIndex(result.index, result.index.player_id)
  lastReward.value = normalizeReward(result.reward_delta)
}

async function handleClaim() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法领取 VIP 宝箱。"
    return
  }
  if (!index.value || index.value.daily_claimed) {
    return
  }

  isSubmitting.value = true
  try {
    applyClaimResult(await claimVIPDaily(playerId))
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "VIP 宝箱领取失败。"
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  await loadVIP()
})
</script>

<style scoped>
.vip-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(247, 239, 225, 0.15);
  background:
    radial-gradient(circle at 18% 18%, rgba(251, 191, 36, 0.18), transparent 34%),
    radial-gradient(circle at 82% 14%, rgba(244, 114, 182, 0.16), transparent 30%),
    linear-gradient(180deg, rgba(34, 16, 28, 0.94), rgba(14, 10, 18, 0.96));
}

.hero-tag {
  margin: 0;
  font-size: 12px;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: rgba(251, 191, 36, 0.92);
}

.hero-panel h1,
.hero-copy,
.label,
.reward-copy,
.vip-card h3,
.vip-card p {
  margin: 0;
}

.hero-panel h1 {
  margin-top: 8px;
  font-size: clamp(28px, 4.8vw, 44px);
}

.hero-copy {
  margin-top: 10px;
  max-width: 640px;
  color: rgba(247, 239, 225, 0.74);
  line-height: 1.6;
}

.status-text {
  margin: 0;
  color: #ffd4d4;
}

.vip-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.vip-summary,
.vip-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(18, 11, 24, 0.78);
}

.vip-summary {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.label {
  color: rgba(251, 191, 36, 0.92);
  font-size: 13px;
  letter-spacing: 0.08em;
}

.vip-summary h2 {
  margin: 0;
  font-size: clamp(24px, 3vw, 34px);
}

.reward-copy,
.next-meta,
.wallet-text,
.benefit-list {
  color: rgba(247, 239, 225, 0.74);
}

.actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.claim-button {
  padding: 12px 18px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #fbbf24, #fb7185);
  color: #2b1203;
  font-weight: 700;
  cursor: pointer;
}

.claim-button:disabled {
  cursor: not-allowed;
  opacity: 0.58;
}

.benefit-list {
  margin: 12px 0 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

@media (max-width: 720px) {
  .vip-grid {
    grid-template-columns: 1fr;
  }
}
</style>
