import { mount } from '@vue/test-utils'

import InventoryPage from '../InventoryPage.vue'

test('renders inventory title and init tip', () => {
  const wrapper = mount(InventoryPage)
  expect(wrapper.text()).toContain('背包')
  expect(wrapper.text()).toContain('初始化提示')
})
