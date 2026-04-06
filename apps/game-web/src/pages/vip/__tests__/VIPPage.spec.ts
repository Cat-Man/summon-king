import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { claimVIPDaily, getVIPIndex } from "@/api/modules/vip"
import { useSessionStore } from "@/stores/session"
import VIPPage from "../VIPPage.vue"

vi.mock("@/api/modules/vip", () => ({
  getVIPIndex: vi.fn(),
  claimVIPDaily: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders vip index and updates after daily claim", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 7001,
    nickname: "贵宾旅人",
  })

  vi.mocked(getVIPIndex).mockResolvedValue({
    player_id: 7001,
    vip_level: 0,
    total_gem_spent: 0,
    daily_claimed: false,
    daily_reward: {
      spirit_power: 5,
      soul_pieces: 0,
    },
    current_benefits: ["每日宝箱：灵力 +5"],
    next_level: {
      level: 1,
      required_gem_spent: 100,
      benefits: ["每日宝箱：灵力 +10"],
    },
  } as any)
  vi.mocked(claimVIPDaily).mockResolvedValue({
    index: {
      player_id: 7001,
      vip_level: 0,
      total_gem_spent: 0,
      daily_claimed: true,
      daily_reward: {
        spirit_power: 5,
        soul_pieces: 0,
      },
      current_benefits: ["每日宝箱：灵力 +5"],
      next_level: {
        level: 1,
        required_gem_spent: 100,
        benefits: ["每日宝箱：灵力 +10"],
      },
    },
    reward_delta: {
      spirit_power: 5,
      soul_pieces: 0,
    },
  } as any)

  const wrapper = mount(VIPPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("VIP0")
  expect(wrapper.text()).toContain("领取每日宝箱")

  await wrapper.get('[data-testid="vip-daily-claim"]').trigger("click")
  await flushPromises()

  expect(wrapper.text()).toContain("今日已领取")
  expect(wrapper.text()).toContain("灵力 +5")
})
