import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { APIError } from "@/api/http"
import { enterDungeon, getDungeonStatus, rollDungeonDice } from "@/api/modules/dungeon"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import DungeonRunPage from "../DungeonRunPage.vue"

let mockedRouteQuery: Record<string, string> = {}

vi.mock("@/api/modules/dungeon", () => ({
  getDungeonStatus: vi.fn(),
  enterDungeon: vi.fn(),
  rollDungeonDice: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

vi.mock("vue-router", () => ({
  useRoute: () => ({
    query: mockedRouteQuery,
  }),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockedRouteQuery = {}
})

test("enters dungeon when status is missing and rolls forward", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3003,
    nickname: "副本旅人",
  })

  vi.mocked(getDungeonStatus).mockRejectedValue(
    new APIError("dungeon run not found", { status: 404, code: 4041 }),
  )
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 3003,
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
      player_id: 3003,
      total_power: 176,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 176,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 8,
            soul: 24,
            total: 176,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
  vi.mocked(enterDungeon).mockResolvedValue({
    player_id: 3003,
    dungeon_id: 1,
    remain_dice: 15,
    current_floor: 1,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "无掉落",
      spirit_power: 0,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3003,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })
  vi.mocked(rollDungeonDice).mockResolvedValue({
    player_id: 3003,
    dungeon_id: 1,
    remain_dice: 11,
    current_floor: 5,
    status: "boss",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "Boss掉落",
      spirit_power: 8,
      soul_pieces: 1,
    },
    last_battle: {
      battle_type: "dungeon",
      result: "success",
      rounds: 1,
      attacker_power: 180,
      defender_power: 140,
    },
    wallet_snapshot: {
      player_id: 3003,
      spirit_power: 108,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 1,
      manor_plots: 2,
    },
    pet_growth: {
      exp: 50,
      team_total_power: 120,
    },
  } as any)

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("1")
  expect(enterDungeon).toHaveBeenCalledWith(3003, 1)
  expect(wrapper.text()).toContain("战斗前摘要")
  expect(wrapper.text()).toContain("当前队伍战力 152")
  expect(wrapper.text()).toContain("成长总加成 +32")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +0")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(rollDungeonDice).toHaveBeenCalledWith(3003)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("5")
  expect(wrapper.text()).toContain("Boss掉落")
  expect(wrapper.text()).toContain("灵力 +8")
  expect(wrapper.text()).toContain("魂力 +1")
  expect(wrapper.text()).toContain("幻兽经验 +50")
  expect(wrapper.text()).toContain("当前灵力 108")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(wrapper.text()).toContain("当前队伍战力 176")
  expect(wrapper.text()).toContain("成长总加成 +56")
})

test("switches dungeon and restarts with selected dungeon", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3004,
    nickname: "裂隙旅人",
  })

  vi.mocked(getDungeonStatus).mockRejectedValue(
    new APIError("dungeon run not found", { status: 404, code: 4041 }),
  )
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 3004,
    total_power: 128,
    team_size: 1,
    active_team: [
      {
        pet_id: 10011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        exp: 0,
        next_level_exp: 100,
        power: 128,
        power_breakdown: {
          base: 120,
          level: 0,
          bone: 0,
          spirit: 8,
          soul: 0,
          total: 128,
        },
        is_active: true,
      },
    ],
    roster: [],
  })
  vi.mocked(enterDungeon)
    .mockResolvedValueOnce({
      player_id: 3004,
      dungeon_id: 1,
      remain_dice: 15,
      current_floor: 1,
      status: "ongoing",
      started_at: "2026-04-04T00:00:00Z",
      last_reward: {
        label: "无掉落",
        spirit_power: 0,
        soul_pieces: 0,
      },
      wallet_snapshot: {
        player_id: 3004,
        spirit_power: 100,
        spirit_free_wash: 3,
        bone_level: 1,
        soul_pieces: 0,
        manor_plots: 2,
      },
    })
    .mockResolvedValueOnce({
      player_id: 3004,
      dungeon_id: 2,
      remain_dice: 15,
      current_floor: 1,
      status: "ongoing",
      started_at: "2026-04-04T00:00:00Z",
      last_reward: {
        label: "无掉落",
        spirit_power: 0,
        soul_pieces: 0,
      },
      wallet_snapshot: {
        player_id: 3004,
        spirit_power: 100,
        spirit_free_wash: 3,
        bone_level: 1,
        soul_pieces: 0,
        manor_plots: 2,
      },
    })
  vi.mocked(rollDungeonDice).mockResolvedValue({
    player_id: 3004,
    dungeon_id: 2,
    remain_dice: 14,
    current_floor: 2,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "寒渊裂隙掉落",
      spirit_power: 7,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3004,
      spirit_power: 107,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get('[data-dungeon-id="2"]').trigger("click")
  await wrapper.get("button.restart-btn").trigger("click")
  await flushPromises()

  expect(enterDungeon).toHaveBeenLastCalledWith(3004, 2)

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(wrapper.text()).toContain("寒渊裂隙")
  expect(wrapper.text()).toContain("灵力 +7")
  expect(wrapper.text()).toContain("当前灵力 107")
})

test("reads dungeon id from route query when entering missing run", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3005,
    nickname: "地图旅人",
  })
  mockedRouteQuery = { dungeon_id: "2" }

  vi.mocked(getDungeonStatus).mockRejectedValue(
    new APIError("dungeon run not found", { status: 404, code: 4041 }),
  )
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 3005,
    total_power: 128,
    team_size: 1,
    active_team: [
      {
        pet_id: 10011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        exp: 0,
        next_level_exp: 100,
        power: 128,
        power_breakdown: {
          base: 120,
          level: 0,
          bone: 0,
          spirit: 8,
          soul: 0,
          total: 128,
        },
        is_active: true,
      },
    ],
    roster: [],
  })
  vi.mocked(enterDungeon).mockResolvedValue({
    player_id: 3005,
    dungeon_id: 2,
    remain_dice: 15,
    current_floor: 1,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "无掉落",
      spirit_power: 0,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3005,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(enterDungeon).toHaveBeenCalledWith(3005, 2)
  expect(wrapper.text()).toContain("寒渊裂隙")
})

test("restarts into query dungeon when stale run belongs to another dungeon", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3006,
    nickname: "旧档旅人",
  })
  mockedRouteQuery = { dungeon_id: "2" }

  vi.mocked(getDungeonStatus).mockResolvedValue({
    player_id: 3006,
    dungeon_id: 1,
    remain_dice: 12,
    current_floor: 4,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "怪物掉落",
      spirit_power: 5,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3006,
      spirit_power: 125,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 3006,
    total_power: 128,
    team_size: 1,
    active_team: [
      {
        pet_id: 10011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        exp: 0,
        next_level_exp: 100,
        power: 128,
        power_breakdown: {
          base: 120,
          level: 0,
          bone: 0,
          spirit: 8,
          soul: 0,
          total: 128,
        },
        is_active: true,
      },
    ],
    roster: [],
  })
  vi.mocked(enterDungeon).mockResolvedValue({
    player_id: 3006,
    dungeon_id: 2,
    remain_dice: 15,
    current_floor: 1,
    status: "ongoing",
    started_at: "2026-04-04T00:05:00Z",
    last_reward: {
      label: "无掉落",
      spirit_power: 0,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3006,
      spirit_power: 125,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getDungeonStatus).toHaveBeenCalledWith(3006)
  expect(enterDungeon).toHaveBeenCalledWith(3006, 2)
  expect(wrapper.text()).toContain("寒渊裂隙")
  expect(wrapper.text()).toContain("15")
})

test("refreshes prebattle growth summary when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3007,
    nickname: "成长联动",
  })

  vi.mocked(getDungeonStatus).mockResolvedValue({
    player_id: 3007,
    dungeon_id: 1,
    remain_dice: 12,
    current_floor: 4,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "怪物掉落",
      spirit_power: 5,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3007,
      spirit_power: 125,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 3007,
      total_power: 128,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 128,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 0,
            spirit: 8,
            soul: 0,
            total: 128,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
    .mockResolvedValueOnce({
      player_id: 3007,
      total_power: 176,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 176,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 8,
            soul: 24,
            total: 176,
          },
          is_active: true,
        },
      ],
      roster: [],
    })

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("当前队伍战力 128")
  expect(wrapper.text()).toContain("成长总加成 +8")

  syncStore.touch()
  await flushPromises()

  expect(getPetCollection).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("当前队伍战力 176")
  expect(wrapper.text()).toContain("成长总加成 +56")
})
