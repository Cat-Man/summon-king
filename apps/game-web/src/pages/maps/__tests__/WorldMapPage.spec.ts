import { mount } from '@vue/test-utils'

import WorldMapPage from '../WorldMapPage.vue'

test('renders world map growth path modules', () => {
  const wrapper = mount(WorldMapPage)

  expect(wrapper.text()).toContain('当前位置')
  expect(wrapper.text()).toContain('当前城市副本')
  expect(wrapper.text()).toContain('临近城市')
  expect(wrapper.text()).toContain('移动/传送')
  expect(wrapper.text()).toContain('等级段提示')
  expect(wrapper.text()).toContain('当前推荐目标')
})
