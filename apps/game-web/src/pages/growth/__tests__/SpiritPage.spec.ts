import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getGrowthWallet, washSpirit } from "@/api/modules/growth"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import SpiritPage from "../SpiritPage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  getSoulState: vi.fn(),
  getManorPlots: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("loads wallet and refreshes after spirit wash", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8101,
    nickname: "战灵旅人",
  })

  vi.mocked(getGrowthWallet)
    .mockResolvedValueOnce({
      player_id: 8101,
      spirit_power: 120,
      spirit_free_wash: 3,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    })
    .mockResolvedValueOnce({
      player_id: 8101,
      spirit_power: 120,
      spirit_free_wash: 2,
      bone_level: 2,
      soul_pieces: 4,
      manor_plots: 2,
    })
  vi.mocked(washSpirit).mockResolvedValue({ washed: true })

  const wrapper = mount(SpiritPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getGrowthWallet).toHaveBeenCalledWith(8101)
  expect(wrapper.text()).toContain("120")
  expect(wrapper.text()).toContain("3")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(washSpirit).toHaveBeenCalledWith(8101)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("2")
})

test("refreshes wallet when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8102,
    nickname: "战灵联动",
  })

  vi.mocked(getGrowthWallet)
    .mockResolvedValueOnce({
      player_id: 8102,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    })
    .mockResolvedValueOnce({
      player_id: 8102,
      spirit_power: 108,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 1,
      manor_plots: 2,
    })

  const wrapper = mount(SpiritPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  syncStore.touch()
  await flushPromises()

  expect(getGrowthWallet).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("108")
  expect(wrapper.text()).toContain("1")
})
