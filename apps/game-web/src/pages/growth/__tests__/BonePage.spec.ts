import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getBoneState } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import BonePage from "../BonePage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  getSoulState: vi.fn(),
  getManorPlots: vi.fn(),
}))

test("loads bone state from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8202,
    nickname: "战骨旅人",
  })

  vi.mocked(getBoneState).mockResolvedValue({
    name: "战骨",
    level: 3,
  })

  const wrapper = mount(BonePage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getBoneState).toHaveBeenCalledWith(8202)
  expect(wrapper.text()).toContain("战骨")
  expect(wrapper.text()).toContain("3")
})
