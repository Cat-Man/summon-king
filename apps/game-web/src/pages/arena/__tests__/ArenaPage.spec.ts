import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { challengeArena, getArenaIndex, refreshArenaOpponents } from "@/api/modules/arena"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import ArenaPage from "../ArenaPage.vue"

vi.mock("@/api/modules/arena", () => ({
  getArenaIndex: vi.fn(),
  refreshArenaOpponents: vi.fn(),
  challengeArena: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

test("loads arena opponents refreshes and updates after battle", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 7101,
    nickname: "斗法旅人",
  })

  vi.mocked(getArenaIndex).mockResolvedValue({
    record: {
      player_id: 7101,
      current_streak: 0,
      last_win: false,
    },
    opponents: [
      {
        opponent_id: 9001,
        name: "炎甲狼将",
        power: 132,
        current_streak: 2,
      },
      {
        opponent_id: 9002,
        name: "霜羽灵使",
        power: 168,
        current_streak: 4,
      },
    ],
  })
  vi.mocked(refreshArenaOpponents).mockResolvedValue({
    record: {
      player_id: 7101,
      current_streak: 0,
      last_win: false,
    },
    opponents: [
      {
        opponent_id: 9011,
        name: "玄甲斗魁",
        power: 140,
        current_streak: 1,
      },
      {
        opponent_id: 9012,
        name: "赤瞳夜姬",
        power: 174,
        current_streak: 5,
      },
    ],
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
        spirit_power: 18,
        soul_pieces: 1,
      },
      battle: {
        battle_type: "arena",
        result: "success",
        rounds: 1,
        attacker_power: 152,
        defender_power: 140,
      },
      wallet_snapshot: {
        spirit_power: 118,
        bone_level: 1,
        soul_pieces: 1,
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

  expect(getArenaIndex).toHaveBeenCalledWith(7101)
  expect(wrapper.text()).toContain("0")
  expect(wrapper.text()).toContain("败")
  expect(wrapper.text()).toContain("斗法对手")
  expect(wrapper.text()).toContain("炎甲狼将")
  expect(wrapper.text()).toContain("霜羽灵使")
  expect(wrapper.text()).toContain("132")
  expect(wrapper.text()).toContain("168")
  expect(wrapper.text()).toContain("战斗前摘要")
  expect(wrapper.text()).toContain("当前队伍战力 152")
  expect(wrapper.text()).toContain("成长总加成 +32")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +0")

  await wrapper.get('button[data-refresh-opponents="true"]').trigger("click")
  await flushPromises()

  expect(refreshArenaOpponents).toHaveBeenCalledWith(7101)
  expect(wrapper.text()).toContain("玄甲斗魁")
  expect(wrapper.text()).toContain("赤瞳夜姬")

  await wrapper.get('button[data-opponent-id="9011"]').trigger("click")
  await flushPromises()

  expect(challengeArena).toHaveBeenCalledWith(7101, 9011)
  expect(wrapper.text()).toContain("1")
  expect(wrapper.text()).toContain("胜")
  expect(wrapper.text()).toContain("奖励拆分")
  expect(wrapper.text()).toContain("灵力 +18")
  expect(wrapper.text()).toContain("魔魂碎片 +1")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(wrapper.text()).toContain("当前队伍战力 162")
  expect(wrapper.text()).toContain("成长总加成 +42")
  expect(wrapper.text()).toContain("查看排行榜")
  expect(resourceSyncStore.version).toBe(1)
})
