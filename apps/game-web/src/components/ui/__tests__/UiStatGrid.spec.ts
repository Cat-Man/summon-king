import { mount } from '@vue/test-utils'

import UiStatGrid from '../UiStatGrid.vue'

test('renders label value items with tone', () => {
  const wrapper = mount(UiStatGrid, {
    props: {
      tone: 'warm',
      items: [
        { label: '等级', value: 'Lv.36' },
        { label: '战力', value: '18,620' }
      ]
    }
  })

  expect(wrapper.text()).toContain('等级')
  expect(wrapper.text()).toContain('Lv.36')
  expect(wrapper.text()).toContain('战力')
  expect(wrapper.text()).toContain('18,620')
  expect(wrapper.classes()).toContain('ui-stat-grid--warm')
})
