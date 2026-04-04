import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getLeaderboard } from "@/api/modules/ranking"
import { useSessionStore } from "@/stores/session"

import RankingPage from "../RankingPage.vue"

vi.mock("@/api/modules/ranking", () => ({
  getLeaderboard: vi.fn(),
}))

test("renders leaderboard from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 5005,
    nickname: "榜单旅人",
  })

  vi.mocked(getLeaderboard).mockResolvedValue([
    { rank: 1, player_id: 5005, name: "榜单旅人", score: 1680, arena_streak: 2, updated: Date.now(), is_self: true },
    { rank: 2, player_id: 9001, name: "星痕丶苍穹", score: 1480, arena_streak: 9, updated: Date.now() - 60_000, is_self: false },
  ])

  const wrapper = mount(RankingPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("榜单旅人")
  expect(wrapper.text()).toContain("1680")
  expect(wrapper.text()).toContain("NO. 1")
  expect(wrapper.text()).toContain("我的排名")
  expect(wrapper.text()).toContain("当前连胜 2 场")
  expect(wrapper.text()).toContain("返回首页")
  expect(wrapper.find(".self-entry").exists()).toBe(true)
})
