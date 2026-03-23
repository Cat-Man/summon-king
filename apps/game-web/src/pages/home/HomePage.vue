<script setup lang="ts">
import { onMounted, ref } from 'vue'

import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

import { legacyUiAssets } from '@/assets/legacy'
import { runtimeConfig } from '@/config/runtime'
import { createInitialHomeDashboard, loadHomeDashboard } from '@/services/home-dashboard'
import { useSessionStore } from '@/stores/session'

const sessionStore = useSessionStore()
const dashboard = ref(createInitialHomeDashboard())
const loadError = ref('')

onMounted(async () => {
  try {
    dashboard.value = await loadHomeDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore
    })
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '首页加载失败'
  }
})
</script>

<template>
  <section class="home-page">
    <UiPageHero
      :eyebrow="dashboard.hero.eyebrow"
      :title="dashboard.hero.title"
      :description="dashboard.hero.description"
      :tone="dashboard.hero.tone"
      :meta-label="dashboard.hero.metaLabel"
      :meta-value="dashboard.hero.metaValue"
    />

    <p v-if="loadError" class="error-banner">{{ loadError }}</p>

    <UiPanelCard title="每日必做">
      <div class="todo-list">
        <article v-for="item in dashboard.dailyTodos" :key="item.title" class="todo-item">
          <strong>{{ item.title }}</strong>
          <span>{{ item.value }}</span>
          <em>{{ item.action }}</em>
        </article>
      </div>
    </UiPanelCard>

    <section class="resource-strip">
      <article v-for="res in dashboard.resourceIcons" :key="res.label" class="resource-chip">
        <img :src="res.icon" :alt="`${res.label}图标`" :data-testid="`resource-icon-${res.label}`" />
        <div>
          <strong>{{ res.label }}</strong>
          <span>{{ res.value }}</span>
        </div>
      </article>
      <article class="activity-entry" data-testid="activity-entry-icon" :style="{ backgroundImage: `url(${legacyUiAssets.sectionTitleBg})` }">
        <img :src="dashboard.activityEntry.icon" alt="活动入口图标" />
        <div>
          <strong>{{ dashboard.activityEntry.title }}</strong>
          <p>{{ dashboard.activityEntry.description }}</p>
          <small>{{ dashboard.activityEntry.note }}</small>
        </div>
      </article>
    </section>

    <section class="content-grid">
      <UiPanelCard title="资源总览">
        <UiStatGrid :items="dashboard.resources" min-width="110px" />
      </UiPanelCard>

      <UiPanelCard title="消息流入口">
        <ul class="message-list">
          <li v-for="item in dashboard.messages" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="功能矩阵" wide>
        <UiChipGroup :items="dashboard.entries" tone="blue" min-width="100px" />
      </UiPanelCard>
    </section>
  </section>
</template>
<style scoped>
.home-page {
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

.resource-strip {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

.resource-chip,
.activity-entry {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px 16px;
  border-radius: 14px;
  background: #fff;
  box-shadow: 0 5px 18px rgb(15 23 42 / 12%);
}

.resource-chip img,
.activity-entry img {
  width: 42px;
  height: 42px;
  object-fit: contain;
  padding: 4px;
  border-radius: 12px;
  background: #fff;
  border: 2px solid #f8fafc;
}

.resource-chip strong {
  font-size: 14px;
}

.resource-chip span {
  color: #374151;
  font-size: 12px;
}

.activity-entry {
  background-size: cover;
  background-position: center;
  padding: 16px;
}

.activity-entry div strong {
  display: block;
  font-size: 14px;
}

.activity-entry div p {
  margin: 0;
  font-size: 12px;
  color: #1f2937;
}

.activity-entry div small {
  color: #0f172a;
  font-size: 11px;
}

.todo-list,
.content-grid {
  display: grid;
  gap: 14px;
}

.todo-list {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.todo-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.todo-item span {
  color: #4b5563;
}

.todo-item em {
  font-style: normal;
  color: #2563eb;
  font-size: 13px;
}

.content-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.message-list {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
}

@media (max-width: 768px) {
  .resource-strip {
    grid-template-columns: 1fr;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
