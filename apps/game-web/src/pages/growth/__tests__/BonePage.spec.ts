import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getBoneState, upgradeBone } from "@/api/modules/growth"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import BonePage from "../BonePage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  upgradeBone: vi.fn(),
  getSoulState: vi.fn(),
  getManorPlots: vi.fn(),
  harvestManor: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("loads bone state from api and upgrades bone", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8202,
    nickname: "战骨旅人",
  })

  vi.mocked(getBoneState)
    .mockResolvedValueOnce({
      name: "战骨",
      level: 3,
    })
    .mockResolvedValueOnce({
      name: "战骨",
      level: 4,
    })
  vi.mocked(upgradeBone).mockResolvedValue({
    name: "战骨",
    level: 4,
  })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 8202,
      total_power: 120,
      team_size: 1,
      active_team: [],
      roster: [],
    })
    .mockResolvedValueOnce({
      player_id: 8202,
      total_power: 144,
      team_size: 1,
      active_team: [],
      roster: [],
    })

  const wrapper = mount(BonePage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getBoneState).toHaveBeenCalledWith(8202)
  expect(wrapper.text()).toContain("战骨")
  expect(wrapper.text()).toContain("3")
  expect(wrapper.text()).toContain("阵容战力 120")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(upgradeBone).toHaveBeenCalledWith(8202)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("4")
  expect(wrapper.text()).toContain("阵容战力 144")
})
