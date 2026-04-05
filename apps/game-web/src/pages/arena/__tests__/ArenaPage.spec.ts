import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { challengeArena, getArenaStatus } from "@/api/modules/arena"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import ArenaPage from "../ArenaPage.vue"

vi.mock("@/api/modules/arena", () => ({
  getArenaStatus: vi.fn(),
  challengeArena: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
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
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 7101,
      total_power: 152,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 152,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 8,
            soul: 0,
            total: 152,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
    .mockResolvedValueOnce({
      player_id: 7101,
      total_power: 162,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 162,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 18,
            soul: 0,
            total: 162,
          },
          is_active: true,
        },
      ],
      roster: [],
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
  expect(wrapper.text()).toContain("战斗前摘要")
  expect(wrapper.text()).toContain("当前队伍战力 152")
  expect(wrapper.text()).toContain("成长总加成 +32")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +0")

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
  expect(wrapper.text()).toContain("当前队伍战力 162")
  expect(wrapper.text()).toContain("成长总加成 +42")
  expect(wrapper.text()).toContain("查看排行榜")
  expect(resourceSyncStore.version).toBe(1)

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(challengeArena).toHaveBeenCalledWith(7101, false)
  expect(wrapper.text()).toContain("0")
})
