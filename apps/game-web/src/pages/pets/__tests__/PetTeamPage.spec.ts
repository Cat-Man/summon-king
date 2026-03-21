import { mount } from '@vue/test-utils'

import PetTeamPage from '../PetTeamPage.vue'

test('renders team title', () => {
  const wrapper = mount(PetTeamPage)
  expect(wrapper.text()).toContain('阵容')
})
