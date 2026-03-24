import { beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import { loadPetTeamDashboard, savePetTeamSelection } from '@/services/pet-dashboard'
import PetTeamPage from '../PetTeamPage.vue'

const routerPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: routerPush
  })
}))

vi.mock('@/services/pet-dashboard', () => ({
  createInitialPetTeamDashboard: () => ({
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '',
      tone: 'blue' as const,
      metaLabel: '综合战力',
      metaValue: ''
    },
    overview: [],
    team: [],
    roster: [],
    strategies: [],
    focus: {
      name: '',
      description: '',
      nextStep: ''
    }
  }),
  loadPetTeamDashboard: vi.fn(),
  savePetTeamSelection: vi.fn()
}))

function mountPetTeamPage() {
  const pinia = createPinia()
  setActivePinia(pinia)

  return mount(PetTeamPage, {
    global: {
      plugins: [pinia]
    }
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  routerPush.mockReset()
  vi.mocked(loadPetTeamDashboard).mockResolvedValue({
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '接口联调后的阵容页。',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: '270'
    },
    overview: [
      { label: '已上阵', value: '2 / 5' },
      { label: '队伍核心', value: '寒枝鹿灵' },
      { label: '平均等级', value: 'Lv.14' },
      { label: '可替补', value: '1 只' }
    ],
    team: [
      { petId: 2, slot: '1 号位', name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.16', power: '150', icon: '/deer.png' },
      { petId: 1, slot: '2 号位', name: '烈焰狼王', role: '主战输出', level: 'Lv.12', power: '120', icon: '/wolf.png' }
    ],
    roster: [
      { petId: 2, name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.16', status: '已上阵', power: '150', icon: '/deer.png' },
      { petId: 1, name: '烈焰狼王', role: '主战输出', level: 'Lv.12', status: '已上阵', power: '120', icon: '/wolf.png' },
      { petId: 3, name: '雷角牛', role: '前排承伤', level: 'Lv.20', status: '可替补', power: '180' }
    ],
    strategies: ['保存阵容', '支持上下阵'],
    focus: {
      name: '寒枝鹿灵',
      description: '当前控制与辅助能力最稳定。',
      nextStep: '保存阵容后可用于首页战力联动。'
    }
  })
  vi.mocked(savePetTeamSelection).mockResolvedValue({ message: '阵容已保存' })
})

test('renders team overview and saves edited selection', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  expect(wrapper.text()).toContain('主力阵容')
  expect(wrapper.text()).toContain('当前战斗队')
  expect(wrapper.text()).toContain('幻兽栏')
  expect(wrapper.text()).toContain('主养成目标')
  expect(wrapper.text()).toContain('上阵策略')
  expect(wrapper.text()).toContain('270')
  expect(wrapper.find('[data-testid="team-slot-1号位"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="roster-item-烈焰狼王"] img').exists()).toBe(true)

  await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')
  await wrapper.get('[data-testid="save-team"]').trigger('click')

  expect(savePetTeamSelection).toHaveBeenCalledTimes(1)
  expect(wrapper.text()).toContain('阵容已保存')
})

test('adds a bench pet into the next empty slot and saves edited order', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  expect(wrapper.text()).toContain('可替补')

  await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')

  expect(wrapper.get('[data-testid="team-slot-card-3"]').text()).toContain('雷角牛')
  expect(wrapper.get('[data-testid="roster-item-雷角牛"]').text()).toContain('已上阵')

  await wrapper.get('[data-testid="save-team"]').trigger('click')

  expect(savePetTeamSelection).toHaveBeenCalledWith(
    expect.objectContaining({
      petIds: [2, 1, 3]
    })
  )
  expect(wrapper.text()).toContain('450')
  expect(wrapper.get('[data-testid="team-member-3"]').text()).toContain('前排承伤')
})

test('removes an equipped pet back to bench', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  await wrapper.get('[data-testid="team-remove-2"]').trigger('click')

  expect(wrapper.find('[data-testid="team-member-2"]').exists()).toBe(false)
  expect(wrapper.get('[data-testid="roster-item-寒枝鹿灵"]').text()).toContain('可替补')
  expect(wrapper.text()).toContain('120')
})

test('reorders current team before save', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  await wrapper.get('[data-testid="team-move-down-2"]').trigger('click')

  expect(wrapper.get('[data-testid="team-member-1"]').text()).toContain('烈焰狼王')
  expect(wrapper.get('[data-testid="team-member-2"]').text()).toContain('寒枝鹿灵')

  await wrapper.get('[data-testid="save-team"]').trigger('click')

  expect(savePetTeamSelection).toHaveBeenCalledWith(
    expect.objectContaining({
      petIds: [1, 2]
    })
  )
})

test('disables save button when there are no unsaved changes', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  expect(wrapper.get('[data-testid="save-team"]').attributes('disabled')).toBeDefined()

  await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')

  expect(wrapper.get('[data-testid="save-team"]').attributes('disabled')).toBeUndefined()

  await wrapper.get('[data-testid="save-team"]').trigger('click')
  await flushPromises()

  expect(wrapper.get('[data-testid="save-team"]').attributes('disabled')).toBeDefined()
})

test('shows a guard message when trying to add a bench pet to a full team', async () => {
  vi.mocked(loadPetTeamDashboard).mockResolvedValueOnce({
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '满编阵容测试。',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: '780'
    },
    overview: [
      { label: '已上阵', value: '5 / 5' },
      { label: '队伍核心', value: '寒枝鹿灵' },
      { label: '平均等级', value: 'Lv.15' },
      { label: '可替补', value: '1 只' }
    ],
    team: [
      { petId: 2, slot: '1 号位', name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.16', power: '150' },
      { petId: 1, slot: '2 号位', name: '烈焰狼王', role: '主战输出', level: 'Lv.12', power: '120' },
      { petId: 3, slot: '3 号位', name: '雷角牛', role: '前排承伤', level: 'Lv.20', power: '180' },
      { petId: 4, slot: '4 号位', name: '青羽雀', role: '后排输出', level: 'Lv.14', power: '140' },
      { petId: 5, slot: '5 号位', name: '磐石猿', role: '防御守护', level: 'Lv.19', power: '190' }
    ],
    roster: [
      { petId: 2, name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.16', status: '已上阵', power: '150' },
      { petId: 1, name: '烈焰狼王', role: '主战输出', level: 'Lv.12', status: '已上阵', power: '120' },
      { petId: 3, name: '雷角牛', role: '前排承伤', level: 'Lv.20', status: '已上阵', power: '180' },
      { petId: 4, name: '青羽雀', role: '后排输出', level: 'Lv.14', status: '已上阵', power: '140' },
      { petId: 5, name: '磐石猿', role: '防御守护', level: 'Lv.19', status: '已上阵', power: '190' },
      { petId: 6, name: '玄甲龟', role: '减伤守卫', level: 'Lv.11', status: '可替补', power: '90' }
    ],
    strategies: ['保存阵容', '支持上下阵'],
    focus: {
      name: '磐石猿',
      description: '当前为满编阵容。',
      nextStep: '先下阵再替换。'
    }
  })

  const wrapper = mountPetTeamPage()

  await flushPromises()

  await wrapper.get('[data-testid="roster-action-玄甲龟"]').trigger('click')

  expect(wrapper.text()).toContain('当前战斗队已满，请先下阵后再上阵')
  expect(wrapper.get('[data-testid="roster-item-玄甲龟"]').text()).toContain('可替补')
  expect(wrapper.get('[data-testid="team-slot-card-5"]').text()).not.toContain('玄甲龟')
})

test('offers a return-home action after save success', async () => {
  const wrapper = mountPetTeamPage()

  await flushPromises()

  await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')
  await wrapper.get('[data-testid="save-team"]').trigger('click')
  await flushPromises()

  await wrapper.get('[data-testid="return-home-after-save"]').trigger('click')

  expect(routerPush).toHaveBeenCalledWith({ name: 'home' })
})
