import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { useSessionStore } from "@/stores/session"
import HomePage from "../HomePage.vue"
import { getHomeOverview } from "@/api/modules/home"

vi.mock("@/api/modules/home", () => ({
  getHomeOverview: vi.fn(),
}))

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
    },
  })

  const wrapper = mount(HomePage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("测试玩家")
  expect(wrapper.text()).toContain("玄境 · 已开放 2 城")
  expect(wrapper.text()).toContain("第 3 层")
})
