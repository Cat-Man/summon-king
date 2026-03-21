import { mount } from '@vue/test-utils'

import UiChipGroup from '../UiChipGroup.vue'

test('renders chips with tone variant', () => {
  const wrapper = mount(UiChipGroup, {
    props: {
      tone: 'blue',
      items: ['世界地图', '联盟', '修行']
    }
  })

  expect(wrapper.text()).toContain('世界地图')
  expect(wrapper.text()).toContain('联盟')
  expect(wrapper.text()).toContain('修行')
  expect(wrapper.classes()).toContain('ui-chip-group--blue')
})
