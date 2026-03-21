import { mount } from '@vue/test-utils'

import AllianceWarPage from '../AllianceWarPage.vue'

test('renders alliance war phase and redeem sections', () => {
  const wrapper = mount(AllianceWarPage)

  expect(wrapper.text()).toContain('报名阶段')
  expect(wrapper.text()).toContain('成员签到')
  expect(wrapper.text()).toContain('对阵信息')
  expect(wrapper.text()).toContain('战功兑换')
})
