import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getSoulState, upgradeSoul } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import SoulPage from "../SoulPage.vue"

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

test("loads soul state from api and upgrades soul", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8303,
    nickname: "魔魂旅人",
  })

  vi.mocked(getSoulState)
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 9,
    })
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 10,
    })
  vi.mocked(upgradeSoul).mockResolvedValue({
    name: "魔魂",
    power: 10,
  })

  const wrapper = mount(SoulPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getSoulState).toHaveBeenCalledWith(8303)
  expect(wrapper.text()).toContain("魔魂")
  expect(wrapper.text()).toContain("9")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(upgradeSoul).toHaveBeenCalledWith(8303)
  expect(wrapper.text()).toContain("10")
})
