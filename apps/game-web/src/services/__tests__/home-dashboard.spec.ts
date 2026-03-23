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
      last_login_at: '2026-03-23T10:00:00Z',
      daily_todos: [{ title: '签到状态', value: '今日已签', action: '查看奖励' }],
      resources: [
        { label: '等级', value: 'Lv.7' },
        { label: '战力', value: '1,680' },
        { label: '活力', value: '120 / 120' },
        { label: '铜钱', value: '1,280' },
        { label: '元宝', value: '96' },
        { label: '声望', value: '30' }
      ],
      resource_icons: [
        { label: '铜钱', value: '1,280' },
        { label: '元宝', value: '96' }
      ],
      messages: ['系统消息：VIP 每日宝箱可领取'],
      entries: ['世界地图', '联盟', '背包'],
      activity_entry: {
        title: '修行收益',
        description: '火焰山修行进行中',
        note: '剩余 1小时59分',
        action: '查看修行'
      },
      cultivation_summary: {
        status: 'running',
        map_name: '火焰山',
        started_at: '2026-03-23T08:00:00Z',
        finished_at: '2026-03-23T10:00:00Z',
        remaining_seconds: 7140,
        reward_coins: 240,
        reward_pet_exp: 160,
        action: '查看修行'
      }
    })
  }

  const result = await loadHomeDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(sessionStore.ensureGuestSession).toHaveBeenCalledTimes(1)
  expect(sessionStore.fetchHomeIndex).toHaveBeenCalledTimes(1)
  expect(result.hero.metaValue).toBe('联调游客2002')
  expect(result.dailyTodos[0]).toEqual({
    title: '签到状态',
    value: '今日已签',
    action: '查看奖励'
  })
  expect(result.resources.find((item) => item.label === '等级')?.value).toBe('Lv.7')
  expect(result.resources.find((item) => item.label === '战力')?.value).toBe('1,680')
  expect(result.resourceIcons.find((item) => item.label === '铜钱')?.value).toBe('1,280')
  expect(result.messages).toEqual(['系统消息：VIP 每日宝箱可领取'])
  expect(result.entries).toEqual(['世界地图', '联盟', '背包'])
  expect(result.activityEntry.title).toBe('修行收益')
  expect(result.activityEntry.description).toContain('火焰山')
  expect(result.cultivationSummary.action).toBe('查看修行')
  expect(result.cultivationSummary.items).toEqual([
    { label: '修行地图', value: '火焰山', subtext: '进行中' },
    { label: '结束时间', value: '2026-03-23 10:00', subtext: '剩余 1小时59分' },
    { label: '铜钱收益', value: '240' },
    { label: '幻兽经验', value: '160' }
  ])
})

test('loadHomeDashboard should map idle cultivation summary to user-facing copy', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue(undefined),
    fetchHomeIndex: vi.fn().mockResolvedValue({
      player_id: 2003,
      nickname: '联调游客2003',
      level: 1,
      coin: 0,
      diamond: 0,
      last_login_at: '2026-03-23T10:00:00Z',
      daily_todos: [{ title: '签到状态', value: '今日未签', action: '前往签到' }],
      resources: [
        { label: '等级', value: 'Lv.1' },
        { label: '战力', value: '450' },
        { label: '活力', value: '120 / 120' },
        { label: '铜钱', value: '0' },
        { label: '元宝', value: '0' },
        { label: '声望', value: '50' }
      ],
      resource_icons: [
        { label: '铜钱', value: '0' },
        { label: '元宝', value: '0' }
      ],
      messages: ['修行消息：当前暂无进行中的修行'],
      entries: ['世界地图', '联盟', '背包'],
      activity_entry: {
        title: '签到奖励',
        description: '今日签到待领取',
        note: '前往签到',
        action: '前往签到'
      },
      cultivation_summary: {
        status: 'idle',
        map_name: '青云城',
        started_at: '0001-01-01T00:00:00Z',
        finished_at: '0001-01-01T00:00:00Z',
        remaining_seconds: 0,
        reward_coins: 0,
        reward_pet_exp: 0,
        action: '前往修行'
      }
    })
  }

  const result = await loadHomeDashboard({
    dataSource: 'api',
    sessionStore
  })

  expect(result.cultivationSummary.action).toBe('前往修行')
  expect(result.cultivationSummary.items).toEqual([
    { label: '修行地图', value: '青云城', subtext: '待开始' },
    { label: '结束时间', value: '未开始', subtext: '可立即开启' },
    { label: '铜钱收益', value: '0' },
    { label: '幻兽经验', value: '0' }
  ])
})
