import { mount } from '@vue/test-utils'

import HomePage from '../HomePage.vue'

test('renders home workstation modules', () => {
  const wrapper = mount(HomePage)

  expect(wrapper.text()).toContain('每日必做')
  expect(wrapper.text()).toContain('签到状态')
  expect(wrapper.text()).toContain('当前修行状态')
  expect(wrapper.text()).toContain('当前推荐副本')
  expect(wrapper.text()).toContain('资源总览')
  expect(wrapper.text()).toContain('功能矩阵')
  expect(wrapper.find('[data-testid="resource-icon-铜钱"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="resource-icon-铜钱"]').attributes('alt')).toContain('铜钱')
  expect(wrapper.find('[data-testid="resource-icon-元宝"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="resource-icon-元宝"]').attributes('alt')).toContain('元宝')
  expect(wrapper.find('[data-testid="activity-entry-icon"]').exists()).toBe(true)
})
