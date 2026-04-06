import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { APIError } from "@/api/http"
import { claimCultivation, startCultivation } from "@/api/modules/dungeon"
import { getHomeOverview } from "@/api/modules/home"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import CultivationPage from "../CultivationPage.vue"

vi.mock("@/api/modules/home", () => ({
  getHomeOverview: vi.fn(),
}))

vi.mock("@/api/modules/dungeon", () => ({
  claimCultivation: vi.fn(),
  startCultivation: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("loads cultivation status and updates after start and claim", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4004,
    nickname: "修行旅人",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 4004,
    nickname: "修行旅人",
    wallet: {
      player_id: 4004,
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
          reward_preview: "战骨锻造",
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
        self_rank: 10,
        self_score: 900,
      },
    },
    next_action: {
      title: "前往修行",
      description: "先启动当前修行循环。",
      route: "/cultivation",
      cta: "继续修行",
    },
  })
  vi.mocked(startCultivation).mockResolvedValue({
    player_id: 4004,
    spirit_power: 0,
    state: "cultivating",
    start_at: "2026-04-04T00:00:00Z",
    claimable_at: "2026-04-04T01:00:00Z",
  })
  vi.mocked(claimCultivation).mockResolvedValue({
    player_id: 4004,
    spirit_power: 10,
    state: "idle",
    start_at: "2026-04-04T00:00:00Z",
    claimable_at: "",
    pet_growth: {
      exp: 100,
      team_total_power: 144,
    },
  } as any)

  const wrapper = mount(CultivationPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("idle")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()
  expect(startCultivation).toHaveBeenCalledWith(4004)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("cultivating")

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()
  expect(claimCultivation).toHaveBeenCalledWith(4004)
  expect(syncStore.version).toBe(2)
  expect(wrapper.text()).toContain("110")
  expect(wrapper.text()).toContain("幻兽经验 +100")
  expect(wrapper.get("button.ghost").attributes("disabled")).toBeDefined()
})

test("disables claim button until cultivation reward is ready", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4005,
    nickname: "等待修士",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 4005,
    nickname: "等待修士",
    wallet: {
      player_id: 4005,
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
    },
  } as any)
  vi.mocked(startCultivation).mockResolvedValue({
    player_id: 4005,
    spirit_power: 0,
    state: "cultivating",
    start_at: "2026-04-04T00:00:00Z",
    claimable_at: "2099-04-04T01:00:00Z",
  })
  vi.mocked(claimCultivation).mockResolvedValue({
    player_id: 4005,
    spirit_power: 10,
    state: "idle",
  } as any)

  const wrapper = mount(CultivationPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startCultivation).toHaveBeenCalledWith(4005)
  expect(wrapper.get("button.ghost").attributes("disabled")).toBeDefined()
  expect(claimCultivation).not.toHaveBeenCalled()
})

test("resets cultivation state when claim reports cultivation missing", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4006,
    nickname: "失效修士",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 4006,
    nickname: "失效修士",
    wallet: {
      player_id: 4006,
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
        state: "cultivating",
        spirit_power: 10,
        claimable: true,
        claimable_at: "2000-04-04T01:00:00Z",
      },
    },
  } as any)
  vi.mocked(claimCultivation).mockRejectedValue(
    new APIError("cultivation not found", { status: 404, code: 4041 }),
  )

  const wrapper = mount(CultivationPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(claimCultivation).toHaveBeenCalledWith(4006)
  expect(syncStore.version).toBe(0)
  expect(wrapper.text()).toContain("尚未开始修行，请先开始修行。")
  expect(wrapper.text()).toContain("当前状态 idle")
  expect(wrapper.text()).toContain("预计可领取：未开始修行")
  expect(wrapper.get("button.ghost").attributes("disabled")).toBeDefined()
})

test("shows waiting message when claim reports cultivation not ready", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4007,
    nickname: "卡点修士",
  })

  vi.mocked(getHomeOverview).mockResolvedValue({
    player_id: 4007,
    nickname: "卡点修士",
    wallet: {
      player_id: 4007,
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
        state: "cultivating",
        spirit_power: 10,
        claimable: true,
        claimable_at: "2000-04-04T01:00:00Z",
      },
    },
  } as any)
  vi.mocked(claimCultivation).mockRejectedValue(
    new APIError("cultivation not ready", { status: 409, code: 4091 }),
  )

  const wrapper = mount(CultivationPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.ghost").trigger("click")
  await flushPromises()

  expect(claimCultivation).toHaveBeenCalledWith(4007)
  expect(syncStore.version).toBe(0)
  expect(wrapper.text()).toContain("修行尚未完成，请稍后领取。")
  expect(wrapper.text()).toContain("当前状态 cultivating")
  expect(wrapper.text()).toContain("预计可领取：2000-04-04T01:00:00Z")
})
