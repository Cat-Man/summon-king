import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { challengeArena, getArenaStatus } from "@/api/modules/arena"
import { useSessionStore } from "@/stores/session"

import ArenaPage from "../ArenaPage.vue"

vi.mock("@/api/modules/arena", () => ({
  getArenaStatus: vi.fn(),
  challengeArena: vi.fn(),
}))

test("loads arena record and updates after battle", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 7101,
    nickname: "斗法旅人",
  })

  vi.mocked(getArenaStatus).mockResolvedValue({
    player_id: 7101,
    current_streak: 0,
    last_win: false,
  })
  vi.mocked(challengeArena)
    .mockResolvedValueOnce({
      player_id: 7101,
      current_streak: 1,
      last_win: true,
    })
    .mockResolvedValueOnce({
      player_id: 7101,
      current_streak: 0,
      last_win: false,
    })

  const wrapper = mount(ArenaPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getArenaStatus).toHaveBeenCalledWith(7101)
  expect(wrapper.text()).toContain("0")
  expect(wrapper.text()).toContain("败")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(challengeArena).toHaveBeenCalledWith(7101, true)
  expect(wrapper.text()).toContain("1")
  expect(wrapper.text()).toContain("胜")

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(challengeArena).toHaveBeenCalledWith(7101, false)
  expect(wrapper.text()).toContain("0")
})
