import { legacyInventoryIconMap } from '@/assets/legacy'

import type { InventoryDashboardData } from '@/services/inventory-dashboard.types'

interface MockInventoryItemState {
  itemId: string
  quantity: number
  sellPrice: number
}

interface MockInventoryState {
  coins: number
  diamonds: number
  items: MockInventoryItemState[]
}

interface MockCatalogMeta {
  name: string
  action: string
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
    action: '可召唤随机幻兽'
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
}

export function sellMockInventoryItem(itemId: string, count = 1) {
  const item = mockInventoryState.items.find((entry) => entry.itemId === itemId)
  if (!item || item.quantity < count) {
    throw new Error('inventory not enough')
  }

  item.quantity -= count
  mockInventoryState.coins += item.sellPrice * count
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
