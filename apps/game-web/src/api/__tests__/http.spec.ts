import { afterEach, expect, test, vi } from 'vitest'

import { request } from '../http'

afterEach(() => {
  vi.restoreAllMocks()
})

test('unwraps backend api response data', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: { player_id: 1001 },
        trace_id: 'trace-1'
      })
    })
  )

  const data = await request<{ player_id: number }>('/api/v1/player/home/index')

  expect(data.player_id).toBe(1001)
})

test('prefixes default api base url for relative endpoints', async () => {
  const fetchMock = vi.fn().mockResolvedValue({
    ok: true,
    json: async () => ({
      code: 0,
      message: 'ok',
      data: { ok: true },
      trace_id: 'trace-2'
    })
  })
  vi.stubGlobal('fetch', fetchMock)

  await request<{ ok: boolean }>('/player/home/index')

  expect(fetchMock).toHaveBeenCalledWith(
    '/api/v1/player/home/index',
    expect.objectContaining({
      headers: expect.objectContaining({
        'Content-Type': 'application/json'
      })
    })
  )
})
