import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getManorPlots } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import ManorPage from "../ManorPage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  getSoulState: vi.fn(),
  getManorPlots: vi.fn(),
}))

test("loads manor plots from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8404,
    nickname: "庄园旅人",
  })

  vi.mocked(getManorPlots).mockResolvedValue([
    { plot_id: 1, state: "空闲" },
    { plot_id: 2, state: "成长中" },
  ])

  const wrapper = mount(ManorPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getManorPlots).toHaveBeenCalledWith(8404)
  expect(wrapper.text()).toContain("地块 1")
  expect(wrapper.text()).toContain("成长中")
})
