import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { useSessionStore } from "@/stores/session"
import { useResourceSyncStore } from "@/stores/resourceSync"
import HomePage from "../HomePage.vue"
import { getHomeOverview } from "@/api/modules/home"
import { getPetCollection } from "@/api/modules/pet"

vi.mock("@/api/modules/home", () => ({
  getHomeOverview: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders overview from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 2001,
    nickname: "测试玩家",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 2001,
    nickname: "测试玩家",
    wallet: {
      player_id: 2001,
      spirit_power: 120,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    },
    modules: {
      map_label: "玄境 · 已开放 2 城",
      map_city_count: 2,
      dungeon: {
        status: "ongoing",
        current_floor: 3,
        remain_dice: 12,
        dungeon_id: 1,
      },
      cultivation: {
        state: "cultivating",
        spirit_power: 10,
        claimable: false,
        claimable_at: "",
      },
      pet: {
        total_power: 120,
        active_count: 1,
        starter_pet_name: "初始灵狐",
      },
      tower: {
        pagoda: {
          current_floor: 0,
          remaining_challenges: 5,
          reward_preview: "战骨强韧",
        },
        spirit: {
          current_floor: 3,
          remaining_challenges: 2,
          reward_preview: "灵魂碎片",
        },
      },
      arena: {
        current_streak: 2,
        last_win: true,
      },
      ranking: {
        self_rank: 1,
        self_score: 1680,
        arena_streak: 2,
      },
    },
    next_action: {
      title: "前往通天塔",
      description: "今日还有 5 次挑战，先拿战骨强化。",
      route: "/tower/pagoda",
      cta: "继续挑战",
    },
  } as any)
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 2001,
    total_power: 168,
    team_size: 1,
    active_team: [
      {
        pet_id: 10011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        exp: 0,
        next_level_exp: 100,
        power: 168,
        power_breakdown: {
          base: 120,
          level: 0,
          bone: 24,
          spirit: 8,
          soul: 16,
          total: 168,
        },
        is_active: true,
      },
    ],
    roster: [],
  })

  const wrapper = mount(HomePage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("测试玩家")
  expect(wrapper.text()).toContain("下一步推荐")
  expect(wrapper.text()).toContain("通天塔")
  expect(wrapper.text()).toContain("战灵塔")
  expect(wrapper.text()).toContain("竞技场")
  expect(wrapper.text()).toContain("排行榜")
  expect(wrapper.text()).toContain("当前连胜 2 场")
  expect(wrapper.text()).toContain("第 1 名")
  expect(wrapper.text()).toContain("幻兽阵容")
  expect(wrapper.text()).toContain("初始灵狐")
  expect(wrapper.text()).toContain("养成总加成 +48")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +16")
  expect(wrapper.text()).toContain("玄境 · 已开放 2 城")
  expect(wrapper.text()).toContain("第 3 层")
  const routes = wrapper.findAllComponents(RouterLinkStub).map((component) => component.props("to"))
  expect(routes).toContain("/arena")
  expect(routes).toContain("/ranking")
  expect(routes).toContain("/growth/bone")
  expect(routes).toContain("/growth/spirit")
  expect(routes).toContain("/growth/soul")
})

test("refreshes overview when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 2002,
    nickname: "联动玩家",
  })

  const overview = (floor: number, remainDice: number) => ({
    player_id: 2002,
    nickname: "联动玩家",
    wallet: {
      player_id: 2002,
      spirit_power: 120 + floor,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    },
    modules: {
      map_label: "玄境 · 已开放 2 城",
      map_city_count: 2,
      dungeon: {
        status: floor >= 5 ? "boss" : "ongoing",
        current_floor: floor,
        remain_dice: remainDice,
        dungeon_id: 1,
      },
      cultivation: {
        state: "cultivating",
        spirit_power: 10,
        claimable: false,
        claimable_at: "",
      },
      pet: {
        total_power: 120 + floor,
        active_count: 1,
        starter_pet_name: "初始灵狐",
      },
      tower: {
        pagoda: {
          current_floor: floor - 2,
          remaining_challenges: 5 - floor,
          reward_preview: "战骨锻造",
        },
        spirit: {
          current_floor: 3,
          remaining_challenges: 2,
          reward_preview: "灵魂碎片",
        },
      },
      arena: {
        current_streak: floor - 1,
        last_win: floor >= 5,
      },
      ranking: {
        self_rank: 4 - Math.min(floor, 3),
        self_score: 1400 + floor * 40,
        arena_streak: floor - 1,
      },
    },
    next_action: {
      title: "前往通天塔",
      description: "继续挑战获取战骨。",
      route: "/tower/pagoda",
      cta: "继续挑战",
    },
  } as any)

  vi.mocked(getHomeOverview)
    .mockResolvedValueOnce(overview(3, 12))
    .mockResolvedValueOnce(overview(5, 11))
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 2002,
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
      player_id: 2002,
      total_power: 168,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 168,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 24,
            spirit: 8,
            soul: 16,
            total: 168,
          },
          is_active: true,
        },
      ],
      roster: [],
    })

  const wrapper = mount(HomePage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  syncStore.touch()
  await flushPromises()

  expect(getHomeOverview).toHaveBeenCalledTimes(2)
  expect(getPetCollection).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("第 5 层")
  expect(wrapper.text()).toContain("养成总加成 +48")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +16")
})

test("renders cultivation card as ready to claim when overview says claimable", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 2003,
    nickname: "可领修士",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 2003,
    nickname: "可领修士",
    wallet: {
      player_id: 2003,
      spirit_power: 130,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    },
    modules: {
      map_label: "玄境 · 已开放 2 城",
      map_city_count: 2,
      dungeon: {
        status: "idle",
        current_floor: 0,
        remain_dice: 0,
        dungeon_id: 0,
      },
      cultivation: {
        state: "cultivating",
        spirit_power: 10,
        claimable: true,
        claimable_at: "2026-04-06T11:00:00Z",
      },
      pet: {
        total_power: 120,
        active_count: 1,
        starter_pet_name: "初始灵狐",
      },
      tower: {
        pagoda: {
          current_floor: 0,
          remaining_challenges: 5,
          reward_preview: "战骨强韧",
        },
        spirit: {
          current_floor: 0,
          remaining_challenges: 5,
          reward_preview: "灵魂碎片",
        },
      },
      arena: {
        current_streak: 0,
        last_win: false,
      },
      ranking: {
        self_rank: 9,
        self_score: 950,
      },
    },
    next_action: {
      title: "领取修行收益",
      description: "修行已完成，先把灵力写回钱包。",
      route: "/cultivation",
      cta: "立即领取",
    },
  } as any)
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 2003,
    total_power: 120,
    team_size: 1,
    active_team: [],
    roster: [],
  })

  const wrapper = mount(HomePage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  const cultivationCard = wrapper
    .findAll(".overview-card")
    .find((card) => card.text().includes("修行进度"))
  expect(cultivationCard?.text()).toContain("可领取")
  expect(cultivationCard?.text()).toContain("立即领取")
})

test("renders cultivation card as start action when cultivation is idle", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 2004,
    nickname: "起步修士",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 2004,
    nickname: "起步修士",
    wallet: {
      player_id: 2004,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
    modules: {
      map_label: "玄境 · 已开放 2 城",
      map_city_count: 2,
      dungeon: {
        status: "idle",
        current_floor: 0,
        remain_dice: 0,
        dungeon_id: 0,
      },
      cultivation: {
        state: "idle",
        spirit_power: 0,
        claimable: false,
        claimable_at: "",
      },
      pet: {
        total_power: 120,
        active_count: 1,
        starter_pet_name: "初始灵狐",
      },
      tower: {
        pagoda: {
          current_floor: 0,
          remaining_challenges: 5,
          reward_preview: "战骨强韧",
        },
        spirit: {
          current_floor: 0,
          remaining_challenges: 5,
          reward_preview: "灵魂碎片",
        },
      },
      arena: {
        current_streak: 0,
        last_win: false,
      },
      ranking: {
        self_rank: 9,
        self_score: 900,
      },
    },
    next_action: {
      title: "开始修行",
      description: "先启动当前修行循环。",
      route: "/cultivation",
      cta: "开始修行",
    },
  } as any)
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 2004,
    total_power: 120,
    team_size: 1,
    active_team: [],
    roster: [],
  })

  const wrapper = mount(HomePage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  const cultivationCard = wrapper
    .findAll(".overview-card")
    .find((card) => card.text().includes("修行进度"))
  expect(cultivationCard?.text()).toContain("未开始")
  expect(cultivationCard?.text()).toContain("开始修行")
})
