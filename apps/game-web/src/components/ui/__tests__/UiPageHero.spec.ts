import { mount } from '@vue/test-utils'

import UiPageHero from '../UiPageHero.vue'

test('renders hero content with tone and meta info', () => {
  const wrapper = mount(UiPageHero, {
    props: {
      eyebrow: '今日工作台',
      title: '召唤之王',
      description: '展示首页关键决策信息',
      tone: 'navy',
      metaLabel: '收益聚合',
      metaValue: '3 项待处理'
    }
  })

  expect(wrapper.text()).toContain('今日工作台')
  expect(wrapper.text()).toContain('召唤之王')
  expect(wrapper.text()).toContain('展示首页关键决策信息')
  expect(wrapper.text()).toContain('收益聚合')
  expect(wrapper.text()).toContain('3 项待处理')
  expect(wrapper.classes()).toContain('ui-page-hero--navy')
})
