import { mount } from '@vue/test-utils'

import RankingPage from '../RankingPage.vue'

test('renders ranking list with top player avatars', () => {
  const wrapper = mount(RankingPage)

  expect(wrapper.text()).toContain('排行榜')
  expect(wrapper.find('[data-testid="ranking-player-1"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="ranking-player-2"] img').exists()).toBe(true)
})
