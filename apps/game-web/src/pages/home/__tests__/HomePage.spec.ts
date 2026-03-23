import { beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import { createMockHomeDashboard } from '@/mocks/home-dashboard'
import { loadHomeDashboard } from '@/services/home-dashboard'
import HomePage from '../HomePage.vue'

vi.mock('@/services/home-dashboard', () => ({
  loadHomeDashboard: vi.fn()
}))

function mountHomePage() {
  const pinia = createPinia()
  setActivePinia(pinia)

  return mount(HomePage, {
    global: {
      plugins: [pinia]
    }
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(loadHomeDashboard).mockResolvedValue(createMockHomeDashboard())
})

test('renders home workstation modules', () => {
  const wrapper = mountHomePage()

  expect(wrapper.text()).toContain('每日必做')
  expect(wrapper.text()).toContain('签到状态')
  expect(wrapper.text()).toContain('当前修行状态')
  expect(wrapper.text()).toContain('当前推荐副本')
  expect(wrapper.text()).toContain('资源总览')
  expect(wrapper.text()).toContain('功能矩阵')
  expect(wrapper.find('[data-testid="resource-icon-铜钱"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="resource-icon-铜钱"]').attributes('alt')).toContain('铜钱')
  expect(wrapper.find('[data-testid="resource-icon-元宝"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="resource-icon-元宝"]').attributes('alt')).toContain('元宝')
  expect(wrapper.find('[data-testid="activity-entry-icon"]').exists()).toBe(true)
})

test('loads guest summary asynchronously and renders dynamic profile', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.hero.metaValue = '联调游客2002'
  loadedDashboard.hero.description = '欢迎回来，联调游客2002。第一屏直接告诉玩家今天能做什么、有哪些收益可领、当前主推进目标是什么。'
  loadedDashboard.resources = loadedDashboard.resources.map((item) =>
    item.label === '等级' ? { ...item, value: 'Lv.7' } : item
  )
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()

  expect(loadHomeDashboard).toHaveBeenCalledTimes(1)
  expect(wrapper.text()).toContain('联调游客2002')
  expect(wrapper.text()).toContain('Lv.7')
})
