import { mount } from '@vue/test-utils'

import CultivationPage from '../CultivationPage.vue'

test('renders cultivation management modules', () => {
  const wrapper = mount(CultivationPage)

  expect(wrapper.text()).toContain('可修行地图')
  expect(wrapper.text()).toContain('修行时长')
  expect(wrapper.text()).toContain('12 小时 VIP2 开放')
  expect(wrapper.text()).toContain('24 小时 VIP5 开放')
  expect(wrapper.text()).toContain('当前修行状态')
  expect(wrapper.text()).toContain('收益规则')
  expect(wrapper.text()).toContain('修行队伍 = 战斗队伍')
  expect(wrapper.text()).toContain('到期领取')
})
