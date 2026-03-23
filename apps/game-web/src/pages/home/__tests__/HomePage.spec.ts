import { beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import { createMockHomeDashboard } from '@/mocks/home-dashboard'
import { loadHomeDashboard } from '@/services/home-dashboard'
import HomePage from '../HomePage.vue'

const routerPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: routerPush
  })
}))

vi.mock('@/services/home-dashboard', () => ({
  createInitialHomeDashboard: () => ({
    hero: {
      eyebrow: '今日工作台',
      title: '召唤之王',
      description: '',
      tone: 'navy' as const,
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
      icon: '',
      actionKey: ''
    },
    cultivationSummary: {
      action: '',
      actionKey: '',
      items: []
    },
    entries: [],
    messages: []
  }),
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
  routerPush.mockReset()
})

test('renders home workstation modules', async () => {
  const wrapper = mountHomePage()

  await flushPromises()

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
  expect(wrapper.find('[data-testid="activity-entry-action"]').exists()).toBe(true)
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

test('renders cultivation summary from backend driven dashboard', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.activityEntry.title = '修行收益'
  loadedDashboard.activityEntry.description = '火焰山修行进行中'
  loadedDashboard.activityEntry.note = '剩余 1小时59分'
  loadedDashboard.cultivationSummary = {
    action: '查看修行',
    actionKey: 'cultivation',
    items: [
      { label: '修行地图', value: '火焰山', subtext: '进行中' },
      { label: '结束时间', value: '2026-03-23 10:00', subtext: '剩余 1小时59分' },
      { label: '铜钱收益', value: '240' },
      { label: '幻兽经验', value: '160' }
    ]
  }
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()

  expect(wrapper.text()).toContain('修行收益速览')
  expect(wrapper.text()).toContain('火焰山')
  expect(wrapper.text()).toContain('铜钱收益')
  expect(wrapper.text()).toContain('240')
  expect(wrapper.text()).toContain('查看修行')
})

test('does not render mock guest copy before dashboard load resolves', () => {
  vi.mocked(loadHomeDashboard).mockReturnValueOnce(new Promise(() => {}) as ReturnType<typeof loadHomeDashboard>)

  const wrapper = mountHomePage()

  expect(wrapper.text()).not.toContain('游客1001')
  expect(wrapper.text()).not.toContain('今日未签')
  expect(wrapper.text()).not.toContain('世界消息：青木林地今日双倍经验已开启')
})

test('navigates to mapped page when clicking a matrix entry', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.entries = [{ label: '主线入口', actionKey: 'pet_catalog' }] as never
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()
  await wrapper.get('[data-testid="home-entry-主线入口"]').trigger('click')

  expect(routerPush).toHaveBeenCalledWith({ name: 'pet-catalog' })
})

test('navigates from activity entry using action key instead of copy', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.activityEntry.title = '限时入口'
  loadedDashboard.activityEntry.description = '不是固定文案'
  loadedDashboard.activityEntry.note = '交给 action key'
  ;(loadedDashboard.activityEntry as any).actionKey = 'signin'
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()
  await wrapper.get('[data-testid="activity-entry-action"]').trigger('click')

  expect(routerPush).toHaveBeenCalledWith({ name: 'signin' })
})

test('navigates from daily todo card using action key instead of copy', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.dailyTodos = [
    {
      title: '自定义任务',
      value: '任意描述',
      action: '任意动作',
      actionKey: 'dungeon_run'
    } as never
  ]
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()
  await wrapper.get('[data-testid="todo-action-自定义任务"]').trigger('click')

  expect(routerPush).toHaveBeenCalledWith({ name: 'dungeon-run' })
})

test('navigates cultivation summary action using action key instead of copy', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.cultivationSummary = {
    action: '自定义入口',
    actionKey: 'cultivation',
    items: [
      { label: '修行地图', value: '青木林地', subtext: '进行中' }
    ]
  } as never
  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()

  await flushPromises()
  await wrapper.get('.panel-action-button').trigger('click')

  expect(routerPush).toHaveBeenCalledWith({ name: 'cultivation' })
})
