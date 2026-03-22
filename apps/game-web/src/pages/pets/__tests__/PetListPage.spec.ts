import { mount } from '@vue/test-utils'

import PetListPage from '../PetListPage.vue'

test('renders pet catalog overview and template cards', () => {
  const wrapper = mount(PetListPage)

  expect(wrapper.text()).toContain('幻兽图鉴')
  expect(wrapper.text()).toContain('图鉴统计')
  expect(wrapper.text()).toContain('已拥有')
  expect(wrapper.text()).toContain('未拥有')
  expect(wrapper.text()).toContain('所属地图')
  expect(wrapper.text()).toContain('技能池')
  expect(wrapper.text()).toContain('来源')
  expect(wrapper.text()).toContain('最低/满资质')
})
