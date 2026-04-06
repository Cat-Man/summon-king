import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getPetCollection } from "@/api/modules/pet"
import { getTowerStatus, startTowerChallenge } from "@/api/modules/tower"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import PagodaPage from "../PagodaPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

test("loads pagoda status and refreshes after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
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
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 6101,
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
      player_id: 6101,
      total_power: 181,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 181,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 48,
            spirit: 13,
            soul: 0,
            total: 181,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6101,
    tower: "pagoda",
    floor: 1,
    reward: "战骨锻造",
    remaining_challenges: 4,
    reward_delta: {
      bone_level: 1,
      spirit_power: 5,
    },
    battle: {
      battle_no: "tower-pagoda-6101-1",
      battle_type: "tower",
      result: "success",
      winner_side: "attacker",
      rounds: 1,
      attacker_power: 180,
      defender_power: 120,
    },
    wallet_snapshot: {
      spirit_power: 125,
      bone_level: 3,
      soul_pieces: 0,
    },
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
  expect(wrapper.text()).toContain("战斗前摘要")
  expect(wrapper.text()).toContain("当前队伍战力 152")
  expect(wrapper.text()).toContain("成长总加成 +32")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +0")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startTowerChallenge).toHaveBeenCalledWith("pagoda", 6101)
  expect(wrapper.text()).toContain("第 1 层")
  expect(wrapper.text()).toContain("4/5")
  expect(wrapper.text()).toContain("战骨 +1")
  expect(wrapper.text()).toContain("灵力 +5")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(wrapper.text()).toContain("当前队伍战力 181")
  expect(wrapper.text()).toContain("成长总加成 +61")
  expect(resourceSyncStore.version).toBe(1)
})

test("still syncs resources when pagoda status refresh fails after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6102,
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
    .mockRejectedValueOnce(new Error("refresh failed"))
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 6102,
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
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6102,
    tower: "pagoda",
    floor: 1,
    reward: "战骨锻造",
    remaining_challenges: 4,
    reward_delta: {
      bone_level: 1,
    },
    battle: {
      battle_no: "tower-pagoda-6102-1",
      battle_type: "tower",
      result: "success",
      winner_side: "attacker",
      rounds: 1,
      attacker_power: 170,
      defender_power: 130,
    },
    wallet_snapshot: {
      spirit_power: 100,
      bone_level: 2,
      soul_pieces: 0,
    },
  })

  const wrapper = mount(PagodaPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(resourceSyncStore.version).toBe(1)
  expect(wrapper.text()).toContain("战骨 +1")
})
