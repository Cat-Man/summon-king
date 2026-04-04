import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getTowerStatus } from "@/api/modules/tower"
import { useSessionStore } from "@/stores/session"

import SpiritTowerPage from "../SpiritTowerPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

test("loads spirit tower status from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6202,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus).mockResolvedValue({
    tower: "spirit",
    label: "战灵塔",
    current_floor: 3,
    max_floor: 12,
    remaining_challenges: 2,
    reward_preview: "灵魂碎片",
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getTowerStatus).toHaveBeenCalledWith("spirit", 6202)
  expect(wrapper.text()).toContain("第 3 层")
  expect(wrapper.text()).toContain("2/5")
  expect(wrapper.text()).toContain("灵魂碎片")
})
