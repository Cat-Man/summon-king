import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

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
  })

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
})
