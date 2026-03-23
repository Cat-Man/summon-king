import { afterEach, expect, test, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { legacyInventoryIconMap } from '@/assets/legacy'
import { resetMockInventoryDashboardState } from '@/mocks/inventory-dashboard'
import { useSessionStore } from '@/stores/session'

import {
  loadInventoryDashboard,
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
