import { mount } from '@vue/test-utils'

import SigninPage from '../SigninPage.vue'

test('renders signin rewards and redeem sections', () => {
  const wrapper = mount(SigninPage)

  expect(wrapper.text()).toContain('今日签到')
  expect(wrapper.text()).toContain('连续签到')
  expect(wrapper.text()).toContain('礼包中心')
  expect(wrapper.text()).toContain('兑换码')
})
