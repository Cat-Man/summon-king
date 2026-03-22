import { mount } from '@vue/test-utils'

import PetDetailPage from '../PetDetailPage.vue'

test('renders pet detail growth hub modules', () => {
  const wrapper = mount(PetDetailPage)

  expect(wrapper.text()).toContain('技能')
  expect(wrapper.text()).toContain('战骨')
  expect(wrapper.text()).toContain('战灵')
  expect(wrapper.text()).toContain('魔魂')
  expect(wrapper.text()).toContain('进化/升境')
  expect(wrapper.text()).toContain('重生/放生')
  expect(wrapper.text()).toContain('来源说明')
  expect(wrapper.find('[data-testid="pet-detail-hero-art"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="pet-detail-evolution-cost"]').text()).toContain('火系进化石')
})
