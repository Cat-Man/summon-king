import { afterEach, expect, test, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { legacyInventoryIconMap } from '@/assets/legacy'
import { resetMockInventoryDashboardState } from '@/mocks/inventory-dashboard'
import { useSessionStore } from '@/stores/session'

import {
  loadInventoryDashboard,
  loadRecentAssetLogs,
  sellInventoryItem,
  useInventoryItem
} from '../inventory-dashboard'

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
  resetMockInventoryDashboardState()
})

test('loadInventoryDashboard should return inventory view model in mock mode', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  const result = await loadInventoryDashboard({
    dataSource: 'mock',
    sessionStore
  })

  expect(result.hero.metaValue).toBe('4 / 30')
  expect(result.wallet.find((item) => item.label === '铜钱')?.value).toBe('0')
  expect(result.items.find((item) => item.name === '火系进化石')?.icon).toBe(
    legacyInventoryIconMap['火系进化石']
  )
})

test('loadInventoryDashboard should map backend wallet and inventory in api mode', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  const fetchMock = vi
    .fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          player_id: 1001,
          token: 'guest-token-1001',
          channel: 'web'
        },
        trace_id: 'trace-login-1'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          player_id: 1001,
          coins: 0,
          diamonds: 0,
          updated_at: '2026-03-23T09:00:00Z'
        },
        trace_id: 'trace-wallet-1'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          {
            player_id: 1001,
            item_id: 'potion_small',
            item_name: 'Small Potion',
            quantity: 5,
            sell_price: 10
          },
          {
            player_id: 1001,
            item_id: 'summon_scroll',
            item_name: 'Summon Scroll',
            quantity: 1,
            sell_price: 50
          }
        ],
        trace_id: 'trace-inventory-1'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const result = await loadInventoryDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(fetchMock).toHaveBeenCalledTimes(3)
  expect(result.hero.metaValue).toBe('2 / 30')
  expect(result.wallet.find((item) => item.label === '铜钱')?.value).toBe('0')
  expect(result.items.find((item) => item.name === '中级经验丹')?.count).toBe('x5')
  expect(result.items.find((item) => item.name === '中级经验丹')?.icon).toBe(
    legacyInventoryIconMap['中级经验丹']
  )
  expect(result.items.find((item) => item.name === '召唤卷轴')?.count).toBe('x1')
})

test('useInventoryItem should update mock dashboard and decrement item quantity', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  const result = await useInventoryItem({
    dataSource: 'mock',
    sessionStore,
    itemId: 'potion_small'
  })

  expect(result.message).toContain('中级经验丹')
  expect(result.dashboard.items.find((item) => item.name === '中级经验丹')?.count).toBe('x11')
})

test('loadRecentAssetLogs should return paged logs in mock mode', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  const result = await loadRecentAssetLogs({
    dataSource: 'mock',
    sessionStore,
    page: 1,
    pageSize: 2,
    changeType: 'all'
  })

  expect(result.total).toBe(5)
  expect(result.page).toBe(1)
  expect(result.pageSize).toBe(2)
  expect(result.items).toHaveLength(2)
  expect(result.items[0]?.changeType).toBe('inventory_sell')
  expect(result.items[1]?.changeType).toBe('inventory_use')
})

test('loadRecentAssetLogs should support changeType filter in mock mode', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  const result = await loadRecentAssetLogs({
    dataSource: 'mock',
    sessionStore,
    page: 1,
    pageSize: 5,
    changeType: 'inventory_use'
  })

  expect(result.total).toBe(2)
  expect(result.page).toBe(1)
  expect(result.pageSize).toBe(5)
  expect(result.items).toHaveLength(2)
  expect(result.items.every((entry) => entry.changeType === 'inventory_use')).toBe(true)
  expect(result.items.map((entry) => entry.changeTypeText)).toEqual(['使用道具', '使用道具'])
})

test('loadRecentAssetLogs should map backend resource logs for assets page', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()
  sessionStore.setSession(1001, 'guest-token-1001')

  const fetchMock = vi
    .fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          items: [
            {
              player_id: 1001,
              change_type: 'grant_reward',
              biz_id: 'signin-1',
              coins_delta: 100,
              diamonds_delta: 0,
              reason: 'daily_signin',
              created_at: '2026-03-23T09:00:00Z'
            },
            {
              player_id: 1001,
              change_type: 'inventory_sell',
              biz_id: 'summon_scroll',
              coins_delta: 50,
              diamonds_delta: 0,
              reason: 'inventory_sell',
              created_at: '2026-03-23T09:00:01Z'
            }
          ],
          total: 7,
          page: 2,
          page_size: 3
        },
        trace_id: 'trace-logs-1'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const result = await loadRecentAssetLogs({
    dataSource: 'api',
    sessionStore,
    page: 2,
    pageSize: 3,
    changeType: 'inventory_sell'
  })

  expect(fetchMock).toHaveBeenCalledTimes(1)
  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/player/assets/logs?player_id=1001&page=2&page_size=3&change_type=inventory_sell',
    expect.objectContaining({
      headers: expect.objectContaining({
        Authorization: 'Bearer guest-token-1001'
      })
    })
  )
  expect(result.total).toBe(7)
  expect(result.page).toBe(2)
  expect(result.pageSize).toBe(3)
  expect(result.items[0]?.changeType).toBe('inventory_sell')
  expect(result.items[0]?.changeTypeText).toBe('出售道具')
  expect(result.items[0]?.reason).toBe('inventory_sell')
  expect(result.items[0]?.coinsDelta).toBe(50)
  expect(result.items[0]?.diamondsDelta).toBe(0)
  expect(result.items[0]?.title).toBe('出售召唤卷轴')
  expect(result.items[0]?.delta).toBe('+50 铜钱')
  expect(result.items[1]?.title).toBe('发放奖励')
})

test('loadRecentAssetLogs should not send all sentinel to backend filter', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()
  sessionStore.setSession(1001, 'guest-token-1001')

  const fetchMock = vi.fn().mockResolvedValueOnce({
    ok: true,
    json: async () => ({
      code: 0,
      message: 'ok',
      data: {
        items: [],
        total: 0,
        page: 1,
        page_size: 2
      },
      trace_id: 'trace-logs-all-1'
    })
  })

  vi.stubGlobal('fetch', fetchMock)

  await loadRecentAssetLogs({
    dataSource: 'api',
    sessionStore,
    page: 1,
    pageSize: 2,
    changeType: 'all'
  })

  expect(fetchMock).toHaveBeenCalledTimes(1)
  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/player/assets/logs?player_id=1001&page=1&page_size=2&change_type=',
    expect.objectContaining({
      headers: expect.objectContaining({
        Authorization: 'Bearer guest-token-1001'
      })
    })
  )
})

test('sellInventoryItem should post to backend and refresh wallet and inventory dashboard', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()
  sessionStore.setSession(1001, 'guest-token-1001')

  const fetchMock = vi
    .fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          wallet: {
            player_id: 1001,
            coins: 50,
            diamonds: 0,
            updated_at: '2026-03-23T09:00:00Z'
          },
          item: {
            player_id: 1001,
            item_id: 'summon_scroll',
            item_name: 'Summon Scroll',
            quantity: 0,
            sell_price: 50
          }
        },
        trace_id: 'trace-sell-1'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          player_id: 1001,
          coins: 50,
          diamonds: 0,
          updated_at: '2026-03-23T09:00:01Z'
        },
        trace_id: 'trace-wallet-2'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          {
            player_id: 1001,
            item_id: 'potion_small',
            item_name: 'Small Potion',
            quantity: 5,
            sell_price: 10
          }
        ],
        trace_id: 'trace-inventory-2'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const result = await sellInventoryItem({
    dataSource: 'api',
    sessionStore,
    itemId: 'summon_scroll'
  })

  expect(fetchMock).toHaveBeenCalledTimes(3)
  expect(result.message).toContain('召唤卷轴')
  expect(result.dashboard.wallet.find((item) => item.label === '铜钱')?.value).toBe('50')
  expect(result.dashboard.items.find((item) => item.name === '召唤卷轴')).toBeUndefined()
})
