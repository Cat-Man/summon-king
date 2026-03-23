import { afterEach, expect, test, vi } from 'vitest'

import {
  loadPetCatalogDashboard,
  loadPetDetailDashboard,
  loadPetTeamDashboard,
  savePetTeamSelection
} from '../pet-dashboard'

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
})

test('loadPetCatalogDashboard should map backend pet catalog in api mode', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue({
      player_id: 3001,
      token: 'guest-token-3001'
    })
  }

  const fetchMock = vi.fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          { pet_id: 1, name: '赤焰狼', rarity: 2, element: 'fire', locked: false, owned: true, power: 120, portrait: 'pet_1.png', story_text: '火山边缘的狩猎者' },
          { pet_id: 2, name: '寒霜鹤', rarity: 3, element: 'water', locked: false, owned: true, power: 150, portrait: 'pet_2.png', story_text: '来自北境冰湖的守望者' },
          { pet_id: 3, name: '雷角牛', rarity: 4, element: 'thunder', locked: false, owned: true, power: 180, portrait: 'pet_3.png', story_text: '暴风中的冲锋者' }
        ],
        trace_id: 'trace-pet-catalog'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          { player_id: 3001, pet_id: 1, level: 12, star: 1, power: 120 },
          { player_id: 3001, pet_id: 2, level: 16, star: 2, power: 150 },
          { player_id: 3001, pet_id: 3, level: 20, star: 3, power: 180 }
        ],
        trace_id: 'trace-pet-list'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const result = await loadPetCatalogDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(sessionStore.ensureGuestSession).toHaveBeenCalledTimes(1)
  expect(result.hero.title).toBe('幻兽图鉴')
  expect(result.hero.metaValue).toBe('3 / 3')
  expect(result.overview.find((item) => item.label === '已拥有')?.value).toBe('3')
  expect(result.overview.find((item) => item.label === '队伍推荐')?.value).toBe('雷角牛')
  expect(result.pets[0]).toMatchObject({
    petId: 1,
    name: '烈焰狼王',
    status: '已拥有',
    level: 'Lv.12'
  })
})

test('loadPetDetailDashboard should merge backend detail and player pet state', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue({
      player_id: 3001,
      token: 'guest-token-3001'
    })
  }

  const fetchMock = vi.fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          { player_id: 3001, pet_id: 1, level: 12, star: 1, power: 120 },
          { player_id: 3001, pet_id: 2, level: 16, star: 2, power: 150 }
        ],
        trace_id: 'trace-pet-list'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: { pet_id: 2, name: '寒霜鹤', rarity: 3, element: 'water', locked: false, owned: true, power: 150, portrait: 'pet_2.png', story_text: '来自北境冰湖的守望者' },
        trace_id: 'trace-pet-detail'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const result = await loadPetDetailDashboard({
    dataSource: 'api',
    sessionStore,
    petId: 2
  })

  expect(result.hero.title).toBe('寒枝鹿灵')
  expect(result.hero.metaValue).toBe('150')
  expect(result.overview.find((item) => item.label === '等级')?.value).toBe('Lv.16')
  expect(result.overview.find((item) => item.label === '星级')?.value).toBe('2 星')
  expect(result.skills[0]).toContain('寒枝鹿灵')
  expect(result.sources[0]).toContain('北境冰湖')
})

test('loadPetTeamDashboard should read team and savePetTeamSelection should persist it in api mode', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue({
      player_id: 3001,
      token: 'guest-token-3001'
    })
  }

  const fetchMock = vi.fn()
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: [
          { player_id: 3001, pet_id: 1, level: 12, star: 1, power: 120 },
          { player_id: 3001, pet_id: 2, level: 16, star: 2, power: 150 },
          { player_id: 3001, pet_id: 3, level: 20, star: 3, power: 180 }
        ],
        trace_id: 'trace-pet-list'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: { player_id: 3001, pet_ids: [3, 2] },
        trace_id: 'trace-pet-team'
      })
    })
    .mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: { saved: true },
        trace_id: 'trace-pet-team-save'
      })
    })

  vi.stubGlobal('fetch', fetchMock)

  const dashboard = await loadPetTeamDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(dashboard.hero.metaValue).toBe('330')
  expect(dashboard.team.map((item) => item.petId)).toEqual([3, 2])
  expect(dashboard.roster.find((item) => item.petId === 3)?.status).toBe('已上阵')

  const result = await savePetTeamSelection({
    dataSource: 'api',
    sessionStore,
    petIds: [2, 1]
  })

  expect(result.message).toBe('阵容已保存')
  expect(fetchMock).toHaveBeenLastCalledWith(
    '/api/v1/player/pets/team/save',
    expect.objectContaining({
      method: 'POST',
      body: JSON.stringify({ player_id: 3001, pet_ids: [2, 1] }),
      headers: expect.objectContaining({ Authorization: 'Bearer guest-token-3001' })
    })
  )
})
