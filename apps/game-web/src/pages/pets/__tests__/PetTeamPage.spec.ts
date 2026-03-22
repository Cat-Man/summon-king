import { mount } from '@vue/test-utils'

import PetTeamPage from '../PetTeamPage.vue'

test('renders team overview and pet roster modules', () => {
  const wrapper = mount(PetTeamPage)

  expect(wrapper.text()).toContain('主力阵容')
  expect(wrapper.text()).toContain('当前战斗队')
  expect(wrapper.text()).toContain('幻兽栏')
  expect(wrapper.text()).toContain('主养成目标')
  expect(wrapper.text()).toContain('上阵策略')
  expect(wrapper.text()).toContain('综合战力')
})
