import { runtimeConfig, type GameDataSource } from '@/config/runtime'
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

function createApiHomeDashboard(homeIndex: HomeIndexResponse): HomeDashboardData {
  const dashboard = createMockHomeDashboard()

  dashboard.hero = {
    ...dashboard.hero,
    description: `欢迎回来，${homeIndex.nickname}。第一屏直接告诉玩家今天能做什么、有哪些收益可领、当前主推进目标是什么。`,
    metaValue: homeIndex.nickname
  }

  dashboard.resources = dashboard.resources.map((item) => {
    switch (item.label) {
      case '等级':
        return { ...item, value: `Lv.${homeIndex.level}` }
      case '铜钱':
        return { ...item, value: formatNumber(homeIndex.coin) }
      case '元宝':
        return { ...item, value: formatNumber(homeIndex.diamond) }
      default:
        return item
    }
  })

  dashboard.resourceIcons = dashboard.resourceIcons.map((item) => {
    switch (item.label) {
      case '铜钱':
        return { ...item, value: formatNumber(homeIndex.coin) }
      case '元宝':
        return { ...item, value: formatNumber(homeIndex.diamond) }
      default:
        return item
    }
  })

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
