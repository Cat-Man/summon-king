<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { runtimeConfig } from '@/config/runtime'
import {
  createInitialPetTeamDashboard,
  loadPetTeamDashboard,
  savePetTeamSelection
} from '@/services/pet-dashboard'
import type { PetTeamMember } from '@/services/pet-dashboard.types'
import { useSessionStore } from '@/stores/session'

const TEAM_SLOT_COUNT = 5

const sessionStore = useSessionStore()
const router = useRouter()
const dashboard = ref(createInitialPetTeamDashboard())
const loadError = ref('')
const operationMessage = ref('')
const pendingSave = ref(false)
const editableTeamPetIds = ref<Array<number | null>>(createEmptyTeamPetIds())
const savedTeamPetIds = ref<Array<number | null>>(createEmptyTeamPetIds())

function formatNumber(value: number) {
  return new Intl.NumberFormat('en-US').format(value)
}

function parseNumberLabel(value: string) {
  return Number(value.replace(/[^\d]/g, '')) || 0
}

function createEmptyTeamPetIds() {
  return Array.from({ length: TEAM_SLOT_COUNT }, () => null as number | null)
}

function normalizeSlotId(slot: string) {
  return slot.replace(/\s+/g, '')
}

function normalizeTeamPetIds(petIds: Array<number | null>) {
  const activePetIds = petIds.filter((petId): petId is number => petId !== null)
  return [
    ...activePetIds,
    ...Array.from({ length: Math.max(0, TEAM_SLOT_COUNT - activePetIds.length) }, () => null as number | null)
  ]
}

function createBenchTeamMember(petId: number, slotIndex: number): PetTeamMember | null {
  const rosterPet = dashboard.value.roster.find((item) => item.petId === petId)
  if (!rosterPet) {
    return null
  }

  return {
    petId,
    slot: `${slotIndex + 1} 号位`,
    name: rosterPet.name,
    role: rosterPet.role,
    level: rosterPet.level,
    power: rosterPet.power,
    icon: rosterPet.icon
  }
}

function createTeamSlotMember(petId: number | null, slotIndex: number): PetTeamMember | null {
  if (!petId) {
    return null
  }

  const currentMember = dashboard.value.team.find((item) => item.petId === petId)
  if (currentMember) {
    return {
      ...currentMember,
      slot: `${slotIndex + 1} 号位`
    }
  }

  return createBenchTeamMember(petId, slotIndex)
}

const equippedPetIds = computed(
  () =>
    new Set(
      editableTeamPetIds.value.flatMap((petId) => (petId ? [petId] : []))
    )
)

const teamSlots = computed(() =>
  Array.from({ length: TEAM_SLOT_COUNT }, (_, slotIndex) => ({
    slotIndex,
    slotLabel: `${slotIndex + 1} 号位`,
    member: createTeamSlotMember(editableTeamPetIds.value[slotIndex] ?? null, slotIndex)
  }))
)

const currentTeam = computed(() =>
  teamSlots.value
    .map((slot) => slot.member)
    .filter((member): member is PetTeamMember => member !== null)
)

const isDirty = computed(
  () => JSON.stringify(editableTeamPetIds.value) !== JSON.stringify(savedTeamPetIds.value)
)

const rosterWithIcons = computed(() =>
  dashboard.value.roster.map((item) => ({
    ...item,
    status: equippedPetIds.value.has(item.petId) ? '已上阵' : '可替补'
  }))
)

const totalPower = computed(() =>
  currentTeam.value.reduce((sum, item) => sum + parseNumberLabel(item.power), 0)
)

const averageLevel = computed(() => {
  if (currentTeam.value.length === 0) {
    return 0
  }

  const totalLevel = currentTeam.value.reduce((sum, item) => sum + parseNumberLabel(item.level), 0)
  return Math.round(totalLevel / currentTeam.value.length)
})

const heroMetaValue = computed(() => formatNumber(totalPower.value))

const overviewItems = computed(() => {
  const focusMember = currentTeam.value[0] ?? rosterWithIcons.value[0]
  const benchCount = rosterWithIcons.value.filter((item) => item.status !== '已上阵').length

  return [
    { label: '已上阵', value: `${currentTeam.value.length} / ${TEAM_SLOT_COUNT}` },
    { label: '队伍核心', value: focusMember?.name ?? '-' },
    { label: '平均等级', value: averageLevel.value > 0 ? `Lv.${averageLevel.value}` : '-' },
    { label: '可替补', value: `${benchCount} 只` }
  ]
})

const focusCard = computed(() => {
  const focusMember = currentTeam.value[0] ?? rosterWithIcons.value[0]
  if (!focusMember) {
    return dashboard.value.focus
  }

  return {
    name: focusMember.name,
    description: `${focusMember.name}当前承担${focusMember.role}定位，建议围绕当前阵容继续补位。`,
    nextStep: '保存阵容后可用于首页战力联动。'
  }
})

function syncEditableTeamFromDashboard() {
  const normalizedTeamPetIds = normalizeTeamPetIds(dashboard.value.team.map((item) => item.petId))
  editableTeamPetIds.value = [...normalizedTeamPetIds]
  savedTeamPetIds.value = [...normalizedTeamPetIds]
}

function clearFeedbackMessages() {
  loadError.value = ''
  operationMessage.value = ''
}

function assignPetToNextSlot(petId: number) {
  if (equippedPetIds.value.has(petId)) {
    return
  }

  if (!editableTeamPetIds.value.some((currentPetId) => currentPetId === null)) {
    loadError.value = '当前战斗队已满，请先下阵后再上阵'
    operationMessage.value = ''
    return
  }

  clearFeedbackMessages()
  editableTeamPetIds.value = normalizeTeamPetIds([...editableTeamPetIds.value, petId])
}

function removePetFromTeam(petId: number) {
  clearFeedbackMessages()
  editableTeamPetIds.value = normalizeTeamPetIds(
    editableTeamPetIds.value.filter((currentPetId) => currentPetId !== petId)
  )
}

function movePet(petId: number, direction: -1 | 1) {
  const activePetIds = editableTeamPetIds.value.filter((currentPetId): currentPetId is number => currentPetId !== null)
  const currentIndex = activePetIds.findIndex((currentPetId) => currentPetId === petId)
  const targetIndex = currentIndex + direction

  if (currentIndex === -1 || targetIndex < 0 || targetIndex >= activePetIds.length) {
    return
  }

  clearFeedbackMessages()
  ;[activePetIds[currentIndex], activePetIds[targetIndex]] = [activePetIds[targetIndex], activePetIds[currentIndex]]
  editableTeamPetIds.value = normalizeTeamPetIds(activePetIds)
}

function navigateHome() {
  router.push({ name: 'home' })
}

onMounted(async () => {
  try {
    dashboard.value = await loadPetTeamDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore
    })
    syncEditableTeamFromDashboard()
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '阵容加载失败'
  }
})

async function saveCurrentTeam() {
  if (pendingSave.value || !isDirty.value) {
    return
  }

  pendingSave.value = true
  operationMessage.value = ''
  loadError.value = ''

  try {
    const activePetIds = editableTeamPetIds.value.flatMap((petId) => (petId ? [petId] : []))
    const result = await savePetTeamSelection({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore,
      petIds: activePetIds
    })
    savedTeamPetIds.value = normalizeTeamPetIds(activePetIds)
    operationMessage.value = result.message
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '阵容保存失败'
  } finally {
    pendingSave.value = false
  }
}
</script>

<template>
  <section class="pet-team-page">
    <UiPageHero
      :eyebrow="dashboard.hero.eyebrow"
      :title="dashboard.hero.title"
      :description="dashboard.hero.description"
      :tone="dashboard.hero.tone"
      :meta-label="dashboard.hero.metaLabel"
      :meta-value="heroMetaValue"
    />

    <p v-if="loadError" class="error-banner">{{ loadError }}</p>
    <div v-if="operationMessage" class="success-banner">
      <span>{{ operationMessage }}</span>
      <button
        type="button"
        class="success-banner__action"
        data-testid="return-home-after-save"
        @click="navigateHome"
      >
        返回首页查看战力
      </button>
    </div>

    <UiStatGrid :items="overviewItems" min-width="140px" />

    <section class="grid">
      <UiPanelCard title="当前战斗队">
        <template #actions>
          <button
            type="button"
            class="save-button"
            data-testid="save-team"
            :disabled="pendingSave || !isDirty"
            @click="saveCurrentTeam"
          >
            {{ pendingSave ? '保存中...' : '保存阵容' }}
          </button>
        </template>
        <div class="team-list">
          <article
            v-for="slot in teamSlots"
            :key="slot.slotIndex"
            class="team-item"
            :data-testid="`team-slot-card-${slot.slotIndex + 1}`"
          >
            <template v-if="slot.member">
              <div
                v-if="slot.member.icon"
                class="team-item__avatar-wrapper"
                :data-testid="`team-slot-${normalizeSlotId(slot.member.slot)}`"
              >
                <img :src="slot.member.icon" :alt="`${slot.member.name} 头像`" class="team-item__avatar" />
              </div>
              <div class="team-item__content" :data-testid="`team-member-${slot.slotIndex + 1}`">
                <span class="subtle">{{ slot.member.slot }}</span>
                <strong>{{ slot.member.name }}</strong>
                <p>{{ slot.member.role }} · {{ slot.member.level }}</p>
              </div>
              <span class="power">战力 {{ slot.member.power }}</span>
              <div class="team-item__actions">
                <button
                  type="button"
                  class="team-action"
                  :data-testid="`team-remove-${slot.member.petId}`"
                  @click="removePetFromTeam(slot.member.petId)"
                >
                  下阵
                </button>
                <button
                  type="button"
                  class="team-action"
                  :data-testid="`team-move-up-${slot.member.petId}`"
                  @click="movePet(slot.member.petId, -1)"
                >
                  前移
                </button>
                <button
                  type="button"
                  class="team-action"
                  :data-testid="`team-move-down-${slot.member.petId}`"
                  @click="movePet(slot.member.petId, 1)"
                >
                  后移
                </button>
              </div>
            </template>
            <template v-else>
              <div class="team-item__empty">
                <span class="subtle">{{ slot.slotLabel }}</span>
                <strong>待上阵</strong>
                <p>从幻兽栏补入当前主力队伍</p>
              </div>
            </template>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="主养成目标">
        <div class="focus-card">
          <strong>{{ focusCard.name }}</strong>
          <p>{{ focusCard.description }}</p>
          <span>{{ focusCard.nextStep }}</span>
        </div>
      </UiPanelCard>

      <UiPanelCard title="幻兽栏" wide>
        <div class="roster-grid">
          <article
            v-for="pet in rosterWithIcons"
            :key="pet.petId"
            class="roster-item"
            :data-testid="`roster-item-${pet.name}`"
          >
            <div v-if="pet.icon" class="roster-item__avatar-wrapper">
              <img :src="pet.icon" :alt="`${pet.name} 头像`" class="roster-item__avatar" />
            </div>
            <strong>{{ pet.name }}</strong>
            <span>{{ pet.level }}</span>
            <em>{{ pet.status }}</em>
            <small>综合战力 {{ pet.power }}</small>
            <button
              v-if="pet.status !== '已上阵'"
              type="button"
              class="roster-action"
              :data-testid="`roster-action-${pet.name}`"
              @click="assignPetToNextSlot(pet.petId)"
            >
              上阵
            </button>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="上阵策略" wide>
        <UiChipGroup :items="dashboard.strategies" tone="blue" min-width="150px" />
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-team-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.error-banner,
.success-banner {
  margin: 0;
  padding: 10px 14px;
  border-radius: 12px;
}

.error-banner {
  background: #fef2f2;
  color: #b91c1c;
}

.success-banner {
  background: #ecfdf5;
  color: #047857;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.team-list,
.roster-grid {
  display: grid;
  gap: 12px;
}

.team-item,
.roster-item,
.focus-card {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.team-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.team-item__content,
.team-item__empty {
  display: grid;
  gap: 6px;
}

.team-item__actions {
  display: grid;
  gap: 6px;
}

.save-button {
  border: 0;
  border-radius: 999px;
  padding: 8px 14px;
  background: #1d4ed8;
  color: #fff;
  font-size: 13px;
  cursor: pointer;
}

.save-button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.success-banner__action {
  border: 0;
  border-radius: 999px;
  padding: 6px 12px;
  background: #047857;
  color: #fff;
  font-size: 12px;
  cursor: pointer;
}

.roster-grid {
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
}

.roster-item {
  display: grid;
  gap: 6px;
}

.team-item__avatar-wrapper,
.roster-item__avatar-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 6px;
}

.team-item__avatar,
.roster-item__avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.focus-card p,
.team-item p {
  margin: 6px 0 0;
  color: #4b5563;
}

.focus-card span,
.roster-item small,
.subtle {
  color: #6b7280;
  font-size: 13px;
}

.power,
.roster-item em {
  font-style: normal;
  color: #1d4ed8;
  font-size: 13px;
}

.roster-action,
.team-action {
  width: fit-content;
  border: 0;
  border-radius: 999px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
}

.roster-action {
  background: #dbeafe;
  color: #1d4ed8;
}

.team-action {
  background: #e2e8f0;
  color: #0f172a;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .team-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
