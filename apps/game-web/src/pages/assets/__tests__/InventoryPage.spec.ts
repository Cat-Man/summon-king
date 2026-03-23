import { flushPromises, mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach } from 'vitest'

import { resetMockInventoryDashboardState } from '@/mocks/inventory-dashboard'

import InventoryPage from '../InventoryPage.vue'

afterEach(() => {
  resetMockInventoryDashboardState()
})

function mountInventoryPage() {
  const pinia = createPinia()
  setActivePinia(pinia)

  return mount(InventoryPage, {
    global: {
      plugins: [pinia]
    }
  })
}

test('renders wallet and inventory management modules', () => {
  const wrapper = mountInventoryPage()

  expect(wrapper.text()).toContain('资源钱包')
  expect(wrapper.text()).toContain('普通背包')
  expect(wrapper.text()).toContain('使用/出售说明')
  expect(wrapper.text()).toContain('资源流水入口')
  expect(wrapper.text()).toContain('背包容量')
  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="wallet-icon-元宝"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="inventory-item-中级经验丹"] img').exists()).toBe(true)
  expect(wrapper.find('[data-testid="inventory-item-火系进化石"] img').exists()).toBe(true)
})

test('loads inventory dashboard asynchronously and renders mapped mock values', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()

  expect(wrapper.text()).toContain('4 / 30')
  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').text()).toContain('0')
  expect(wrapper.find('[data-testid="inventory-item-召唤卷轴"]').exists()).toBe(true)
})

test('clicking use action updates inventory quantity and shows feedback', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()
  await wrapper.get('[data-testid="inventory-use-中级经验丹"]').trigger('click')
  await flushPromises()

  expect(wrapper.find('[data-testid="inventory-item-中级经验丹"]').text()).toContain('x11')
  expect(wrapper.text()).toContain('已使用')
})

test('clicking sell action updates wallet and removes sold item', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()
  await wrapper.get('[data-testid="inventory-sell-召唤卷轴"]').trigger('click')
  await flushPromises()

  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').text()).toContain('50')
  expect(wrapper.find('[data-testid="inventory-item-召唤卷轴"]').exists()).toBe(false)
  expect(wrapper.text()).toContain('已出售')
})
