import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { claimSignin, getSigninIndex } from "@/api/modules/signin"
import SigninPage from "../SigninPage.vue"
import { useSessionStore } from "@/stores/session"

vi.mock("@/api/modules/signin", () => ({
  getSigninIndex: vi.fn(),
  claimSignin: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders signin index and updates after claim", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 7101,
    nickname: "签到旅人",
  })

  vi.mocked(getSigninIndex).mockResolvedValue({
    player_id: 7101,
    signed_today: false,
    current_streak: 0,
    next_reward: {
      spirit_power: 12,
      soul_pieces: 0,
    },
  } as any)
  vi.mocked(claimSignin).mockResolvedValue({
    index: {
      player_id: 7101,
      signed_today: true,
      current_streak: 1,
      next_reward: {
        spirit_power: 15,
        soul_pieces: 0,
      },
    },
    reward_delta: {
      spirit_power: 12,
      soul_pieces: 0,
    },
  } as any)

  const wrapper = mount(SigninPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("今日可签到")
  expect(wrapper.text()).toContain("连续签到 0 天")

  await wrapper.get('[data-testid="signin-claim"]').trigger("click")
  await flushPromises()

  expect(wrapper.text()).toContain("今日已签到")
  expect(wrapper.text()).toContain("连续签到 1 天")
})
