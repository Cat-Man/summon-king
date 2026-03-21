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
})
