import { beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import { loadPetCatalogDashboard } from '@/services/pet-dashboard'
import PetListPage from '../PetListPage.vue'

vi.mock('@/services/pet-dashboard', () => ({
  createInitialPetCatalogDashboard: () => ({
    hero: {
      eyebrow: '图鉴与来源导航',
      title: '幻兽图鉴',
      description: '',
      tone: 'violet' as const,
      metaLabel: '已拥有',
      metaValue: ''
    },
    overview: [],
    sources: [],
    pets: []
  }),
  loadPetCatalogDashboard: vi.fn()
}))

function mountPetListPage() {
  const pinia = createPinia()
  setActivePinia(pinia)

  return mount(PetListPage, {
    global: {
      plugins: [pinia]
    }
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(loadPetCatalogDashboard).mockResolvedValue({
    hero: {
      eyebrow: '图鉴与来源导航',
      title: '幻兽图鉴',
      description: '接口联调后的幻兽列表。',
      tone: 'violet',
      metaLabel: '已拥有',
      metaValue: '3 / 3'
    },
    overview: [
      { label: '已拥有', value: '3' },
      { label: '未拥有', value: '0' },
      { label: '队伍推荐', value: '雷角牛' },
      { label: '获取来源', value: '副本 / 修行 / 礼包' }
    ],
    sources: ['地图副本：按地图掉落对应幻兽与召唤球'],
    pets: [
      {
        petId: 1,
        name: '烈焰狼王',
        status: '已拥有',
        map: '青木林地',
        skillPool: '烈焰撕咬 / 焚风追击',
        aptitude: '战力 120 / 星级 1 星',
        source: '地图副本 / 召唤球',
        level: 'Lv.12',
        power: '120',
        icon: '/wolf.png'
      }
    ]
  })
})

test('renders pet catalog overview from loaded dashboard', async () => {
  const wrapper = mountPetListPage()

  await flushPromises()

  expect(wrapper.text()).toContain('幻兽图鉴')
  expect(wrapper.text()).toContain('图鉴统计')
  expect(wrapper.text()).toContain('已拥有')
  expect(wrapper.text()).toContain('未拥有')
  expect(wrapper.text()).toContain('雷角牛')
  expect(wrapper.text()).toContain('所属地图')
  expect(wrapper.text()).toContain('技能池')
  expect(wrapper.text()).toContain('来源')
  expect(wrapper.text()).toContain('战力 120 / 星级 1 星')
  expect(loadPetCatalogDashboard).toHaveBeenCalledTimes(1)
  expect(wrapper.find('[data-testid="pet-card-烈焰狼王"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="pet-card-烈焰狼王"] img').attributes('alt')).toContain('烈焰狼王')
})
