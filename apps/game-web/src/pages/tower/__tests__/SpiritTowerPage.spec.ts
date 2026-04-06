import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getPetCollection } from "@/api/modules/pet"
import { getTowerStatus, startTowerChallenge } from "@/api/modules/tower"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import SpiritTowerPage from "../SpiritTowerPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

test("loads spirit tower status from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
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
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 6202,
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
  expect(wrapper.text()).toContain("战斗前摘要")
  expect(wrapper.text()).toContain("当前队伍战力 152")
  expect(wrapper.text()).toContain("成长总加成 +32")
})

test("challenges spirit tower and refreshes resources", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6205,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 3,
      max_floor: 12,
      remaining_challenges: 2,
      reward_preview: "灵魂碎片",
    })
    .mockResolvedValue({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 4,
      max_floor: 12,
      remaining_challenges: 1,
      reward_preview: "灵魂碎片",
    })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 6205,
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
      player_id: 6205,
      total_power: 172,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 172,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 20,
            soul: 8,
            total: 172,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6205,
    tower: "spirit",
    floor: 4,
    reward: "灵魂碎片",
    remaining_challenges: 1,
    reward_delta: {
      spirit_power: 12,
      soul_pieces: 1,
    },
    battle: {
      battle_no: "tower-spirit-6205-1",
      battle_type: "tower",
      result: "success",
      winner_side: "attacker",
      rounds: 1,
      attacker_power: 190,
      defender_power: 150,
    },
    wallet_snapshot: {
      spirit_power: 132,
      bone_level: 2,
      soul_pieces: 5,
    },
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startTowerChallenge).toHaveBeenCalledWith("spirit", 6205)
  expect(wrapper.text()).toContain("灵力 +12")
  expect(wrapper.text()).toContain("魔魂碎片 +1")
  expect(wrapper.text()).toContain("当前队伍战力 172")
  expect(wrapper.text()).toContain("成长总加成 +52")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(resourceSyncStore.version).toBe(1)
})

test("still syncs resources when spirit tower status refresh fails after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6206,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 3,
      max_floor: 12,
      remaining_challenges: 2,
      reward_preview: "灵魂碎片",
    })
    .mockRejectedValueOnce(new Error("refresh failed"))
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 6206,
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
    player_id: 6206,
    tower: "spirit",
    floor: 4,
    reward: "灵魂碎片",
    remaining_challenges: 1,
    reward_delta: {
      spirit_power: 12,
      soul_pieces: 1,
    },
    battle: {
      battle_no: "tower-spirit-6206-1",
      battle_type: "tower",
      result: "success",
      winner_side: "attacker",
      rounds: 1,
      attacker_power: 195,
      defender_power: 152,
    },
    wallet_snapshot: {
      spirit_power: 132,
      bone_level: 2,
      soul_pieces: 5,
    },
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(resourceSyncStore.version).toBe(1)
  expect(wrapper.text()).toContain("灵力 +12")
  expect(wrapper.text()).toContain("魔魂碎片 +1")
})
