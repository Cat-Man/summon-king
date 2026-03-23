import { afterEach, expect, test, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { createMockHomeDashboard } from '@/mocks/home-dashboard'
import { useSessionStore } from '@/stores/session'
import { loadHomeDashboard } from '../home-dashboard'

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

test('loadHomeDashboard should return home view model in mock mode', async () => {
  setActivePinia(createPinia())
  const sessionStore = useSessionStore()

  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
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
  )

  const result = await loadHomeDashboard({
    dataSource: 'mock',
    sessionStore
  })

  expect(result.hero.title).toContain('召唤之王')
  expect(result.resources.find((item) => item.label === '等级')?.value).toBe('Lv.1')
  expect(result.resourceIcons.find((item) => item.label === '铜钱')?.value).toBe('0')
})

test('loadHomeDashboard should hydrate dynamic profile in api mode', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue(undefined),
    fetchHomeIndex: vi.fn().mockResolvedValue({
      player_id: 2002,
      nickname: '联调游客2002',
      level: 7,
      coin: 1280,
      diamond: 96,
      last_login_at: '2026-03-23T10:00:00Z'
    })
  }

  const result = await loadHomeDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(sessionStore.ensureGuestSession).toHaveBeenCalledTimes(1)
  expect(sessionStore.fetchHomeIndex).toHaveBeenCalledTimes(1)
  expect(result.hero.metaValue).toBe('联调游客2002')
  expect(result.resources.find((item) => item.label === '等级')?.value).toBe('Lv.7')
  expect(result.resourceIcons.find((item) => item.label === '铜钱')?.value).toBe('1,280')
})
