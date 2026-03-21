import { mount } from '@vue/test-utils'

import VIPPage from '../VIPPage.vue'

test('renders vip value and privilege sections', () => {
  const wrapper = mount(VIPPage)

  expect(wrapper.text()).toContain('当前 VIP')
  expect(wrapper.text()).toContain('距离下一档')
  expect(wrapper.text()).toContain('每日宝箱')
  expect(wrapper.text()).toContain('当前权益')
  expect(wrapper.text()).toContain('下一档新增权益')
})
