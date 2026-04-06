import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getManorPlots, harvestManor, plantManor } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import ManorPage from "../ManorPage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  upgradeBone: vi.fn(),
  getSoulState: vi.fn(),
  upgradeSoul: vi.fn(),
  getManorPlots: vi.fn(),
  harvestManor: vi.fn(),
  plantManor: vi.fn(),
}))

test("loads manor plots from api, harvests and plants plots", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8404,
    nickname: "庄园旅人",
  })

  vi.mocked(getManorPlots)
    .mockResolvedValueOnce([
      { plot_id: 1, state: "空闲" },
      { plot_id: 2, state: "成长中" },
    ])
    .mockResolvedValueOnce([
      { plot_id: 1, state: "冷却中" },
      { plot_id: 2, state: "冷却中" },
    ])
    .mockResolvedValueOnce([
      { plot_id: 1, state: "成长中" },
      { plot_id: 2, state: "成长中" },
    ])
  vi.mocked(harvestManor).mockResolvedValue({
    message: "庄园收获完成",
    plots: [
      { plot_id: 1, state: "冷却中" },
      { plot_id: 2, state: "冷却中" },
    ],
  })
  vi.mocked(plantManor).mockResolvedValue({
    message: "庄园种植完成",
    plots: [
      { plot_id: 1, state: "成长中" },
      { plot_id: 2, state: "成长中" },
    ],
  })

  const wrapper = mount(ManorPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getManorPlots).toHaveBeenCalledWith(8404)
  expect(wrapper.text()).toContain("地块 1")
  expect(wrapper.text()).toContain("成长中")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(harvestManor).toHaveBeenCalledWith(8404)
  expect(wrapper.text()).toContain("冷却中")

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(plantManor).toHaveBeenCalledWith(8404)
  expect(wrapper.text()).toContain("成长中")
})
