<script setup lang="ts">
import { onMounted, ref } from 'vue'

import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

import { runtimeConfig } from '@/config/runtime'
import { createMockInventoryDashboard } from '@/mocks/inventory-dashboard'
import {
  loadInventoryDashboard,
  sellInventoryItem,
  useInventoryItem
} from '@/services/inventory-dashboard'
import { useSessionStore } from '@/stores/session'

const sessionStore = useSessionStore()
const dashboard = ref(createMockInventoryDashboard())
const loadError = ref('')
const operationMessage = ref('')
const pendingItemId = ref('')

onMounted(async () => {
  try {
    dashboard.value = await loadInventoryDashboard({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore
    })
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '背包加载失败'
  }
})

async function mutateInventory(action: 'use' | 'sell', itemId: string) {
  pendingItemId.value = itemId
  operationMessage.value = ''
  loadError.value = ''

  try {
    const result =
      action === 'use'
        ? await useInventoryItem({
            dataSource: runtimeConfig.gameDataSource,
            sessionStore,
            itemId
          })
        : await sellInventoryItem({
            dataSource: runtimeConfig.gameDataSource,
            sessionStore,
            itemId
          })

    dashboard.value = result.dashboard
    operationMessage.value = result.message
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '背包操作失败'
  } finally {
    pendingItemId.value = ''
  }
}
</script>

<template>
  <section class="inventory-page">
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

    <section class="grid">
      <UiPanelCard title="资源钱包">
        <div class="wallet-icons">
          <article
            v-for="res in dashboard.walletResources"
            :key="res.label"
            class="wallet-icon"
            :data-testid="`wallet-icon-${res.label}`"
          >
            <img :src="res.icon" :alt="`${res.label} 图标`" />
            <div>
              <strong>{{ res.label }}</strong>
              <span>{{ res.value }}</span>
            </div>
          </article>
        </div>
        <UiStatGrid :items="dashboard.wallet" tone="warm" min-width="130px" />
      </UiPanelCard>

      <UiPanelCard title="普通背包">
        <div class="item-list">
          <article
            v-for="item in dashboard.items"
            :key="item.name"
            class="item-card"
            :data-testid="`inventory-item-${item.name}`"
          >
            <div class="item-card__top">
              <img v-if="item.icon" :src="item.icon" :alt="`${item.name} 图标`" />
              <div>
                <strong>{{ item.name }}</strong>
                <span>{{ item.count }}</span>
              </div>
            </div>
            <p>{{ item.action }}</p>
            <div class="item-card__actions">
              <button
                type="button"
                class="action-button"
                :data-testid="`inventory-use-${item.name}`"
                :disabled="pendingItemId === item.itemId"
                @click="mutateInventory('use', item.itemId)"
              >
                使用1个
              </button>
              <button
                type="button"
                class="action-button action-button--sell"
                :data-testid="`inventory-sell-${item.name}`"
                :disabled="pendingItemId === item.itemId"
                @click="mutateInventory('sell', item.itemId)"
              >
                出售1个
              </button>
            </div>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="使用/出售说明">
        <ul class="text-list">
          <li v-for="rule in dashboard.rules" :key="rule">{{ rule }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="资源流水入口">
        <p>后续这里直接挂接资产流水查询，可按铜钱、元宝、道具获得与消耗来源回溯最近变更。</p>
        <button type="button" class="ghost-button">查看最近流水</button>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.inventory-page {
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

.success-banner {
  margin: 0;
  padding: 10px 14px;
  border-radius: 12px;
  background: #ecfdf5;
  color: #047857;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.item-list,
.text-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-list {
  margin: 0;
}

.item-card {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.item-card__actions {
  display: flex;
  gap: 10px;
  margin-top: 12px;
}

.wallet-icons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.wallet-icon {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border-radius: 14px;
  background: #faf5ff;
}

.wallet-icon img {
  width: 32px;
  height: 32px;
  object-fit: contain;
  padding: 4px;
  border-radius: 8px;
  background: #fff;
}

.item-card__top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}

.item-card__top img {
  width: 40px;
  height: 40px;
  object-fit: contain;
  padding: 4px;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 6px 16px rgb(15 23 42 / 15%);
}

.item-card p,
.text-list,
.grid p {
  margin: 0;
  color: #4b5563;
  line-height: 1.6;
}

.text-list {
  padding-left: 18px;
}

.ghost-button {
  margin-top: 12px;
  padding: 10px 14px;
  border: 1px solid #fcd34d;
  border-radius: 12px;
  background: #fffbeb;
  color: #b45309;
  cursor: pointer;
}

.action-button {
  padding: 8px 12px;
  border: 0;
  border-radius: 10px;
  background: #dbeafe;
  color: #1d4ed8;
  cursor: pointer;
}

.action-button--sell {
  background: #fee2e2;
  color: #b91c1c;
}

.action-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
