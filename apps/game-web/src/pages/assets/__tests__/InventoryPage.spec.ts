import { mount } from '@vue/test-utils'

import InventoryPage from '../InventoryPage.vue'

test('renders wallet and inventory management modules', () => {
  const wrapper = mount(InventoryPage)

  expect(wrapper.text()).toContain('资源钱包')
  expect(wrapper.text()).toContain('普通背包')
  expect(wrapper.text()).toContain('使用/出售说明')
  expect(wrapper.text()).toContain('资源流水入口')
  expect(wrapper.text()).toContain('背包容量')
  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="wallet-icon-元宝"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="inventory-item-中级经验丹"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="inventory-item-火系进化石"] img').exists()).toBe(true)
})
