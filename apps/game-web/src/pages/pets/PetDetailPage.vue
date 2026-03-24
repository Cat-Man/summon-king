<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { legacyInventoryIconMap } from '@/assets/legacy'
import { runtimeConfig } from '@/config/runtime'
import { createInitialPetDetailDashboard, loadPetDetailDashboard } from '@/services/pet-dashboard'
import { useSessionStore } from '@/stores/session'

const route = useRoute()
const sessionStore = useSessionStore()
const dashboard = ref(createInitialPetDetailDashboard())
const loadError = ref('')

const resolvedPetId = computed(() => {
  const value = Array.isArray(route.query.pet_id) ? route.query.pet_id[0] : route.query.pet_id
  const petId = Number(value)
  return Number.isFinite(petId) && petId > 0 ? petId : undefined
})

const evolutionStoneIcon = computed(() => {
  const itemName = dashboard.value.evolutionCost.itemName
  return legacyInventoryIconMap[itemName] ?? (itemName.includes('进化石') ? legacyInventoryIconMap['火系进化石'] : undefined)
})

onMounted(async () => {
  try {
    dashboard.value = await loadPetDetailDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore,
      petId: resolvedPetId.value
    })
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '幻兽详情加载失败'
  }
})
</script>

<template>
  <section class="pet-detail-page">
    <UiPageHero
      :eyebrow="dashboard.hero.eyebrow"
      :title="dashboard.hero.title"
      :description="dashboard.hero.description"
      :tone="dashboard.hero.tone"
      :meta-label="dashboard.hero.metaLabel"
      :meta-value="dashboard.hero.metaValue"
    />

    <p v-if="loadError" class="error-banner">{{ loadError }}</p>

    <UiStatGrid :items="dashboard.overview" tone="warm" min-width="130px" />

    <div v-if="dashboard.heroIcon" class="pet-detail-hero" data-testid="pet-detail-hero-art">
      <img :src="dashboard.heroIcon" :alt="`${dashboard.hero.title} 头像`" />
      <strong>{{ dashboard.hero.title }}</strong>
    </div>

    <section class="grid">
      <UiPanelCard title="技能">
        <ul class="text-list">
          <li v-for="item in dashboard.skills" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="战骨">
        <UiChipGroup :items="dashboard.bones" tone="orange" min-width="120px" />
      </UiPanelCard>

      <UiPanelCard title="战灵">
        <UiChipGroup :items="dashboard.spirits" tone="blue" min-width="140px" />
      </UiPanelCard>

      <UiPanelCard title="魔魂">
        <UiChipGroup :items="dashboard.souls" tone="green" min-width="140px" />
      </UiPanelCard>

      <UiPanelCard title="进化/升境">
        <div class="evolution-cost" data-testid="pet-detail-evolution-cost">
          <img
            v-if="evolutionStoneIcon"
            :src="evolutionStoneIcon"
            :alt="dashboard.evolutionCost.itemName"
          />
          <div>
            <strong>{{ dashboard.evolutionCost.itemName }} {{ dashboard.evolutionCost.itemProgress }}</strong>
            <span>铜钱 {{ dashboard.evolutionCost.coinProgress }}</span>
          </div>
        </div>
        <p>
          当前升境条件：{{ dashboard.evolutionCost.itemName }} {{ dashboard.evolutionCost.itemProgress }}、
          铜钱 {{ dashboard.evolutionCost.coinProgress }}。满足后可直接提交升境并刷新最终属性。
        </p>
      </UiPanelCard>

      <UiPanelCard title="重生/放生">
        <p>重生将返还大部分培养材料；放生将永久失去该幻兽实例，提交前必须二次确认。</p>
      </UiPanelCard>

      <UiPanelCard title="来源说明" wide>
        <ul class="text-list">
          <li v-for="item in dashboard.sources" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-detail-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.error-banner {
  margin: 0;
  padding: 10px 14px;
  border-radius: 12px;
  background: #fef2f2;
  color: #b91c1c;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.text-list {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
  line-height: 1.6;
}

.grid p {
  margin: 0;
  color: #4b5563;
  line-height: 1.7;
}

.pet-detail-hero {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px;
  border-radius: 20px;
  background: #fff7ed;
  color: #b45309;
}

.pet-detail-hero img {
  width: 80px;
  height: 80px;
  border-radius: 999px;
  object-fit: cover;
  border: 2px solid #fb923c;
}

.evolution-cost {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 10px;
}

.evolution-cost img {
  width: 40px;
  height: 40px;
  object-fit: contain;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
