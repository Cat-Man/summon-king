import { mount } from '@vue/test-utils'

import AlliancePage from '../AlliancePage.vue'

test('renders alliance workstation modules', () => {
  const wrapper = mount(AlliancePage)

  expect(wrapper.text()).toContain('联盟公告')
  expect(wrapper.text()).toContain('盟战阶段')
  expect(wrapper.text()).toContain('火能修行')
  expect(wrapper.text()).toContain('联盟建筑')
  expect(wrapper.text()).toContain('联盟动态')
  expect(wrapper.text()).toContain('聊天室入口')
  expect(wrapper.find('[data-testid="alliance-member-avatar-盟主"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="alliance-member-avatar-盟主"] img').attributes('alt')).toContain('盟主')
  expect(wrapper.find('[data-testid="alliance-member-avatar-副盟主"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="alliance-chat-entry-icon"]').exists()).toBe(true)
})
