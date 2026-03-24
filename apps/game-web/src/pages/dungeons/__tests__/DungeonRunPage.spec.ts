import { mount } from '@vue/test-utils'

import DungeonRunPage from '../DungeonRunPage.vue'

test('renders dungeon expectation management modules', () => {
  const wrapper = mount(DungeonRunPage)

  expect(wrapper.text()).toContain('副本名称')
  expect(wrapper.text()).toContain('所属城市')
  expect(wrapper.text()).toContain('推荐等级段')
  expect(wrapper.text()).toContain('地图简介')
  expect(wrapper.text()).toContain('Boss 主要奖励')
  expect(wrapper.text()).toContain('经验门槛')
  expect(wrapper.text()).toContain('地图主题')
  expect(wrapper.text()).toContain('剩余次数')
  expect(wrapper.find('[data-testid="dungeon-reward-火系进化石"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="dungeon-reward-幻兽经验丹"] img').exists()).toBe(true)
})
