import { beforeEach, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'

import { loadPetDetailDashboard } from '@/services/pet-dashboard'
import PetDetailPage from '../PetDetailPage.vue'

vi.mock('vue-router', () => ({
  useRoute: () => ({
    query: { pet_id: '2' }
  })
}))

vi.mock('@/services/pet-dashboard', () => ({
  createInitialPetDetailDashboard: () => ({
    hero: {
      eyebrow: '养成中枢',
      title: '',
      description: '',
      tone: 'orange' as const,
      metaLabel: '综合战力',
      metaValue: ''
    },
    overview: [],
    skills: [],
    bones: [],
    spirits: [],
    souls: [],
    sources: [],
    heroIcon: '',
    evolutionCost: {
      itemName: '',
      itemProgress: '',
      coinProgress: ''
    }
  }),
  loadPetDetailDashboard: vi.fn()
}))

function mountPetDetailPage() {
  const pinia = createPinia()
  setActivePinia(pinia)

  return mount(PetDetailPage, {
    global: {
      plugins: [pinia]
    }
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(loadPetDetailDashboard).mockResolvedValue({
    hero: {
      eyebrow: '养成中枢',
      title: '寒枝鹿灵',
      description: '接口联调后的幻兽详情。',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: '150'
    },
    overview: [
      { label: '等级', value: 'Lv.16' },
      { label: '元素', value: '水系' },
      { label: '星级', value: '2 星' },
      { label: '综合战力', value: '150' }
    ],
    skills: ['寒枝鹿灵技能池：寒枝缠绕 / 冰露护体'],
    bones: ['头骨 Lv.4'],
    spirits: ['水系灵·共鸣'],
    souls: ['水系魂·加护'],
    sources: ['来自北境冰湖的守望者'],
    heroIcon: '/deer.png',
    evolutionCost: {
      itemName: '进化石',
      itemProgress: '6 / 8',
      coinProgress: '180,000 / 200,000'
    }
  })
})

test('renders pet detail growth hub modules from loaded dashboard', async () => {
  const wrapper = mountPetDetailPage()

  await flushPromises()

  expect(wrapper.text()).toContain('技能')
  expect(wrapper.text()).toContain('战骨')
  expect(wrapper.text()).toContain('战灵')
  expect(wrapper.text()).toContain('魔魂')
  expect(wrapper.text()).toContain('进化/升境')
  expect(wrapper.text()).toContain('重生/放生')
  expect(wrapper.text()).toContain('来源说明')
  expect(wrapper.text()).toContain('寒枝鹿灵')
  expect(wrapper.text()).toContain('Lv.16')
  expect(loadPetDetailDashboard).toHaveBeenCalledWith(
    expect.objectContaining({
      petId: 2
    })
  )
  expect(wrapper.find('[data-testid="pet-detail-hero-art"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="pet-detail-evolution-cost"]').text()).toContain('进化石')
})
