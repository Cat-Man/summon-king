import { mount } from '@vue/test-utils'

import ArenaPage from '../ArenaPage.vue'

test('renders arena opponents and reward card', () => {
  const wrapper = mount(ArenaPage)

  expect(wrapper.text()).toContain('竞技场')
  expect(wrapper.find('[data-testid="arena-opponent-1"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="arena-reward-card"] img').exists()).toBe(true)
})
