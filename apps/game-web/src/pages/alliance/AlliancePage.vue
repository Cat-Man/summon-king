<template>
  <section class="alliance-page">
    <header class="hero-panel">
      <p class="hero-tag">Alliance</p>
      <h1>联盟大厅</h1>
      <p class="hero-copy">
        先把联盟最小闭环跑通：可创建、可申请、可审批。火能修行和盟战在这一版先保留独立入口，不把逻辑塞进同一轮。
      </p>
    </header>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <article v-if="index?.has_alliance && index.alliance" class="card">
      <header class="card-header">
        <div>
          <h2>{{ index.alliance.name }}</h2>
          <p>Lv.{{ index.alliance.level }} · {{ roleLabel(index.current_role) }}</p>
        </div>
        <RouterLink class="war-link" to="/alliance-war">查看盟战</RouterLink>
      </header>
      <p class="notice-text">{{ index.alliance.notice }}</p>
      <p class="meta-text">成员 {{ index.alliance.member_count }} / {{ index.alliance.member_limit }}</p>

      <div class="grid">
        <section class="panel">
          <h3>成员</h3>
          <ul class="simple-list">
            <li v-for="member in index.alliance.members" :key="member.player_id">
              <span>玩家 {{ member.player_id }}</span>
              <strong>{{ roleLabel(member.role) }}</strong>
            </li>
          </ul>
        </section>

        <section class="panel">
          <h3>建筑</h3>
          <ul class="simple-list">
            <li v-for="building in index.alliance.buildings" :key="building.building_type">
              <span>{{ buildingLabel(building.building_type) }}</span>
              <strong>Lv.{{ building.level }}</strong>
            </li>
          </ul>
        </section>
      </div>

      <section v-if="index.pending_applications?.length" class="panel apply-panel">
        <header class="panel-header">
          <h3>待审批申请</h3>
          <span>{{ index.pending_applications.length }} 条</span>
        </header>
        <ul class="simple-list">
          <li v-for="application in index.pending_applications" :key="application.player_id">
            <span>玩家 {{ application.player_id }}</span>
            <button type="button" class="approve-button" :disabled="isSubmitting" @click="handleApprove(application.player_id)">
              通过
            </button>
          </li>
        </ul>
      </section>
    </article>

    <template v-else>
      <article class="card create-card">
        <header class="card-header">
          <div>
            <h2>创建联盟</h2>
            <p>当前没有联盟，先建立一个最小社交节点。</p>
          </div>
        </header>
        <div class="create-form">
          <input v-model="allianceName" type="text" maxlength="12" placeholder="输入联盟名称" />
          <button type="button" class="create-button" :disabled="isSubmitting" @click="handleCreateAlliance">
            创建联盟
          </button>
        </div>
      </article>

      <article class="card">
        <header class="card-header">
          <div>
            <h2>联盟大厅</h2>
            <p>先加入一个联盟，再扩火能修行和盟战。</p>
          </div>
        </header>
        <ul class="hall-list">
          <li v-for="entry in hall" :key="entry.alliance_id" class="hall-item">
            <div>
              <p>{{ entry.name }}</p>
              <span>Lv.{{ entry.level }} · 成员 {{ entry.member_count }} / {{ entry.member_limit }}</span>
              <small>{{ entry.notice }}</small>
            </div>
            <button type="button" class="apply-button" :disabled="isSubmitting || entry.has_applied" @click="handleApply(entry.alliance_id)">
              {{ entry.has_applied ? "已申请" : "申请加入" }}
            </button>
          </li>
        </ul>
      </article>
    </template>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import {
  applyToAlliance,
  approveAllianceApplication,
  createAlliance,
  getAllianceHall,
  getAllianceIndex,
  type AllianceHallEntry,
  type AllianceIndex,
} from "@/api/modules/alliance"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const errorMessage = ref("")
const isSubmitting = ref(false)
const allianceName = ref("")
const index = ref<AllianceIndex | null>(null)
const hall = ref<AllianceHallEntry[]>([])

async function loadAllianceState() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法进入联盟大厅。"
    return
  }

  try {
    const [nextIndex, nextHall] = await Promise.all([getAllianceIndex(playerId), getAllianceHall(playerId)])
    index.value = nextIndex
    hall.value = nextHall
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "联盟数据加载失败。"
  }
}

async function handleCreateAlliance() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法创建联盟。"
    return
  }
  if (!allianceName.value.trim()) {
    errorMessage.value = "请输入联盟名称。"
    return
  }

  isSubmitting.value = true
  try {
    index.value = await createAlliance(playerId, allianceName.value.trim())
    hall.value = await getAllianceHall(playerId)
    allianceName.value = ""
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "联盟创建失败。"
  } finally {
    isSubmitting.value = false
  }
}

async function handleApply(allianceId: number) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法申请联盟。"
    return
  }

  isSubmitting.value = true
  try {
    await applyToAlliance(playerId, allianceId)
    await loadAllianceState()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "联盟申请失败。"
  } finally {
    isSubmitting.value = false
  }
}

async function handleApprove(applicantPlayerId: number) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法审批申请。"
    return
  }

  isSubmitting.value = true
  try {
    await approveAllianceApplication(playerId, applicantPlayerId)
    await loadAllianceState()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "联盟审批失败。"
  } finally {
    isSubmitting.value = false
  }
}

function roleLabel(role: string) {
  if (role === "leader") {
    return "盟主"
  }
  if (role === "member") {
    return "成员"
  }
  return "未入盟"
}

function buildingLabel(buildingType: string) {
  if (buildingType === "hall") {
    return "议事厅"
  }
  if (buildingType === "fire_forge") {
    return "焚天炉"
  }
  return buildingType
}

onMounted(async () => {
  await loadAllianceState()
})
</script>

<style scoped>
.alliance-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(247, 239, 225, 0.15);
  background:
    radial-gradient(circle at 18% 18%, rgba(244, 114, 182, 0.16), transparent 34%),
    radial-gradient(circle at 82% 14%, rgba(250, 204, 21, 0.14), transparent 30%),
    linear-gradient(180deg, rgba(18, 15, 32, 0.94), rgba(9, 12, 22, 0.96));
}

.hero-tag {
  margin: 0;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 12px;
  color: rgba(251, 113, 133, 0.92);
}

.hero-panel h1 {
  margin: 8px 0;
  font-size: clamp(28px, 4.8vw, 46px);
}

.hero-copy,
.card-header p,
.notice-text,
.hall-item small,
.meta-text {
  margin: 0;
  color: rgba(247, 239, 225, 0.76);
  line-height: 1.7;
}

.status-text {
  margin: 0;
  color: #ffb6a2;
}

.card,
.panel {
  padding: 20px;
  border-radius: 24px;
  border: 1px solid rgba(247, 239, 225, 0.1);
  background: rgba(10, 13, 22, 0.72);
  box-shadow: 0 18px 48px rgba(0, 0, 0, 0.2);
}

.card-header,
.panel-header,
.hall-item,
.simple-list li,
.create-form,
.grid {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.grid {
  align-items: stretch;
  margin-top: 18px;
}

.panel {
  flex: 1;
  padding: 18px;
}

.card-header h2,
.panel h3,
.panel-header h3,
.hall-item p {
  margin: 0;
}

.create-card p {
  margin-top: 6px;
}

.create-form {
  margin-top: 18px;
}

.create-form input {
  flex: 1;
  min-width: 0;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid rgba(247, 239, 225, 0.14);
  background: rgba(255, 255, 255, 0.04);
  color: #fff8f0;
}

.create-button,
.apply-button,
.approve-button,
.war-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 42px;
  padding: 0 16px;
  border-radius: 14px;
  border: 1px solid rgba(247, 189, 120, 0.28);
  background: rgba(247, 189, 120, 0.12);
  color: #fff8f0;
  cursor: pointer;
}

.create-button:disabled,
.apply-button:disabled,
.approve-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.hall-list,
.simple-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 18px 0 0;
  padding: 0;
  list-style: none;
}

.hall-item,
.simple-list li {
  padding: 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.03);
}

.hall-item span,
.simple-list span {
  display: block;
  margin-top: 4px;
  color: rgba(247, 239, 225, 0.74);
}

.apply-panel {
  margin-top: 18px;
}

@media (max-width: 860px) {
  .card-header,
  .panel-header,
  .hall-item,
  .simple-list li,
  .create-form,
  .grid {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
