import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { challengeArena, getArenaStatus } from "@/api/modules/arena"
import { useResourceSyncStore } from "@/stores/resourceSync"
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
  const resourceSyncStore = useResourceSyncStore()
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
      record: {
        player_id: 7101,
        current_streak: 1,
        last_win: true,
      },
      reward_delta: {
        spirit_power: 10,
      },
      battle: {
        battle_type: "arena",
        result: "success",
        rounds: 1,
        attacker_power: 160,
        defender_power: 140,
      },
      wallet_snapshot: {
        spirit_power: 110,
        bone_level: 1,
        soul_pieces: 0,
      },
    })
    .mockResolvedValueOnce({
      record: {
        player_id: 7101,
        current_streak: 0,
        last_win: false,
      },
      reward_delta: {},
      battle: {
        battle_type: "arena",
        result: "fail",
        rounds: 1,
        attacker_power: 160,
        defender_power: 180,
      },
      wallet_snapshot: {
        spirit_power: 110,
        bone_level: 1,
        soul_pieces: 0,
      },
    })

  const wrapper = mount(ArenaPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
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
  expect(wrapper.text()).toContain("奖励拆分")
  expect(wrapper.text()).toContain("灵力 +10")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(wrapper.text()).toContain("查看排行榜")
  expect(resourceSyncStore.version).toBe(1)

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(challengeArena).toHaveBeenCalledWith(7101, false)
  expect(wrapper.text()).toContain("0")
})
