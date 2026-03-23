<script setup lang="ts">
import { onMounted, ref } from 'vue'

import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { runtimeConfig } from '@/config/runtime'
import { createInitialPetCatalogDashboard, loadPetCatalogDashboard } from '@/services/pet-dashboard'
import { useSessionStore } from '@/stores/session'

const sessionStore = useSessionStore()
const dashboard = ref(createInitialPetCatalogDashboard())
const loadError = ref('')

onMounted(async () => {
  try {
    dashboard.value = await loadPetCatalogDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore
    })
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '图鉴加载失败'
  }
})
</script>

<template>
  <section class="pet-list-page">
    <UiPageHero
      :eyebrow="dashboard.hero.eyebrow"
      :title="dashboard.hero.title"
      :description="dashboard.hero.description"
      :tone="dashboard.hero.tone"
      :meta-label="dashboard.hero.metaLabel"
      :meta-value="dashboard.hero.metaValue"
    />

    <p v-if="loadError" class="error-banner">{{ loadError }}</p>

    <UiPanelCard title="图鉴统计">
      <UiStatGrid :items="dashboard.overview" min-width="150px" />
    </UiPanelCard>

    <section class="grid">
      <UiPanelCard title="获取来源">
        <ul class="text-list">
          <li v-for="item in dashboard.sources" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="图鉴列表" wide>
        <div class="catalog-grid">
          <article
            v-for="pet in dashboard.pets"
            :key="pet.petId"
            class="catalog-item"
          >
            <div
              v-if="pet.icon"
              class="catalog-item__avatar-wrapper"
              :data-testid="`pet-card-${pet.name}`"
            >
              <img :src="pet.icon" :alt="`${pet.name} 头像`" class="catalog-item__avatar" />
            </div>
            <header class="catalog-item__header">
              <strong>{{ pet.name }}</strong>
              <span :class="['status', pet.status === '已拥有' ? 'status--owned' : 'status--locked']">
                {{ pet.status }}
              </span>
            </header>
            <p>
              <span>所属地图</span>
              {{ pet.map }}
            </p>
            <p>
              <span>技能池</span>
              {{ pet.skillPool }}
            </p>
            <p>
              <span>培养摘要</span>
              {{ pet.aptitude }}
            </p>
            <p>
              <span>当前等级</span>
              {{ pet.level }}
            </p>
            <p>
              <span>来源</span>
              {{ pet.source }}
            </p>
          </article>
        </div>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-list-page {
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

.catalog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.catalog-item {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
  display: grid;
  gap: 8px;
}

.catalog-item__avatar-wrapper {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 8px;
}

.catalog-item__avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  object-fit: cover;
}

.catalog-item__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.catalog-item p {
  margin: 0;
  display: grid;
  gap: 4px;
  color: #4b5563;
  line-height: 1.5;
}

.catalog-item p span {
  color: #6b7280;
  font-size: 12px;
}

.status {
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.status--owned {
  background: #dcfce7;
  color: #166534;
}

.status--locked {
  background: #ede9fe;
  color: #6d28d9;
}
</style>
