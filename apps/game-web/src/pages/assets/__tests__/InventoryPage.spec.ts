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

test('high value item is locked by default before selling', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()

  expect(wrapper.find('[data-testid="inventory-lock-召唤卷轴"]').text()).toContain('解除锁定')
  expect(wrapper.get('[data-testid="inventory-sell-召唤卷轴"]').attributes()).toHaveProperty(
    'disabled'
  )
  expect(wrapper.text()).toContain('高价值道具，默认锁定')
})

test('selling item requires confirmation and can be cancelled', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()
  await wrapper.get('[data-testid="inventory-lock-召唤卷轴"]').trigger('click')
  await wrapper.get('[data-testid="inventory-sell-召唤卷轴"]').trigger('click')
  await flushPromises()

  expect(wrapper.find('[data-testid="sell-confirmation"]').exists()).toBe(true)
  expect(wrapper.text()).toContain('预计获得 50 铜钱')

  await wrapper.get('[data-testid="sell-cancel"]').trigger('click')
  await flushPromises()

  expect(wrapper.find('[data-testid="sell-confirmation"]').exists()).toBe(false)
  expect(wrapper.find('[data-testid="inventory-item-召唤卷轴"]').exists()).toBe(true)
  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').text()).toContain('0')
})

test('confirming sell updates wallet and removes sold item', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()
  await wrapper.get('[data-testid="inventory-lock-召唤卷轴"]').trigger('click')
  await wrapper.get('[data-testid="inventory-sell-召唤卷轴"]').trigger('click')
  await flushPromises()
  await wrapper.get('[data-testid="sell-confirm"]').trigger('click')
  await flushPromises()

  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').text()).toContain('50')
  expect(wrapper.find('[data-testid="inventory-item-召唤卷轴"]').exists()).toBe(false)
  expect(wrapper.text()).toContain('已出售')
})

test('clicking recent logs loads and renders asset change history', async () => {
  const wrapper = mountInventoryPage()

  await flushPromises()
  await wrapper.get('[data-testid="view-inventory-logs"]').trigger('click')
  await flushPromises()

  expect(wrapper.text()).toContain('最近流水')
  expect(wrapper.text()).toContain('出售召唤卷轴')
  expect(wrapper.text()).toContain('出售道具')
  expect(wrapper.text()).toContain('原因：inventory_sell')
  expect(wrapper.text()).toContain('+50 铜钱')
  expect(wrapper.text()).toContain('元宝 +0')
})
