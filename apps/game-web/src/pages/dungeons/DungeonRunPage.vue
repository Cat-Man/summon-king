<template>
  <section class="dungeon-page">
    <div class="dungeon-meta">
      <h1>地下城跑动</h1>
      <p>当前副本 {{ currentDungeonName }}，掷骰进入楼层，打响各路 Boss 抢夺装备与材料。</p>
      <div class="dungeon-switch">
        <button
          v-for="option in dungeonOptions"
          :key="option.id"
          :data-dungeon-id="option.id"
          class="switch-btn"
          :class="{ active: option.id === selectedDungeonId }"
          type="button"
          @click="selectedDungeonId = option.id"
        >
          {{ option.name }}
        </button>
      </div>
      <div class="stat-grid">
        <article>
          <h2>{{ run.current_floor }}</h2>
          <p>当前层</p>
        </article>
        <article>
          <h2>{{ run.remain_dice }}</h2>
          <p>剩余骰子</p>
        </article>
        <article>
          <h2>{{ run.status }}</h2>
          <p>状态</p>
        </article>
      </div>
    </div>
    <div class="reward-panel">
      <article>
        <p>{{ run.last_reward.label || "本次掉落" }}</p>
        <strong>灵力 +{{ run.last_reward.spirit_power }}</strong>
        <span>魂力 +{{ run.last_reward.soul_pieces }}</span>
      </article>
      <article>
        <p>当前资源</p>
        <strong>当前灵力 {{ run.wallet_snapshot.spirit_power }}</strong>
        <span>当前魂力 {{ run.wallet_snapshot.soul_pieces }}</span>
      </article>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="action-panel">
      <button class="action-btn primary" type="button" @click="rollForward">掷骰推进</button>
      <button class="action-btn restart-btn" type="button" @click="restartRun">重新进入副本</button>
      <button class="action-btn ghost" type="button" @click="refreshRun">刷新当前状态</button>
    </div>
    <div class="timeline">
      <p>进度</p>
      <div class="floors">
        <div v-for="n in floors" :key="n" class="floor" :class="{ boss: n % 5 === 0 }">
          <span>{{ n }}</span>
          <small>{{ n % 5 === 0 ? 'Boss' : '怪物' }}</small>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { useRoute } from "vue-router"

import { APIError } from "@/api/http"
import { enterDungeon, getDungeonStatus, rollDungeonDice, type DungeonRun } from "@/api/modules/dungeon"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const dungeonOptions = [
  { id: 1, name: "妖窟试炼" },
  { id: 2, name: "寒渊裂隙" },
]

const defaultRun: DungeonRun = {
  player_id: 0,
  dungeon_id: 1,
  remain_dice: 0,
  current_floor: 0,
  status: "idle",
  started_at: "",
  last_reward: {
    label: "",
    spirit_power: 0,
    soul_pieces: 0,
  },
  wallet_snapshot: {
    player_id: 0,
    spirit_power: 0,
    spirit_free_wash: 0,
    bone_level: 0,
    soul_pieces: 0,
    manor_plots: 0,
  },
}

const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const route = useRoute()
const run = ref<DungeonRun>(defaultRun)
const errorMessage = ref("")
const selectedDungeonId = ref(1)

const floors = computed(() => Math.max(6, run.value.current_floor + 2))
const currentDungeonName = computed(() => {
  return dungeonOptions.find((option) => option.id === selectedDungeonId.value)?.name ?? "妖窟试炼"
})

function parseDungeonId(value: unknown) {
  const rawValue = Array.isArray(value) ? value[0] : value
  const dungeonId = Number(rawValue)

  if (dungeonOptions.some((option) => option.id === dungeonId)) {
    return dungeonId
  }

  return null
}

function syncSelectedDungeonId(value: unknown) {
  const dungeonId = parseDungeonId(value)
  if (dungeonId !== null) {
    selectedDungeonId.value = dungeonId
  }
}

function resolveTargetDungeonId() {
  return parseDungeonId(route.query.dungeon_id) ?? selectedDungeonId.value
}

async function refreshRun() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法加载副本。"
    return
  }

  try {
    const currentRun = await getDungeonStatus(playerId)
    const targetDungeonId = resolveTargetDungeonId()

    if (currentRun.dungeon_id && currentRun.dungeon_id !== targetDungeonId) {
      selectedDungeonId.value = targetDungeonId
      run.value = await enterDungeon(playerId, targetDungeonId)
      errorMessage.value = ""
      return
    }

    run.value = currentRun
    selectedDungeonId.value = currentRun.dungeon_id || selectedDungeonId.value
    errorMessage.value = ""
  } catch (error) {
    if (error instanceof APIError && error.status === 404) {
      run.value = await enterDungeon(playerId, selectedDungeonId.value)
      errorMessage.value = ""
      return
    }
    errorMessage.value = error instanceof APIError ? error.message : "副本状态加载失败。"
  }
}

async function rollForward() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法推进副本。"
    return
  }

  try {
    run.value = await rollDungeonDice(playerId)
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "掷骰失败，请稍后重试。"
  }
}

async function restartRun() {
  const playerId = sessionStore.playerId
  if (!playerId) {
    errorMessage.value = "当前未登录，无法进入副本。"
    return
  }

  try {
    run.value = await enterDungeon(playerId, selectedDungeonId.value)
    errorMessage.value = ""
    resourceSyncStore.touch()
  } catch (error) {
    errorMessage.value = error instanceof APIError ? error.message : "进入副本失败，请稍后重试。"
  }
}

onMounted(async () => {
  syncSelectedDungeonId(route.query.dungeon_id)
  await refreshRun()
})

watch(
  () => route.query.dungeon_id,
  async (value) => {
    const dungeonId = parseDungeonId(value)
    if (dungeonId === null) {
      return
    }

    selectedDungeonId.value = dungeonId

    const playerId = sessionStore.playerId
    if (!playerId || run.value.player_id === 0 || run.value.dungeon_id === dungeonId) {
      return
    }

    try {
      run.value = await enterDungeon(playerId, dungeonId)
      errorMessage.value = ""
    } catch (error) {
      errorMessage.value = error instanceof APIError ? error.message : "进入副本失败，请稍后重试。"
    }
  },
)
</script>

<style scoped>
.dungeon-page {
  min-height: 100vh;
  padding: 32px;
  background: linear-gradient(180deg, #120c1f, #090515 70%);
  color: #f8fbff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.status-text {
  margin: 0;
  color: #ffcfb8;
}
.dungeon-switch {
  margin-top: 16px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
.switch-btn {
  padding: 10px 14px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.78);
  cursor: pointer;
}
.switch-btn.active {
  border-color: rgba(255, 204, 51, 0.5);
  background: rgba(255, 204, 51, 0.14);
  color: #ffdd7a;
}
.reward-panel {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}
.reward-panel article {
  padding: 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.reward-panel p,
.reward-panel strong,
.reward-panel span {
  display: block;
}
.reward-panel p {
  margin: 0 0 8px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}
.reward-panel strong {
  font-size: 24px;
}
.reward-panel span {
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.78);
}
.dungeon-meta {
  border-radius: 24px;
  padding: 28px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 30px 60px rgba(5, 5, 20, 0.6);
}
.stat-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}
.stat-grid article {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.stat-grid h2 {
  font-size: 32px;
  margin-bottom: 4px;
}
.action-panel {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}
.action-btn {
  flex: 1;
  min-width: 180px;
  border-radius: 15px;
  padding: 14px 20px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  background: transparent;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, background 0.2s ease;
}
.action-btn.primary {
  background: linear-gradient(135deg, #ffb347, #ffcc33);
  color: #1b1000;
  border: none;
}
.action-btn.ghost {
  border-color: rgba(255, 255, 255, 0.3);
}
.action-btn:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.1);
}
.timeline {
  border-radius: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.floors {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
  gap: 12px;
}
.floor {
  border-radius: 12px;
  padding: 12px;
  text-align: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
.floor.boss {
  background: linear-gradient(180deg, rgba(255, 94, 98, 0.1), rgba(255, 94, 98, 0.3));
  border-color: rgba(255, 94, 98, 0.7);
}
.floor small {
  display: block;
  font-size: 12px;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.7);
}
</style>
