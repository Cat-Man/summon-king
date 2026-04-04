import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getWorldMap } from "@/api/modules/dungeon"
import { getGrowthWallet } from "@/api/modules/growth"
import { useSessionStore } from "@/stores/session"

import WorldMapPage from "../WorldMapPage.vue"

const push = vi.fn()

vi.mock("@/api/modules/dungeon", () => ({
  getWorldMap: vi.fn(),
}))

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
}))

vi.mock("vue-router", () => ({
  useRouter: () => ({
    push,
  }),
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
          },
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
            unlock_spirit_power: 120,
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
          },
          {
            dungeon_id: 1,
            dungeon_name: "妖窟试炼",
            unlock_spirit_power: 0,
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

  const wrapper = mount(WorldMapPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("玄境")
  expect(wrapper.text()).toContain("晨曦城")
  expect(wrapper.text()).toContain("妖窟试炼")
  expect(wrapper.text()).toContain("寒渊裂隙")
  expect(wrapper.text()).toContain("需灵力 120")

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
})
