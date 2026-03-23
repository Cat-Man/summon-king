import { runtimeConfig, type GameDataSource } from '@/config/runtime'
import { legacyActivityAssets, legacyResourceIconMap } from '@/assets/legacy'
import { createMockHomeDashboard } from '@/mocks/home-dashboard'
import type { HomeIndexResponse } from '@/stores/session'

import type { HomeDashboardData } from './home-dashboard.types'

interface SessionStoreLike {
  ensureGuestSession(): Promise<unknown>
  fetchHomeIndex(): Promise<HomeIndexResponse>
}

interface LoadHomeDashboardOptions {
  dataSource?: GameDataSource
  sessionStore: SessionStoreLike
}

function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value)
}

export function createInitialHomeDashboard(): HomeDashboardData {
  return {
    hero: {
      eyebrow: '今日工作台',
      title: '召唤之王',
      description: '',
      tone: 'navy',
      metaLabel: '当前角色',
      metaValue: ''
    },
    dailyTodos: [],
    resources: [
      { label: '等级', value: '' },
      { label: '战力', value: '' },
      { label: '活力', value: '' },
      { label: '铜钱', value: '' },
      { label: '元宝', value: '' },
      { label: '声望', value: '' }
    ],
    resourceIcons: [],
    activityEntry: {
      title: '',
      description: '',
      note: '',
      icon: legacyActivityAssets.activitySparkIcon
    },
    entries: [],
    messages: []
  }
}

function getResponseResourceValue(homeIndex: HomeIndexResponse, label: string): string | undefined {
  return homeIndex.resources.find((item) => item.label === label)?.value
}

function getResponseIconValue(homeIndex: HomeIndexResponse, label: string): string | undefined {
  return homeIndex.resource_icons.find((item) => item.label === label)?.value
}

function createApiHomeDashboard(homeIndex: HomeIndexResponse): HomeDashboardData {
  const dashboard = createMockHomeDashboard()

  dashboard.hero = {
    ...dashboard.hero,
    description: `欢迎回来，${homeIndex.nickname}。第一屏直接告诉玩家今天能做什么、有哪些收益可领、当前主推进目标是什么。`,
    metaValue: homeIndex.nickname
  }

  if (homeIndex.daily_todos.length > 0) {
    dashboard.dailyTodos = homeIndex.daily_todos.map((item) => ({ ...item }))
  }

  dashboard.resources = dashboard.resources.map((item) => {
    const responseValue = getResponseResourceValue(homeIndex, item.label)
    if (responseValue) {
      return { ...item, value: responseValue }
    }

    if (item.label === '等级') {
      return { ...item, value: `Lv.${homeIndex.level}` }
    }

    if (item.label === '铜钱') {
      return { ...item, value: formatNumber(homeIndex.coin) }
    }

    if (item.label === '元宝') {
      return { ...item, value: formatNumber(homeIndex.diamond) }
    }

    return item
  })

  dashboard.resourceIcons = dashboard.resourceIcons.map((item) => {
    const responseValue = getResponseIconValue(homeIndex, item.label)
    if (responseValue) {
      return { ...item, value: responseValue }
    }

    if (item.label === '铜钱') {
      return { ...item, value: formatNumber(homeIndex.coin) }
    }

    if (item.label === '元宝') {
      return { ...item, value: formatNumber(homeIndex.diamond) }
    }

    return item
  })

  if (homeIndex.messages.length > 0) {
    dashboard.messages = [...homeIndex.messages]
  }

  if (homeIndex.entries.length > 0) {
    dashboard.entries = [...homeIndex.entries]
  }

  if (homeIndex.activity_entry.title) {
    dashboard.activityEntry = {
      ...dashboard.activityEntry,
      ...homeIndex.activity_entry,
      icon: legacyActivityAssets.activitySparkIcon
    }
  }

  return dashboard
}

export async function loadHomeDashboard({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore
}: LoadHomeDashboardOptions): Promise<HomeDashboardData> {
  if (dataSource === 'mock') {
    return createMockHomeDashboard()
  }

  await sessionStore.ensureGuestSession()
  const homeIndex = await sessionStore.fetchHomeIndex()

  return createApiHomeDashboard(homeIndex)
}
