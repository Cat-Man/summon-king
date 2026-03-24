import { mount } from '@vue/test-utils'

import AllianceWarPage from '../AllianceWarPage.vue'

test('renders alliance war phase and redeem sections', () => {
  const wrapper = mount(AllianceWarPage)

  expect(wrapper.text()).toContain('报名阶段')
  expect(wrapper.text()).toContain('成员签到')
  expect(wrapper.text()).toContain('对阵信息')
  expect(wrapper.text()).toContain('战功兑换')
  expect(wrapper.find('[data-testid="war-commander-盟主"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="war-redeem-焚火晶礼包"] img').exists()).toBe(true)
})
