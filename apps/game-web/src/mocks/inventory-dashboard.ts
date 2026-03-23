import { legacyInventoryIconMap } from '@/assets/legacy'

import type {
  InventoryDashboardData,
  InventoryLogEntry
} from '@/services/inventory-dashboard.types'

interface MockInventoryItemState {
  itemId: string
  quantity: number
  sellPrice: number
}

interface MockInventoryState {
  coins: number
  diamonds: number
  items: MockInventoryItemState[]
  logs: MockResourceLogState[]
}

interface MockResourceLogState {
  changeType: string
  bizId: string
  reason: string
  coinsDelta: number
  diamondsDelta: number
  createdAt: string
}

interface MockCatalogMeta {
  name: string
  action: string
  highValue?: boolean
  icon?: string
}

const mockCatalog: Record<string, MockCatalogMeta> = {
  potion_small: {
    name: '中级经验丹',
    action: '可直接喂给主力幻兽',
    icon: legacyInventoryIconMap['中级经验丹']
  },
  fire_evolution_stone: {
    name: '火系进化石',
    action: '烈焰狼王升境材料',
    icon: legacyInventoryIconMap['火系进化石']
  },
  strengthen_stone: {
    name: '强化石',
    action: '用于战骨与装备强化'
  },
  summon_scroll: {
    name: '召唤卷轴',
    action: '可召唤随机幻兽',
    highValue: true
  }
}

function createInitialMockInventoryState(): MockInventoryState {
  return {
    coins: 0,
    diamonds: 0,
    items: [
      { itemId: 'potion_small', quantity: 12, sellPrice: 10 },
      { itemId: 'fire_evolution_stone', quantity: 6, sellPrice: 20 },
      { itemId: 'strengthen_stone', quantity: 48, sellPrice: 5 },
      { itemId: 'summon_scroll', quantity: 1, sellPrice: 50 }
    ],
    logs: [
      {
        changeType: 'inventory_sell',
        bizId: 'summon_scroll',
        reason: 'inventory_sell',
        coinsDelta: 50,
        diamondsDelta: 0,
        createdAt: '2026-03-23 09:15'
      },
      {
        changeType: 'inventory_use',
        bizId: 'potion_small',
        reason: 'inventory_use',
        coinsDelta: 0,
        diamondsDelta: 0,
        createdAt: '2026-03-23 09:10'
      },
      {
        changeType: 'grant_reward',
        bizId: 'signin-1',
        reason: 'daily_signin',
        coinsDelta: 100,
        diamondsDelta: 0,
        createdAt: '2026-03-23 09:00'
      },
      {
        changeType: 'grant_reward',
        bizId: 'mail-1',
        reason: 'mail_reward',
        coinsDelta: 0,
        diamondsDelta: 10,
        createdAt: '2026-03-22 21:20'
      },
      {
        changeType: 'inventory_use',
        bizId: 'strengthen_stone',
        reason: 'inventory_use',
        coinsDelta: 0,
        diamondsDelta: 0,
        createdAt: '2026-03-22 18:45'
      }
    ]
  }
}

let mockInventoryState = createInitialMockInventoryState()

function getMockCatalogMeta(itemId: string): MockCatalogMeta {
  return mockCatalog[itemId] ?? { name: itemId, action: '暂无使用说明' }
}

export function resetMockInventoryDashboardState() {
  mockInventoryState = createInitialMockInventoryState()
}

export function useMockInventoryItem(itemId: string, count = 1) {
  const item = mockInventoryState.items.find((entry) => entry.itemId === itemId)
  if (!item || item.quantity < count) {
    throw new Error('inventory not enough')
  }

  item.quantity -= count
  mockInventoryState.logs.unshift({
    changeType: 'inventory_use',
    bizId: itemId,
    reason: 'inventory_use',
    coinsDelta: 0,
    diamondsDelta: 0,
    createdAt: '2026-03-23 09:10'
  })
}

export function sellMockInventoryItem(itemId: string, count = 1) {
  const item = mockInventoryState.items.find((entry) => entry.itemId === itemId)
  if (!item || item.quantity < count) {
    throw new Error('inventory not enough')
  }

  item.quantity -= count
  mockInventoryState.coins += item.sellPrice * count
  mockInventoryState.logs.unshift({
    changeType: 'inventory_sell',
    bizId: itemId,
    reason: 'inventory_sell',
    coinsDelta: item.sellPrice * count,
    diamondsDelta: 0,
    createdAt: '2026-03-23 09:15'
  })
}

function formatMockLogTitle(changeType: string, bizId: string): string {
  switch (changeType) {
    case 'inventory_sell':
      return `出售${getMockCatalogMeta(bizId).name}`
    case 'inventory_use':
      return `使用${getMockCatalogMeta(bizId).name}`
    case 'grant_reward':
      return '发放奖励'
    default:
      return bizId
  }
}

function formatMockLogChangeTypeText(changeType: string): string {
  switch (changeType) {
    case 'inventory_sell':
      return '出售道具'
    case 'inventory_use':
      return '使用道具'
    case 'grant_reward':
      return '发放奖励'
    default:
      return changeType
  }
}

function formatMockLogDelta(log: MockResourceLogState): string {
  if (log.coinsDelta !== 0) {
    return `${log.coinsDelta > 0 ? '+' : ''}${log.coinsDelta} 铜钱`
  }
  if (log.diamondsDelta !== 0) {
    return `${log.diamondsDelta > 0 ? '+' : ''}${log.diamondsDelta} 元宝`
  }
  return '无货币变化'
}

export function createMockInventoryLogs(limit = 5): InventoryLogEntry[] {
  return mockInventoryState.logs
    .slice()
    .sort((left, right) => right.createdAt.localeCompare(left.createdAt))
    .slice(0, limit)
    .map((log) => ({
      id: `${log.changeType}-${log.bizId}-${log.createdAt}`,
      createdAt: log.createdAt,
      changeType: log.changeType,
      changeTypeText: formatMockLogChangeTypeText(log.changeType),
      reason: log.reason,
      coinsDelta: log.coinsDelta,
      diamondsDelta: log.diamondsDelta,
      title: formatMockLogTitle(log.changeType, log.bizId),
      delta: formatMockLogDelta(log),
      createdAtLabel: log.createdAt
    }))
}

export function createMockInventoryDashboard(): InventoryDashboardData {
  const visibleItems = mockInventoryState.items.filter((item) => item.quantity > 0)

  return {
    hero: {
      eyebrow: '资源与道具管理',
      title: '背包',
      description: '把资源钱包、普通背包与资源流水入口放在同一页，方便玩家快速判断当前是否缺钱、缺材料、缺活力。',
      tone: 'amber',
      metaLabel: '背包容量',
      metaValue: `${visibleItems.length} / 30`
    },
    wallet: [
      { label: '铜钱', value: String(mockInventoryState.coins) },
      { label: '元宝', value: String(mockInventoryState.diamonds) },
      { label: '声望', value: '0' },
      { label: '活力', value: '120 / 120' }
    ],
    walletResources: [
      { label: '铜钱', value: String(mockInventoryState.coins), icon: legacyInventoryIconMap['铜钱'] },
      { label: '元宝', value: String(mockInventoryState.diamonds), icon: legacyInventoryIconMap['元宝'] }
    ],
    items: visibleItems.map((item) => {
      const meta = getMockCatalogMeta(item.itemId)

      return {
        itemId: item.itemId,
        name: meta.name,
        count: `x${item.quantity}`,
        action: meta.action,
        sellPrice: item.sellPrice,
        isHighValue: Boolean(meta.highValue),
        icon: meta.icon
      }
    }),
    rules: [
      '使用道具前先校验目标幻兽与当前玩法是否匹配',
      '出售前需展示预计获得的铜钱并二次确认',
      '高价值道具默认锁定，避免误操作批量出售'
    ]
  }
}
