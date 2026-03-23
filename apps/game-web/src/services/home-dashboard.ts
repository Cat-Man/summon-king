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
    cultivationSummary: {
      action: '',
      items: []
    },
    entries: [],
    messages: []
  }
}

function formatIsoMinute(value: string): string {
  if (!value) {
    return ''
  }
  return value.replace('T', ' ').replace('Z', '').slice(0, 16)
}

function formatDurationCN(seconds: number): string {
  if (seconds <= 0) {
    return '0分'
  }

  const totalMinutes = Math.ceil(seconds / 60)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours === 0) {
    return `${totalMinutes}分`
  }
  if (minutes === 0) {
    return `${hours}小时`
  }
  return `${hours}小时${minutes}分`
}

function formatCultivationStatus(status: string): string {
  switch (status) {
    case 'running':
      return '进行中'
    case 'claimable':
      return '可领取'
    case 'claimed':
      return '已领取'
    case 'idle':
      return '待开始'
    default:
      return status
  }
}

function createApiCultivationSummary(homeIndex: HomeIndexResponse): HomeDashboardData['cultivationSummary'] {
  const summary = homeIndex.cultivation_summary
  if (!summary?.status) {
    return {
      action: '',
      items: []
    }
  }

  const statusLabel = formatCultivationStatus(summary.status)
  const endValue =
    summary.status === 'idle'
      ? '未开始'
      : formatIsoMinute(summary.finished_at)
  const endSubtext =
    summary.status === 'running'
      ? `剩余 ${formatDurationCN(summary.remaining_seconds)}`
      : summary.status === 'claimable'
        ? '收益已到期'
        : summary.status === 'idle'
          ? '可立即开启'
          : undefined

  return {
    action: summary.action,
    items: [
      { label: '修行地图', value: summary.map_name, subtext: statusLabel },
      { label: '结束时间', value: endValue, subtext: endSubtext },
      { label: '铜钱收益', value: formatNumber(summary.reward_coins) },
      { label: '幻兽经验', value: formatNumber(summary.reward_pet_exp) }
    ]
  }
}

function getResponseResourceValue(homeIndex: HomeIndexResponse, label: string): string | undefined {
  return homeIndex.resources.find((item) => item.label === label)?.value
}

function getResponseIconValue(homeIndex: HomeIndexResponse, label: string): string | undefined {
  return homeIndex.resource_icons.find((item) => item.label === label)?.value
}

function createApiHomeDashboard(homeIndex: HomeIndexResponse): HomeDashboardData {
  const dashboard = createInitialHomeDashboard()

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

  dashboard.resourceIcons = ['铜钱', '元宝'].map((label) => {
    const responseValue = getResponseIconValue(homeIndex, label)
    if (responseValue) {
      return { label, value: responseValue, icon: legacyResourceIconMap[label] }
    }

    if (label === '铜钱') {
      return { label, value: formatNumber(homeIndex.coin), icon: legacyResourceIconMap[label] }
    }

    if (label === '元宝') {
      return { label, value: formatNumber(homeIndex.diamond), icon: legacyResourceIconMap[label] }
    }

    return { label, value: '', icon: legacyResourceIconMap[label] }
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
  dashboard.cultivationSummary = createApiCultivationSummary(homeIndex)

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
