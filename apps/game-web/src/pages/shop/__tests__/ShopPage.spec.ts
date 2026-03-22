import { mount } from '@vue/test-utils'

import ShopPage from '../ShopPage.vue'

test('renders shop goods and payment flow sections', () => {
  const wrapper = mount(ShopPage)

  expect(wrapper.text()).toContain('推荐商品')
  expect(wrapper.text()).toContain('首充礼包')
  expect(wrapper.text()).toContain('支付流程')
  expect(wrapper.text()).toContain('幂等发货')
  expect(wrapper.find('[data-testid="shop-item-首充礼包"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="shop-item-焚火补给包"] img').exists()).toBe(true)
})
