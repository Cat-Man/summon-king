<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

import { runtimeConfig } from '@/config/runtime'
import { createMockInventoryDashboard } from '@/mocks/inventory-dashboard'
import {
  loadInventoryDashboard,
  loadRecentAssetLogs,
  sellInventoryItem,
  useInventoryItem
} from '@/services/inventory-dashboard'
import type {
  InventoryItemCard,
  InventoryLogChangeType,
  InventoryLogItem
} from '@/services/inventory-dashboard.types'
import { useSessionStore } from '@/stores/session'

const sessionStore = useSessionStore()
const dashboard = ref(createMockInventoryDashboard())
const loadError = ref('')
const operationMessage = ref('')
const pendingItemId = ref('')
const unlockedHighValueItemIds = ref<string[]>([])
const pendingSellItemId = ref('')
const inventoryLogs = ref<InventoryLogItem[]>([])
const logsPage = ref(1)
const logsPageSize = ref(2)
const logsTotal = ref(0)
const logsChangeType = ref<InventoryLogChangeType>('all')
const logsVisible = ref(false)
const logsLoading = ref(false)
const logsError = ref('')

const logFilterOptions: Array<{ value: InventoryLogChangeType; label: string }> = [
  { value: 'all', label: '全部' },
  { value: 'grant_reward', label: '发放奖励' },
  { value: 'inventory_use', label: '使用道具' },
  { value: 'inventory_sell', label: '出售道具' }
]

const logsTotalPages = computed(() => {
  const totalPages = Math.ceil(logsTotal.value / logsPageSize.value)
  return Math.max(1, totalPages)
})

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
    if (logsVisible.value) {
      await refreshInventoryLogs()
    }
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : '背包操作失败'
  } finally {
    pendingItemId.value = ''
  }
}

function isItemUnlocked(itemId: string) {
  return unlockedHighValueItemIds.value.includes(itemId)
}

function isItemLocked(item: InventoryItemCard) {
  return item.isHighValue && !isItemUnlocked(item.itemId)
}

function toggleItemLock(item: InventoryItemCard) {
  if (!item.isHighValue) {
    return
  }

  if (isItemUnlocked(item.itemId)) {
    unlockedHighValueItemIds.value = unlockedHighValueItemIds.value.filter((entry) => entry !== item.itemId)
    return
  }

  unlockedHighValueItemIds.value = [...unlockedHighValueItemIds.value, item.itemId]
}

function openSellConfirmation(item: InventoryItemCard) {
  if (isItemLocked(item)) {
    return
  }

  pendingSellItemId.value = item.itemId
  operationMessage.value = ''
  loadError.value = ''
}

function closeSellConfirmation() {
  pendingSellItemId.value = ''
}

function getPendingSellItem() {
  return dashboard.value.items.find((item) => item.itemId === pendingSellItemId.value) ?? null
}

async function refreshInventoryLogs() {
  logsLoading.value = true
  logsError.value = ''

  try {
    const result = await loadRecentAssetLogs({
      dataSource: runtimeConfig.gameDataSource,
      sessionStore,
      page: logsPage.value,
      pageSize: logsPageSize.value,
      changeType: logsChangeType.value
    })
    inventoryLogs.value = result.items
    logsTotal.value = result.total
    logsPage.value = result.page
    logsPageSize.value = result.pageSize
  } catch (error) {
    logsError.value = error instanceof Error ? error.message : '流水加载失败'
  } finally {
    logsLoading.value = false
  }
}

function formatDelta(value: number) {
  return value >= 0 ? `+${value}` : `${value}`
}

async function confirmSell() {
  const pendingSellItem = getPendingSellItem()
  if (!pendingSellItem) {
    return
  }

  closeSellConfirmation()
  await mutateInventory('sell', pendingSellItem.itemId)
}

async function showInventoryLogs() {
  logsVisible.value = true
  logsPage.value = 1
  await refreshInventoryLogs()
}

async function changeLogFilter(changeType: InventoryLogChangeType) {
  if (logsChangeType.value === changeType && logsVisible.value) {
    return
  }

  logsChangeType.value = changeType
  logsPage.value = 1

  if (logsVisible.value) {
    await refreshInventoryLogs()
  }
}

async function goToPreviousLogsPage() {
  if (logsPage.value <= 1 || logsLoading.value) {
    return
  }

  logsPage.value -= 1
  await refreshInventoryLogs()
}

async function goToNextLogsPage() {
  if (logsPage.value >= logsTotalPages.value || logsLoading.value) {
    return
  }

  logsPage.value += 1
  await refreshInventoryLogs()
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
    <section
      v-if="getPendingSellItem()"
      class="sell-confirmation"
      data-testid="sell-confirmation"
    >
      <div>
        <strong>确认出售 {{ getPendingSellItem()?.name }}</strong>
        <p>预计获得 {{ getPendingSellItem()?.sellPrice }} 铜钱</p>
      </div>
      <div class="sell-confirmation__actions">
        <button
          type="button"
          class="action-button action-button--ghost"
          data-testid="sell-cancel"
          @click="closeSellConfirmation"
        >
          取消出售
        </button>
        <button
          type="button"
          class="action-button action-button--sell"
          data-testid="sell-confirm"
          @click="confirmSell"
        >
          确认出售
        </button>
      </div>
    </section>

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
            <p v-if="item.isHighValue" class="item-card__warning">高价值道具，默认锁定</p>
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
                :disabled="pendingItemId === item.itemId || isItemLocked(item)"
                @click="openSellConfirmation(item)"
              >
                出售1个
              </button>
              <button
                v-if="item.isHighValue"
                type="button"
                class="action-button action-button--lock"
                :data-testid="`inventory-lock-${item.name}`"
                @click="toggleItemLock(item)"
              >
                {{ isItemLocked(item) ? '解除锁定' : '重新锁定' }}
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
        <button
          type="button"
          class="ghost-button"
          data-testid="view-inventory-logs"
          @click="showInventoryLogs"
        >
          查看最近流水
        </button>
        <p v-if="logsLoading">流水加载中...</p>
        <p v-else-if="logsError" class="error-banner">{{ logsError }}</p>
        <div v-else-if="logsVisible" class="log-list">
          <h3 class="log-list__title">最近流水</h3>
          <div class="log-filter-row">
            <button
              v-for="option in logFilterOptions"
              :key="option.value"
              type="button"
              class="ghost-button log-filter-button"
              :data-testid="`inventory-log-filter-${option.value}`"
              :disabled="logsLoading"
              @click="changeLogFilter(option.value)"
            >
              {{ option.label }}
            </button>
          </div>
          <article v-for="log in inventoryLogs" :key="log.id" class="log-card">
            <strong>{{ log.title }}</strong>
            <span>{{ log.changeTypeText }}</span>
            <small>原因：{{ log.reason }}</small>
            <span>{{ log.delta }}</span>
            <small>铜钱 {{ formatDelta(log.coinsDelta) }} · 元宝 {{ formatDelta(log.diamondsDelta) }}</small>
            <small>{{ log.createdAtLabel }}</small>
          </article>
          <div class="log-pagination">
            <button
              type="button"
              class="ghost-button log-pager-button"
              data-testid="inventory-log-prev-page"
              :disabled="logsLoading || logsPage <= 1"
              @click="goToPreviousLogsPage"
            >
              上一页
            </button>
            <span data-testid="inventory-log-page">第 {{ logsPage }} / {{ logsTotalPages }} 页</span>
            <button
              type="button"
              class="ghost-button log-pager-button"
              data-testid="inventory-log-next-page"
              :disabled="logsLoading || logsPage >= logsTotalPages"
              @click="goToNextLogsPage"
            >
              下一页
            </button>
          </div>
        </div>
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

.sell-confirmation {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  padding: 14px 16px;
  border-radius: 16px;
  background: linear-gradient(135deg, #fff7ed, #fffbeb);
  border: 1px solid #fdba74;
}

.sell-confirmation__actions {
  display: flex;
  gap: 10px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.item-list,
.text-list,
.log-list {
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

.log-list {
  margin-top: 12px;
}

.log-filter-row,
.log-pagination {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.log-filter-button,
.log-pager-button {
  margin-top: 0;
}

.log-list__title {
  margin: 0;
  font-size: 16px;
  color: #1f2937;
}

.log-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px;
  border-radius: 14px;
  background: #f8fafc;
}

.item-card__warning {
  color: #b45309;
  font-size: 13px;
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

.action-button--ghost {
  background: #e2e8f0;
  color: #334155;
}

.action-button--lock {
  background: #fef3c7;
  color: #b45309;
}

.action-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .sell-confirmation,
  .sell-confirmation__actions {
    flex-direction: column;
    align-items: stretch;
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
