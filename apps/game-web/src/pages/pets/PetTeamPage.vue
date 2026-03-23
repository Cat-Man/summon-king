<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

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
import { useSessionStore } from '@/stores/session'

const sessionStore = useSessionStore()
const dashboard = ref(createInitialPetTeamDashboard())
const loadError = ref('')
const operationMessage = ref('')
const pendingSave = ref(false)

const normalizeSlotId = (slot: string) => slot.replace(/\s+/g, '')
const teamWithIcons = computed(() => dashboard.value.team)
const rosterWithIcons = computed(() => dashboard.value.roster)

onMounted(async () => {
  try {
    dashboard.value = await loadPetTeamDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore
    })
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '阵容加载失败'
  }
})

async function saveCurrentTeam() {
  pendingSave.value = true
  operationMessage.value = ''
  loadError.value = ''

  try {
    const result = await savePetTeamSelection({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore,
      petIds: dashboard.value.team.map((item) => item.petId)
    })
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
      :meta-value="dashboard.hero.metaValue"
    />

    <p v-if="loadError" class="error-banner">{{ loadError }}</p>
    <p v-if="operationMessage" class="success-banner">{{ operationMessage }}</p>

    <UiStatGrid :items="dashboard.overview" min-width="140px" />

    <section class="grid">
      <UiPanelCard title="当前战斗队">
        <template #actions>
          <button
            type="button"
            class="save-button"
            data-testid="save-team"
            :disabled="pendingSave"
            @click="saveCurrentTeam"
          >
            {{ pendingSave ? '保存中...' : '保存阵容' }}
          </button>
        </template>
        <div class="team-list">
          <article v-for="item in teamWithIcons" :key="item.petId" class="team-item">
            <div
              v-if="item.icon"
              class="team-item__avatar-wrapper"
              :data-testid="`team-slot-${normalizeSlotId(item.slot)}`"
            >
              <img :src="item.icon" :alt="`${item.name} 头像`" class="team-item__avatar" />
            </div>
            <div>
              <span class="subtle">{{ item.slot }}</span>
              <strong>{{ item.name }}</strong>
              <p>{{ item.role }} · {{ item.level }}</p>
            </div>
            <span class="power">战力 {{ item.power }}</span>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="主养成目标">
        <div class="focus-card">
          <strong>{{ dashboard.focus.name }}</strong>
          <p>{{ dashboard.focus.description }}</p>
          <span>{{ dashboard.focus.nextStep }}</span>
        </div>
      </UiPanelCard>

      <UiPanelCard title="幻兽栏" wide>
        <div class="roster-grid">
          <article v-for="pet in rosterWithIcons" :key="pet.petId" class="roster-item">
            <div
              v-if="pet.icon"
              class="roster-item__avatar-wrapper"
              :data-testid="`roster-item-${pet.name}`"
            >
              <img :src="pet.icon" :alt="`${pet.name} 头像`" class="roster-item__avatar" />
            </div>
            <strong>{{ pet.name }}</strong>
            <span>{{ pet.level }}</span>
            <em>{{ pet.status }}</em>
            <small>综合战力 {{ pet.power }}</small>
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
  cursor: wait;
  opacity: 0.7;
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
