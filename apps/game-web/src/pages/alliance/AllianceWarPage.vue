<template>
  <section class="war-page">
    <header class="hero-panel">
      <p class="hero-tag">Alliance War</p>
      <h1>盟战前线</h1>
      <p class="hero-copy">
        这一版先把盟战最小闭环跑通：可查看当前阶段、当前目标，并支持盟主登记本轮目标。
      </p>
    </header>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <article v-if="index" class="card">
      <h2>当前阶段</h2>
      <p>{{ phaseLabel(index.phase) }}</p>
      <p class="copy-text">当前目标：{{ index.target_label }}</p>

      <div class="status-grid">
        <p>当前目标</p>
        <strong>{{ index.target_label }}</strong>
        <span>{{ index.can_register ? "可登记目标" : "当前角色不可登记目标" }}</span>
      </div>

      <section v-if="index.has_alliance" class="target-panel">
        <h3>可选目标</h3>
        <div class="target-actions">
          <button
            v-for="target in index.available_targets"
            :key="target"
            :data-testid="`register-target-${target}`"
            class="target-button"
            type="button"
            :disabled="isSubmitting || !index.can_register"
            @click="handleRegisterTarget(target)"
          >
            {{ target }}
          </button>
        </div>
      </section>

      <RouterLink class="back-link" to="/alliance">返回联盟首页</RouterLink>
    </article>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import { getAllianceWarIndex, registerAllianceWarTarget, type AllianceWarIndex } from "@/api/modules/allianceWar"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const index = ref<AllianceWarIndex | null>(null)
const errorMessage = ref("")
const isSubmitting = ref(false)

function normalizeIndex(next?: Partial<AllianceWarIndex> | null): AllianceWarIndex {
  return {
    has_alliance: next?.has_alliance ?? false,
    current_role: next?.current_role ?? "none",
    phase: next?.phase ?? "preparing",
    target_label: next?.target_label ?? "待开放目标",
    can_register: next?.can_register ?? false,
    available_targets: next?.available_targets ?? [],
  }
}

function phaseLabel(phase: string) {
  if (phase === "preparing") {
    return "筹备中"
  }
  if (phase === "fighting") {
    return "进行中"
  }
  if (phase === "settled") {
    return "已结算"
  }
  return phase
}

async function loadIndex() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载盟战状态。"
    return
  }

  try {
    index.value = normalizeIndex(await getAllianceWarIndex(playerId))
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "盟战状态加载失败。"
  }
}

async function handleRegisterTarget(targetLabel: string) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法登记盟战目标。"
    return
  }
  if (!index.value || !index.value.can_register) {
    return
  }

  isSubmitting.value = true
  try {
    index.value = normalizeIndex(await registerAllianceWarTarget(playerId, targetLabel))
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "盟战目标登记失败。"
  } finally {
    isSubmitting.value = false
  }
}

onMounted(async () => {
  await loadIndex()
})
</script>

<style scoped>
.war-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel,
.card {
  padding: 24px;
  border-radius: 24px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(10, 13, 22, 0.72);
}

.hero-panel {
  background:
    radial-gradient(circle at 20% 20%, rgba(248, 113, 113, 0.15), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(56, 189, 248, 0.16), transparent 32%),
    linear-gradient(180deg, rgba(15, 18, 30, 0.95), rgba(10, 12, 20, 0.95));
}

.hero-tag,
.hero-copy,
.copy-text {
  margin: 0;
  color: rgba(247, 239, 225, 0.76);
}

.status-text {
  margin: 0;
  color: #ffd4d4;
}

.hero-tag {
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 12px;
}

.hero-panel h1,
.card h2,
.card p {
  margin: 0;
}

.hero-panel h1 {
  margin: 8px 0;
  font-size: clamp(28px, 4.8vw, 44px);
}

.card {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.status-grid {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 4px 12px;
  padding: 14px;
  border-radius: 14px;
  background: rgba(247, 239, 225, 0.04);
}

.status-grid p,
.status-grid strong,
.status-grid span {
  margin: 0;
}

.status-grid p {
  grid-column: 1 / 2;
  color: rgba(247, 239, 225, 0.64);
}

.status-grid strong {
  grid-column: 2 / 3;
}

.status-grid span {
  grid-column: 1 / -1;
  color: rgba(247, 239, 225, 0.7);
  font-size: 13px;
}

.target-panel {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 6px;
}

.target-panel h3 {
  margin: 0;
  font-size: 16px;
}

.target-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.target-button {
  border: 1px solid rgba(247, 189, 120, 0.28);
  border-radius: 12px;
  background: rgba(247, 189, 120, 0.12);
  color: #fff8f0;
  padding: 8px 12px;
  cursor: pointer;
}

.target-button:disabled {
  opacity: 0.58;
  cursor: not-allowed;
}

.back-link {
  display: inline-flex;
  align-self: flex-start;
  margin-top: 8px;
  padding: 10px 14px;
  border-radius: 14px;
  border: 1px solid rgba(247, 189, 120, 0.28);
  background: rgba(247, 189, 120, 0.12);
  color: #fff8f0;
  text-decoration: none;
}
</style>
