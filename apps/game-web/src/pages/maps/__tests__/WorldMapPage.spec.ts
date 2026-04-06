import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getWorldMap } from "@/api/modules/dungeon"
import { getGrowthWallet } from "@/api/modules/growth"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import WorldMapPage from "../WorldMapPage.vue"

const push = vi.fn()

vi.mock("@/api/modules/dungeon", () => ({
  getWorldMap: vi.fn(),
}))

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

vi.mock("vue-router", () => ({
  useRouter: () => ({
    push,
  }),
  RouterLink: RouterLinkStub,
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders world map from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4001,
    nickname: "地图守卫",
  })

  vi.mocked(getWorldMap).mockResolvedValue({
    name: "玄境",
    cities: [
      {
        city_id: 1,
        name: "晨曦城",
        region: "东境",
        loc_x: 110.5,
        loc_y: 220.4,
        dungeons: [
          {
            dungeon_id: 1,
            dungeon_name: "妖窟试炼",
            unlock_spirit_power: 0,
            unlock_bone_level: 0,
            unlock_soul_pieces: 0,
          },
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
            unlock_spirit_power: 120,
            unlock_bone_level: 2,
            unlock_soul_pieces: 3,
          },
        ],
      },
      {
        city_id: 2,
        name: "霞光堡",
        region: "南境",
        loc_x: 190.8,
        loc_y: 180.1,
        dungeons: [
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
            unlock_spirit_power: 120,
            unlock_bone_level: 2,
            unlock_soul_pieces: 3,
          },
          {
            dungeon_id: 1,
            dungeon_name: "妖窟试炼",
            unlock_spirit_power: 0,
            unlock_bone_level: 0,
            unlock_soul_pieces: 0,
          },
        ],
      },
    ],
  })
  vi.mocked(getGrowthWallet).mockResolvedValue({
    player_id: 4001,
    spirit_power: 100,
    spirit_free_wash: 3,
    bone_level: 1,
    soul_pieces: 0,
    manor_plots: 2,
  })
  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 4001,
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

  const wrapper = mount(WorldMapPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("玄境")
  expect(wrapper.text()).toContain("晨曦城")
  expect(wrapper.text()).toContain("妖窟试炼")
  expect(wrapper.text()).toContain("寒渊裂隙")
  expect(wrapper.text()).toContain("需灵力 120")
  expect(wrapper.text()).toContain("战骨 2")
  expect(wrapper.text()).toContain("魔魂 3")
  expect(wrapper.text()).toContain("当前成长总加成 +32")
  expect(wrapper.text()).toContain("战骨 +24")
  expect(wrapper.text()).toContain("战灵 +8")
  expect(wrapper.text()).toContain("魔魂 +0")
  expect(wrapper.text()).toContain("还差灵力 20")
  expect(wrapper.text()).toContain("还差战骨 1（约 +24 战力）")
  expect(wrapper.text()).toContain("还差魔魂 3（约 +24 战力）")
  expect(wrapper.text()).toContain("去修行")
  expect(wrapper.text()).toContain("去战骨")
  expect(wrapper.text()).toContain("去魔魂")

  const lockedEntry = wrapper.get('[data-city-id="1"][data-dungeon-id="2"]')
  expect(lockedEntry.attributes("disabled")).toBeDefined()
  await lockedEntry.trigger("click")
  expect(push).not.toHaveBeenCalled()

  await wrapper.get('[data-city-id="1"][data-dungeon-id="1"]').trigger("click")

  expect(push).toHaveBeenCalledWith({
    name: "dungeon",
    query: {
      dungeon_id: "1",
    },
  })

  const ctas = wrapper.findAllComponents(RouterLinkStub)
  const routes = ctas.map((cta) => cta.props("to"))
  expect(routes).toContain("/cultivation")
  expect(routes).toContain("/growth/bone")
  expect(routes).toContain("/growth/soul")
})

test("refreshes growth hints when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 4002,
    nickname: "地图联动",
  })

  vi.mocked(getWorldMap).mockResolvedValue({
    name: "玄境",
    cities: [
      {
        city_id: 1,
        name: "晨曦城",
        region: "东境",
        loc_x: 110.5,
        loc_y: 220.4,
        dungeons: [
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
            unlock_spirit_power: 120,
            unlock_bone_level: 2,
            unlock_soul_pieces: 3,
          },
        ],
      },
    ],
  })
  vi.mocked(getGrowthWallet)
    .mockResolvedValueOnce({
      player_id: 4002,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    })
    .mockResolvedValueOnce({
      player_id: 4002,
      spirit_power: 120,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 3,
      manor_plots: 2,
    })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 4002,
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
      player_id: 4002,
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

  const wrapper = mount(WorldMapPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("当前成长总加成 +8")
  expect(wrapper.text()).toContain("还差战骨 1（约 +24 战力）")
  expect(wrapper.text()).toContain("还差魔魂 3（约 +24 战力）")

  syncStore.touch()
  await flushPromises()

  expect(getGrowthWallet).toHaveBeenCalledTimes(2)
  expect(getPetCollection).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("当前成长总加成 +56")
  expect(wrapper.text()).not.toContain("还差战骨 1（约 +24 战力）")
  expect(wrapper.text()).not.toContain("还差魔魂 3（约 +24 战力）")
})
