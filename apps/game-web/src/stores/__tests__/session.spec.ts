import { afterEach, expect, test, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useSessionStore } from '../session'

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

test('ensureGuestSession should login guest and persist session when empty', async () => {
  setActivePinia(createPinia())
  const store = useSessionStore()

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

  const ensureGuestSession = (store as unknown as { ensureGuestSession: () => Promise<void> })
    .ensureGuestSession

  await ensureGuestSession()

  expect(store.playerId).toBe(1001)
  expect(store.token).toBe('guest-token-1001')
})

test('fetchHomeIndex should return parsed player home summary', async () => {
  setActivePinia(createPinia())
  const store = useSessionStore()
  store.setSession(1001, 'guest-token-1001')

  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          player_id: 1001,
          nickname: '游客1001',
          level: 1,
          coin: 0,
          diamond: 0,
          last_login_at: '2026-03-23T09:00:00Z'
        },
        trace_id: 'trace-home-1'
      })
    })
  )

  const homeIndex = await store.fetchHomeIndex()

  expect(homeIndex.nickname).toBe('游客1001')
  expect(homeIndex.level).toBe(1)
  expect(homeIndex.coin).toBe(0)
  expect(homeIndex.diamond).toBe(0)
})
