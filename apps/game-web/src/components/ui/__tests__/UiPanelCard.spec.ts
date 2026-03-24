import { mount } from '@vue/test-utils'

import UiPanelCard from '../UiPanelCard.vue'

test('renders panel title and wide modifier', () => {
  const wrapper = mount(UiPanelCard, {
    props: {
      title: '功能矩阵',
      wide: true
    },
    slots: {
      default: '<div>矩阵内容</div>'
    }
  })

  expect(wrapper.text()).toContain('功能矩阵')
  expect(wrapper.text()).toContain('矩阵内容')
  expect(wrapper.classes()).toContain('ui-panel-card--wide')
})
