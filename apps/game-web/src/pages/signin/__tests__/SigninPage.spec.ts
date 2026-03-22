import { mount } from '@vue/test-utils'

import SigninPage from '../SigninPage.vue'

test('renders signin rewards and redeem sections', () => {
  const wrapper = mount(SigninPage)

  expect(wrapper.text()).toContain('今日签到')
  expect(wrapper.text()).toContain('连续签到')
  expect(wrapper.text()).toContain('礼包中心')
  expect(wrapper.text()).toContain('兑换码')
  expect(wrapper.find('[data-testid="signin-reward-Day 3"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="signin-reward-Day 4"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="signin-gift-center-icon"]').exists()).toBe(true)
})
