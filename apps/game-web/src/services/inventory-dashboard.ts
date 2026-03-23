import { legacyInventoryIconMap } from '@/assets/legacy'
import { request } from '@/api/http'
import { runtimeConfig, type GameDataSource } from '@/config/runtime'
import {
  createMockInventoryDashboard,
  createMockInventoryLogs,
  sellMockInventoryItem,
  useMockInventoryItem
} from '@/mocks/inventory-dashboard'

import type {
  InventoryDashboardData,
  InventoryItemCard,
  InventoryLogEntry,
  InventoryWalletResource,
  InventoryWalletStat
} from './inventory-dashboard.types'

interface SessionStoreLike {
  playerId: number
  token: string
  ensureGuestSession(): Promise<{ player_id: number; token: string }>
}

interface WalletResponse {
  player_id: number
  coins: number
  diamonds: number
  updated_at: string
}

interface InventoryResponseItem {
  player_id: number
  item_id: string
  item_name: string
  quantity: number
  sell_price: number
}

interface ResourceChangeLogResponse {
  player_id: number
  change_type: string
  biz_id: string
  coins_delta: number
  diamonds_delta: number
  reason: string
  created_at: string
}

interface LoadInventoryDashboardOptions {
  dataSource?: GameDataSource
  sessionStore: SessionStoreLike
}

interface InventoryOperateOptions extends LoadInventoryDashboardOptions {
  itemId: string
  count?: number
}

interface InventoryOperateResponse {
  wallet: WalletResponse
  item: InventoryResponseItem
}

interface InventoryOperateResult {
  dashboard: InventoryDashboardData
  message: string
}

interface InventoryCatalogMeta {
  name: string
  action: string
  highValue?: boolean
  icon?: string
}

const inventoryCatalog: Record<string, InventoryCatalogMeta> = {
  potion_small: {
    name: '中级经验丹',
    action: '可直接喂给主力幻兽',
    icon: legacyInventoryIconMap['中级经验丹']
  },
  summon_scroll: {
    name: '召唤卷轴',
    action: '可召唤随机幻兽',
    highValue: true
  },
  fire_evolution_stone: {
    name: '火系进化石',
    action: '烈焰狼王升境材料',
    icon: legacyInventoryIconMap['火系进化石']
  },
  strengthen_stone: {
    name: '强化石',
    action: '用于战骨与装备强化'
  }
}

function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value)
}

function adaptWallet(wallet: WalletResponse): InventoryWalletStat[] {
  return [
    { label: '铜钱', value: formatNumber(wallet.coins) },
    { label: '元宝', value: formatNumber(wallet.diamonds) },
    { label: '声望', value: '0' },
    { label: '活力', value: '120 / 120' }
  ]
}

function adaptWalletResources(wallet: WalletResponse): InventoryWalletResource[] {
  return [
    { label: '铜钱', value: formatNumber(wallet.coins), icon: legacyInventoryIconMap['铜钱'] },
    { label: '元宝', value: formatNumber(wallet.diamonds), icon: legacyInventoryIconMap['元宝'] }
  ]
}

function adaptInventoryItem(item: InventoryResponseItem): InventoryItemCard {
  const meta = inventoryCatalog[item.item_id]
  return {
    itemId: item.item_id,
    name: meta?.name ?? item.item_name,
    count: `x${item.quantity}`,
    action: meta?.action ?? '暂无使用说明',
    sellPrice: item.sell_price,
    isHighValue: meta?.highValue ?? item.sell_price >= 50,
    icon: meta?.icon
  }
}

function createApiInventoryDashboard(wallet: WalletResponse, inventory: InventoryResponseItem[]): InventoryDashboardData {
  const dashboard = createMockInventoryDashboard()
  const visibleInventory = inventory.filter((item) => item.quantity > 0)

  dashboard.hero = {
    ...dashboard.hero,
    metaValue: `${visibleInventory.length} / 30`
  }
  dashboard.wallet = adaptWallet(wallet)
  dashboard.walletResources = adaptWalletResources(wallet)
  dashboard.items = visibleInventory.map(adaptInventoryItem)

  return dashboard
}

function getInventoryDisplayName(itemId: string): string {
  return inventoryCatalog[itemId]?.name ?? itemId
}

function formatLogTitle(log: ResourceChangeLogResponse): string {
  switch (log.change_type) {
    case 'inventory_sell':
      return `出售${getInventoryDisplayName(log.biz_id)}`
    case 'inventory_use':
      return `使用${getInventoryDisplayName(log.biz_id)}`
    case 'grant_reward':
      return '发放奖励'
    default:
      return log.biz_id
  }
}

function formatLogDelta(log: ResourceChangeLogResponse): string {
  if (log.coins_delta !== 0) {
    return `${log.coins_delta > 0 ? '+' : ''}${formatNumber(log.coins_delta)} 铜钱`
  }
  if (log.diamonds_delta !== 0) {
    return `${log.diamonds_delta > 0 ? '+' : ''}${formatNumber(log.diamonds_delta)} 元宝`
  }
  return '无货币变化'
}

function adaptInventoryLogs(logs: ResourceChangeLogResponse[]): InventoryLogEntry[] {
  return [...logs]
    .reverse()
    .map((log) => ({
      title: formatLogTitle(log),
      delta: formatLogDelta(log),
      createdAtLabel: log.created_at.replace('T', ' ').replace('Z', '')
    }))
}

export async function loadInventoryDashboard({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore
}: LoadInventoryDashboardOptions): Promise<InventoryDashboardData> {
  if (dataSource === 'mock') {
    return createMockInventoryDashboard()
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = session.token ? { Authorization: `Bearer ${session.token}` } : undefined

  const [wallet, inventory] = await Promise.all([
    request<WalletResponse>(`/player/assets/wallet?player_id=${session.player_id}`, { headers }),
    request<InventoryResponseItem[]>(`/player/assets/inventory?player_id=${session.player_id}`, { headers })
  ])

  return createApiInventoryDashboard(wallet, inventory)
}

export async function loadInventoryLogs({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore
}: LoadInventoryDashboardOptions): Promise<InventoryLogEntry[]> {
  if (dataSource === 'mock') {
    return createMockInventoryLogs()
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = session.token ? { Authorization: `Bearer ${session.token}` } : undefined
  const logs = await request<ResourceChangeLogResponse[]>(
    `/player/assets/logs?player_id=${session.player_id}`,
    { headers }
  )

  return adaptInventoryLogs(logs)
}

export async function useInventoryItem({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore,
  itemId,
  count = 1
}: InventoryOperateOptions): Promise<InventoryOperateResult> {
  if (dataSource === 'mock') {
    useMockInventoryItem(itemId, count)
    return {
      dashboard: createMockInventoryDashboard(),
      message: `已使用 ${count} 个${getInventoryDisplayName(itemId)}`
    }
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = session.token ? { Authorization: `Bearer ${session.token}` } : undefined

  await request<InventoryOperateResponse>('/player/assets/inventory/use', {
    method: 'POST',
    headers,
    body: JSON.stringify({
      player_id: session.player_id,
      item_id: itemId,
      count
    })
  })

  return {
    dashboard: await loadInventoryDashboard({ dataSource: 'api', sessionStore }),
    message: `已使用 ${count} 个${getInventoryDisplayName(itemId)}`
  }
}

export async function sellInventoryItem({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore,
  itemId,
  count = 1
}: InventoryOperateOptions): Promise<InventoryOperateResult> {
  if (dataSource === 'mock') {
    sellMockInventoryItem(itemId, count)
    return {
      dashboard: createMockInventoryDashboard(),
      message: `已出售 ${count} 个${getInventoryDisplayName(itemId)}`
    }
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = session.token ? { Authorization: `Bearer ${session.token}` } : undefined

  await request<InventoryOperateResponse>('/player/assets/inventory/sell', {
    method: 'POST',
    headers,
    body: JSON.stringify({
      player_id: session.player_id,
      item_id: itemId,
      count
    })
  })

  return {
    dashboard: await loadInventoryDashboard({ dataSource: 'api', sessionStore }),
    message: `已出售 ${count} 个${getInventoryDisplayName(itemId)}`
  }
}
