import { mount } from '@vue/test-utils'

import DashboardPage from '../DashboardPage.vue'

test('renders admin dashboard title', () => {
  const wrapper = mount(DashboardPage)
  expect(wrapper.text()).toContain('运营后台')
})
