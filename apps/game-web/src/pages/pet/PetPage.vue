<template>
  <section class="pet-page">
    <header class="hero-panel">
      <p class="hero-tag">Pet</p>
      <h1>幻兽阵容</h1>
      <p class="hero-copy">
        当前总战力 {{ collection.total_power }}，已上阵 {{ collection.team_size }} 只。先用最小入口跑通战斗队与幻兽栏，再扩编辑与养成细节。
      </p>
    </header>

    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>

    <article class="card">
      <header class="card-header">
        <h2>战斗队</h2>
        <span>{{ collection.active_team.length }} / 2</span>
      </header>
      <ul class="pet-list">
        <li v-for="pet in collection.active_team" :key="`active-${pet.pet_id}`" class="pet-item">
          <div>
            <p>{{ pet.name }}</p>
            <span>槽位 {{ pet.slot }} · Lv.{{ pet.level }} · EXP {{ pet.exp ?? 0 }}/{{ pet.next_level_exp ?? 100 }}</span>
          </div>
          <strong>{{ pet.power }}</strong>
        </li>
      </ul>
    </article>

    <article class="card">
      <header class="card-header">
        <h2>幻兽栏</h2>
        <span>{{ collection.roster.length }} 只</span>
      </header>
      <ul class="pet-list">
        <li v-for="pet in collection.roster" :key="`roster-${pet.pet_id}`" class="pet-item">
          <div>
            <p>{{ pet.name }}</p>
            <span>
              槽位 {{ pet.slot }} · Lv.{{ pet.level }} · EXP {{ pet.exp ?? 0 }}/{{ pet.next_level_exp ?? 100 }} ·
              {{ pet.is_active ? "上阵中" : "待命" }}
            </span>
          </div>
          <div class="pet-actions">
            <strong>{{ pet.power }}</strong>
            <template v-if="pet.is_active">
              <span v-if="pet.slot === 1" class="status-chip">当前主战</span>
              <div v-else class="action-stack">
                <button
                  :data-testid="`set-main-${pet.pet_id}`"
                  type="button"
                  class="set-main-button"
                  :disabled="isEditing"
                  @click="handleSetMainPet(pet.pet_id)"
                >
                  设为主战
                </button>
                <button
                  :data-testid="`remove-team-${pet.pet_id}`"
                  type="button"
                  class="remove-button"
                  :disabled="isEditing"
                  @click="handleRemoveFromTeam(pet.pet_id)"
                >
                  移出队伍
                </button>
              </div>
            </template>
            <div v-else class="action-stack">
              <button
                v-if="!isTeamFull"
                :data-testid="`join-team-${pet.pet_id}`"
                type="button"
                class="add-button"
                :disabled="isEditing"
                @click="handleAddToTeam(pet.pet_id)"
              >
                加入队伍
              </button>
              <button
                :data-testid="`set-main-${pet.pet_id}`"
                type="button"
                class="set-main-button"
                :disabled="isEditing"
                @click="handleSetMainPet(pet.pet_id)"
              >
                设为主战
              </button>
            </div>
          </div>
        </li>
      </ul>
    </article>

    <article class="card growth-actions">
      <header class="card-header">
        <h2>养成入口</h2>
        <span>最小联动</span>
      </header>
      <div class="actions-grid">
        <RouterLink to="/growth/spirit">前往战灵</RouterLink>
        <RouterLink to="/growth/bone">前往战骨</RouterLink>
        <RouterLink to="/growth/soul">前往魔魂</RouterLink>
      </div>
    </article>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { RouterLink } from "vue-router"

import { APIError } from "@/api/http"
import { getPetCollection, savePetTeam, setMainPet, type PetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const errorMessage = ref("")
const isEditing = ref(false)
const teamCapacity = 2
const collection = ref<PetCollection>({
  player_id: 0,
  total_power: 0,
  team_size: 0,
  active_team: [],
  roster: [],
})

const isTeamFull = computed(() => collection.value.active_team.length >= teamCapacity)
const activePetIds = computed(() => collection.value.active_team.map((pet) => pet.pet_id))

async function loadCollection() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载幻兽阵容。"
    return
  }

  try {
    collection.value = await getPetCollection(playerId)
    errorMessage.value = ""
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "幻兽阵容加载失败。"
  }
}

async function handleSetMainPet(petId: number) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法切换主战幻兽。"
    return
  }

  await runMutation(setMainPet(playerId, petId), "主战幻兽切换失败。")
}

async function handleAddToTeam(petId: number) {
  if (isTeamFull.value) {
    return
  }

  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法保存幻兽阵容。"
    return
  }

  await runMutation(savePetTeam(playerId, [...activePetIds.value, petId]), "幻兽阵容保存失败。")
}

async function handleRemoveFromTeam(petId: number) {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法保存幻兽阵容。"
    return
  }

  await runMutation(
    savePetTeam(playerId, activePetIds.value.filter((id) => id !== petId)),
    "幻兽阵容保存失败。",
  )
}

async function runMutation(promise: Promise<PetCollection>, fallbackMessage: string) {
  isEditing.value = true
  try {
    collection.value = await promise
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : fallbackMessage
  } finally {
    isEditing.value = false
  }
}

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadCollection()
  },
)

onMounted(async () => {
  await loadCollection()
})
</script>

<style scoped>
.pet-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(247, 239, 225, 0.15);
  background:
    radial-gradient(circle at 15% 15%, rgba(125, 211, 252, 0.2), transparent 35%),
    radial-gradient(circle at 85% 15%, rgba(250, 204, 21, 0.15), transparent 32%),
    linear-gradient(180deg, rgba(12, 17, 31, 0.92), rgba(9, 12, 22, 0.95));
}

.hero-tag {
  margin: 0;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 12px;
  color: rgba(147, 197, 253, 0.95);
}

.hero-panel h1 {
  margin: 8px 0;
  font-size: clamp(28px, 4.8vw, 46px);
}

.hero-copy {
  margin: 0;
  color: rgba(247, 239, 225, 0.78);
  line-height: 1.7;
}

.status-text {
  margin: 0;
  color: #ffb6a2;
}

.card {
  padding: 20px;
  border-radius: 20px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(10, 13, 22, 0.86);
  box-shadow: 0 14px 30px rgba(0, 0, 0, 0.25);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 14px;
}

.card-header h2 {
  margin: 0;
}

.card-header span {
  color: rgba(247, 239, 225, 0.7);
  font-size: 13px;
}

.pet-list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pet-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 14px;
  padding: 12px 14px;
  background: rgba(247, 239, 225, 0.05);
}

.pet-item p {
  margin: 0;
  font-size: 16px;
}

.pet-item span {
  font-size: 13px;
  color: rgba(247, 239, 225, 0.68);
}

.pet-item strong {
  color: #bae6fd;
}

.pet-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.action-stack {
  display: flex;
  gap: 6px;
}

.status-chip {
  font-size: 12px;
  color: #bbf7d0;
}

.set-main-button {
  min-width: 86px;
  border: 1px solid rgba(186, 230, 253, 0.45);
  border-radius: 10px;
  padding: 6px 10px;
  background: rgba(125, 211, 252, 0.12);
  color: #dbeafe;
  cursor: pointer;
}

.set-main-button:disabled {
  opacity: 0.5;
  cursor: default;
}

.add-button,
.remove-button {
  border: 1px solid rgba(186, 230, 253, 0.45);
  border-radius: 8px;
  padding: 6px 10px;
  background: rgba(125, 211, 252, 0.12);
  color: #dbeafe;
  cursor: pointer;
  font-size: 12px;
}

.remove-button {
  border-color: rgba(249, 115, 22, 0.65);
  background: rgba(249, 115, 22, 0.12);
  color: #fb923c;
}

.growth-actions .actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 10px;
}

.growth-actions a {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  min-height: 42px;
  border-radius: 12px;
  border: 1px solid rgba(147, 197, 253, 0.45);
  background: rgba(125, 211, 252, 0.1);
  color: #dbeafe;
}

@media (max-width: 720px) {
  .hero-panel {
    padding: 22px;
  }
}
</style>
