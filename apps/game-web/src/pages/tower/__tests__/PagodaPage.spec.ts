import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getTowerStatus, startTowerChallenge } from "@/api/modules/tower"
import { useSessionStore } from "@/stores/session"

import PagodaPage from "../PagodaPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

test("loads pagoda status and refreshes after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6101,
    nickname: "通天试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "pagoda",
      label: "通天塔",
      current_floor: 0,
      max_floor: 10,
      remaining_challenges: 5,
      reward_preview: "战骨锻造",
    })
    .mockResolvedValueOnce({
      tower: "pagoda",
      label: "通天塔",
      current_floor: 1,
      max_floor: 10,
      remaining_challenges: 4,
      reward_preview: "战骨锻造",
    })
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6101,
    tower: "pagoda",
    floor: 1,
    reward: "战骨锻造",
    remaining_challenges: 4,
  })

  const wrapper = mount(PagodaPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getTowerStatus).toHaveBeenCalledWith("pagoda", 6101)
  expect(wrapper.text()).toContain("第 0 层")
  expect(wrapper.text()).toContain("5/5")
  expect(wrapper.text()).toContain("战骨锻造")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startTowerChallenge).toHaveBeenCalledWith("pagoda", 6101)
  expect(wrapper.text()).toContain("第 1 层")
  expect(wrapper.text()).toContain("4/5")
})
