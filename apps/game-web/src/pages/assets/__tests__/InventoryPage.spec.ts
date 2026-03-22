import { mount } from '@vue/test-utils'

import InventoryPage from '../InventoryPage.vue'

test('renders wallet and inventory management modules', () => {
  const wrapper = mount(InventoryPage)

  expect(wrapper.text()).toContain('资源钱包')
  expect(wrapper.text()).toContain('普通背包')
  expect(wrapper.text()).toContain('使用/出售说明')
  expect(wrapper.text()).toContain('资源流水入口')
  expect(wrapper.text()).toContain('背包容量')
})
