<template>
  <section class="signin-page">
    <header class="hero-panel">
      <p class="hero-tag">Daily Signin</p>
      <h1>七日签到</h1>
      <p class="hero-copy">
        先把非支付商业化最小闭环接通：展示今日状态、连续天数和领取结果，不引入商城或额外交互。
      </p>
    </header>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <div v-if="index" class="content-grid">
      <article class="summary-card">
        <p class="label">{{ index.signed_today ? "今日已签到" : "今日可签到" }}</p>
        <h2>连续签到 {{ index.current_streak }} 天</h2>
        <p class="reward-copy">
          今日奖励：灵力 +{{ index.next_reward.spirit_power }}
          <span v-if="index.next_reward.soul_pieces"> · 魔魂碎片 +{{ index.next_reward.soul_pieces }}</span>
        </p>
        <div class="actions">
          <button
            data-testid="signin-claim"
            type="button"
            class="claim-button"
            :disabled="isSubmitting || index.signed_today"
            @click="handleClaim"
          >
            {{ index.signed_today ? "今日已签到" : "立即签到" }}
          </button>
          <span class="wallet-text">钱包灵力 {{ index.wallet_snapshot.spirit_power }}</span>
        </div>
      </article>

      <article class="detail-card">
        <h3>签到节奏</h3>
        <ul class="detail-list">
          <li>今日状态：{{ index.signed_today ? "已完成" : "待领取" }}</li>
          <li>当前连签：{{ index.current_streak }} 天</li>
          <li>战骨收益：+{{ index.next_reward.bone_level }}</li>
          <li>魔魂收益：+{{ index.next_reward.soul_pieces }}</li>
        </ul>
      </article>

      <article class="detail-card" v-if="lastReward">
        <h3>本次到账</h3>
        <p>灵力 +{{ lastReward.spirit_power }}</p>
        <p v-if="lastReward.soul_pieces">魔魂碎片 +{{ lastReward.soul_pieces }}</p>
        <p v-if="lastReward.bone_level">战骨 +{{ lastReward.bone_level }}</p>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"

import { APIError } from "@/api/http"
import { claimSignin, getSigninIndex, type SigninClaimResult, type SigninIndex, type SigninReward } from "@/api/modules/signin"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const index = ref<SigninIndex | null>(null)
const lastReward = ref<SigninReward | null>(null)
const errorMessage = ref("")
const isSubmitting = ref(false)

function normalizeReward(next?: Partial<SigninReward> | null): SigninReward {
  return {
    spirit_power: next?.spirit_power ?? 0,
    bone_level: next?.bone_level ?? 0,
    soul_pieces: next?.soul_pieces ?? 0,
  }
}

function normalizeIndex(next?: Partial<SigninIndex> | null, playerId = 0): SigninIndex {
  return {
    player_id: next?.player_id ?? playerId,
    signed_today: next?.signed_today ?? false,
    current_streak: next?.current_streak ?? 0,
    next_reward: normalizeReward(next?.next_reward),
    wallet_snapshot: {
      spirit_power: next?.wallet_snapshot?.spirit_power ?? 0,
      bone_level: next?.wallet_snapshot?.bone_level ?? 0,
      soul_pieces: next?.wallet_snapshot?.soul_pieces ?? 0,
    },
  }
}

async function loadSignin() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载签到状态。"
    return
  }

  try {
    index.value = normalizeIndex(await getSigninIndex(playerId), playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "签到状态加载失败。"
  }
}

function applyClaimResult(result: SigninClaimResult) {
  index.value = normalizeIndex(result.index, result.index.player_id)
  lastReward.value = normalizeReward(result.reward_delta)
}

async function handleClaim() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法执行签到。"
    return
  }
  if (!index.value || index.value.signed_today) {
    return
  }

  isSubmitting.value = true
  try {
    applyClaimResult(await claimSignin(playerId))
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "签到领取失败。"
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  await loadSignin()
})
</script>

<style scoped>
.signin-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(247, 239, 225, 0.15);
  background:
    radial-gradient(circle at 18% 18%, rgba(96, 165, 250, 0.18), transparent 34%),
    radial-gradient(circle at 82% 14%, rgba(251, 191, 36, 0.18), transparent 30%),
    linear-gradient(180deg, rgba(12, 20, 38, 0.94), rgba(8, 12, 24, 0.96));
}

.hero-tag {
  margin: 0;
  font-size: 12px;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: rgba(125, 211, 252, 0.92);
}

.hero-panel h1,
.hero-copy,
.label,
.reward-copy,
.detail-card h3,
.detail-card p {
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

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(260px, 1fr);
  gap: 18px;
}

.summary-card,
.detail-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(9, 13, 25, 0.78);
}

.summary-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.label {
  color: rgba(125, 211, 252, 0.92);
  font-size: 13px;
  letter-spacing: 0.08em;
}

.summary-card h2 {
  margin: 0;
  font-size: clamp(24px, 3vw, 34px);
}

.reward-copy {
  color: rgba(247, 239, 225, 0.78);
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
  background: linear-gradient(135deg, #60a5fa, #fbbf24);
  color: #111827;
  font-weight: 700;
  cursor: pointer;
}

.claim-button:disabled {
  cursor: not-allowed;
  opacity: 0.58;
}

.wallet-text,
.detail-list {
  color: rgba(247, 239, 225, 0.7);
}

.detail-list {
  margin: 12px 0 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

@media (max-width: 720px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
