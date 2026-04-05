<template>
  <section class="map-page">
    <div class="map-hero">
      <p class="badge">参考策略</p>
      <h1>{{ title }}</h1>
      <p>探索城市、占领区域、训练幻兽前的全局视角。</p>
    </div>
    <p v-if="errorMessage" class="status-text">{{ errorMessage }}</p>
    <div class="map-grid">
      <article v-for="city in cities" :key="city.city_id" class="map-card">
        <header>
          <h2>{{ city.name }}</h2>
          <span>{{ city.region }}</span>
        </header>
        <p>坐标 {{ city.loc_x }} / {{ city.loc_y }}</p>
        <div class="dungeon-list">
        <p>可进入副本</p>
        <span class="spirit-text">当前灵力 {{ currentSpiritPower }}</span>
          <span class="growth-text">
            当前成长总加成 +{{ totalGrowthBonus }} · 战骨 +{{ currentBoneBonus }} · 战灵 +{{ currentSpiritBonus }} · 魔魂 +{{ currentSoulBonus }}
          </span>
          <div
            v-for="dungeon in city.dungeons"
            :key="`${city.city_id}-${dungeon.dungeon_id}`"
            class="dungeon-entry"
          >
            <button
              :data-city-id="city.city_id"
              :data-dungeon-id="dungeon.dungeon_id"
              class="ghost-btn"
              :class="{ locked: !canEnterDungeon(dungeon) }"
              type="button"
              :disabled="!canEnterDungeon(dungeon)"
              @click="goToDungeon(dungeon)"
            >
              进入 {{ dungeon.dungeon_name }}
            </button>
            <p v-if="!canEnterDungeon(dungeon)" class="dungeon-summary">
              {{ lockSummary(dungeon) }}
            </p>
            <div v-if="!canEnterDungeon(dungeon)" class="lock-details">
              <p v-for="missing in missingRequirements(dungeon)" :key="missing.label">{{ missing.label }}</p>
              <div class="lock-ctas">
                <RouterLink
                  v-for="missing in missingRequirements(dungeon)"
                  :key="`${city.city_id}-${dungeon.dungeon_id}-${missing.route}`"
                  class="lock-cta"
                  :to="missing.route"
                >
                  {{ missing.ctaLabel }}
                </RouterLink>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { RouterLink, useRouter } from "vue-router"

import { APIError } from "@/api/http"
import { getWorldMap, type WorldMap } from "@/api/modules/dungeon"
import { getGrowthWallet } from "@/api/modules/growth"
import { getPetCollection, type PetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

const world = ref<WorldMap | null>(null)
const errorMessage = ref("")
const router = useRouter()
const sessionStore = useSessionStore()
const resourceSyncStore = useResourceSyncStore()
const currentSpiritPower = ref(0)
const currentBoneLevel = ref(0)
const currentSoulPieces = ref(0)
const petCollection = ref<PetCollection | null>(null)

const title = computed(() => (world.value ? `${world.value.name}世界地图` : "环天世界地图"))
const cities = computed(() => world.value?.cities ?? [])
const currentBoneBonus = computed(() =>
  (petCollection.value?.active_team ?? []).reduce((total, pet) => total + (pet.power_breakdown?.bone ?? 0), 0),
)
const currentSpiritBonus = computed(() =>
  (petCollection.value?.active_team ?? []).reduce((total, pet) => total + (pet.power_breakdown?.spirit ?? 0), 0),
)
const currentSoulBonus = computed(() =>
  (petCollection.value?.active_team ?? []).reduce((total, pet) => total + (pet.power_breakdown?.soul ?? 0), 0),
)
const totalGrowthBonus = computed(
  () => currentBoneBonus.value + currentSpiritBonus.value + currentSoulBonus.value,
)

function canEnterDungeon(dungeon: WorldMap["cities"][number]["dungeons"][number]) {
  const boneReq = dungeon.unlock_bone_level ?? 0
  const soulReq = dungeon.unlock_soul_pieces ?? 0
  return (
    currentSpiritPower.value >= (dungeon.unlock_spirit_power ?? 0) &&
    currentBoneLevel.value >= boneReq &&
    currentSoulPieces.value >= soulReq
  )
}

function goToDungeon(dungeon: WorldMap["cities"][number]["dungeons"][number]) {
  void router.push({
    name: "dungeon",
    query: {
      dungeon_id: String(dungeon.dungeon_id),
    },
  })
}

function lockSummary(dungeon: WorldMap["cities"][number]["dungeons"][number]) {
  const spirit = dungeon.unlock_spirit_power ?? 0
  const bone = dungeon.unlock_bone_level ?? 0
  const soul = dungeon.unlock_soul_pieces ?? 0
  return `需灵力 ${spirit} · 战骨 ${bone} · 魔魂 ${soul}`
}

type MissingRequirement = {
  label: string
  route: string
  ctaLabel: string
}

function missingRequirements(dungeon: WorldMap["cities"][number]["dungeons"][number]) {
  const requirements: MissingRequirement[] = []
  const spiritGap = Math.max(0, (dungeon.unlock_spirit_power ?? 0) - currentSpiritPower.value)
  if (spiritGap > 0) {
    requirements.push({ label: `还差灵力 ${spiritGap}`, route: "/cultivation", ctaLabel: "去修行" })
  }
  const boneGap = Math.max(0, (dungeon.unlock_bone_level ?? 0) - currentBoneLevel.value)
  if (boneGap > 0) {
    requirements.push({ label: `还差战骨 ${boneGap}（约 +${boneGap * 24} 战力）`, route: "/growth/bone", ctaLabel: "去战骨" })
  }
  const soulGap = Math.max(0, (dungeon.unlock_soul_pieces ?? 0) - currentSoulPieces.value)
  if (soulGap > 0) {
    requirements.push({ label: `还差魔魂 ${soulGap}（约 +${soulGap * 8} 战力）`, route: "/growth/soul", ctaLabel: "去魔魂" })
  }
  return requirements
}

async function loadWallet() {
  if (!sessionStore.playerId) {
    return
  }

  const wallet = await getGrowthWallet(sessionStore.playerId)
  currentSpiritPower.value = wallet.spirit_power
  currentBoneLevel.value = wallet.bone_level
  currentSoulPieces.value = wallet.soul_pieces
}

async function loadPetSummary() {
  if (!sessionStore.playerId) {
    return
  }

  petCollection.value = await getPetCollection(sessionStore.playerId)
}

async function loadGrowthState() {
  await Promise.all([loadWallet(), loadPetSummary()])
}

onMounted(async () => {
  try {
    const [nextWorld] = await Promise.all([getWorldMap(), loadGrowthState()])
    world.value = nextWorld
  } catch (error) {
    if (error instanceof APIError) {
      errorMessage.value = error.message
      return
    }
    errorMessage.value = "世界地图加载失败，请稍后重试。"
  }
})

watch(
  () => resourceSyncStore.version,
  async (next, prev) => {
    if (next === prev) {
      return
    }
    await loadGrowthState()
  },
  { flush: "sync" },
)
</script>

<style scoped>
:global(body) {
  background: radial-gradient(circle at 20% 20%, #1a2a4f, #050715 70%);
}
.map-page {
  min-height: 100vh;
  padding: 32px 24px;
  color: #f4f6ff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.status-text {
  margin: 0;
  color: #ffd3c7;
}
.map-hero {
  max-width: 640px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 20px;
  padding: 24px;
  backdrop-filter: blur(12px);
  box-shadow: 0 25px 45px rgba(10, 10, 40, 0.45);
}
.badge {
  display: inline-flex;
  padding: 6px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.15);
  margin-bottom: 8px;
  font-size: 12px;
}
.map-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}
.map-card {
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.02), rgba(255, 255, 255, 0.08));
  border: 1px solid rgba(255, 255, 255, 0.09);
  backdrop-filter: blur(8px);
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: transform 0.3s ease, border-color 0.3s ease;
}
.map-card:hover {
  transform: translateY(-6px);
  border-color: rgba(255, 255, 255, 0.3);
}
.map-card header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.map-card header span {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
}
.dungeon-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.dungeon-list p {
  margin: 0;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
.spirit-text,
.growth-text,
.dungeon-hint {
  color: rgba(255, 255, 255, 0.68);
  font-size: 13px;
}
.ghost-btn {
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 10px;
  background: transparent;
  color: #f4f6ff;
  padding: 10px;
  cursor: pointer;
  transition: background 0.2s ease;
}
.ghost-btn:hover {
  background: rgba(255, 255, 255, 0.08);
}
.ghost-btn.locked,
.ghost-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
.dungeon-entry {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.dungeon-summary {
  margin: 0;
  color: rgba(255, 255, 255, 0.72);
  font-size: 14px;
}
.lock-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-left: 10px;
}
.lock-ctas {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.lock-cta {
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.22);
  color: #f4f6ff;
  text-decoration: none;
}
</style>
