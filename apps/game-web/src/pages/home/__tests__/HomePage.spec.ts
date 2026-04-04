import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { useSessionStore } from "@/stores/session"
import { useResourceSyncStore } from "@/stores/resourceSync"
import HomePage from "../HomePage.vue"
import { getHomeOverview } from "@/api/modules/home"

vi.mock("@/api/modules/home", () => ({
  getHomeOverview: vi.fn(),
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
    },
    next_action: {
      title: "前往通天塔",
      description: "今日还有 5 次挑战，先拿战骨强化。",
      route: "/tower/pagoda",
      cta: "继续挑战",
    },
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
  expect(wrapper.text()).toContain("幻兽阵容")
  expect(wrapper.text()).toContain("初始灵狐")
  expect(wrapper.text()).toContain("玄境 · 已开放 2 城")
  expect(wrapper.text()).toContain("第 3 层")
  expect(wrapper.findAllComponents(RouterLinkStub).length).toBeGreaterThan(0)
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
    },
    next_action: {
      title: "前往通天塔",
      description: "继续挑战获取战骨。",
      route: "/tower/pagoda",
      cta: "继续挑战",
    },
  })

  vi.mocked(getHomeOverview)
    .mockResolvedValueOnce(overview(3, 12))
    .mockResolvedValueOnce(overview(5, 11))

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
  expect(wrapper.text()).toContain("第 5 层")
})
