import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getGrowthWallet, washSpirit } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import SpiritPage from "../SpiritPage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  getSoulState: vi.fn(),
  getManorPlots: vi.fn(),
}))

test("loads wallet and refreshes after spirit wash", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8101,
    nickname: "战灵旅人",
  })

  vi.mocked(getGrowthWallet)
    .mockResolvedValueOnce({
      player_id: 8101,
      spirit_power: 120,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    })
    .mockResolvedValueOnce({
      player_id: 8101,
      spirit_power: 120,
      spirit_free_wash: 2,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    })
  vi.mocked(washSpirit).mockResolvedValue({ washed: true })

  const wrapper = mount(SpiritPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getGrowthWallet).toHaveBeenCalledWith(8101)
  expect(wrapper.text()).toContain("120")
  expect(wrapper.text()).toContain("3")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(washSpirit).toHaveBeenCalledWith(8101)
  expect(wrapper.text()).toContain("2")
})
