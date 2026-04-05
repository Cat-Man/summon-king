import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getSoulState, upgradeSoul } from "@/api/modules/growth"
import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import SoulPage from "../SoulPage.vue"

vi.mock("@/api/modules/growth", () => ({
  getGrowthWallet: vi.fn(),
  washSpirit: vi.fn(),
  getBoneState: vi.fn(),
  upgradeBone: vi.fn(),
  getSoulState: vi.fn(),
  upgradeSoul: vi.fn(),
  getManorPlots: vi.fn(),
  harvestManor: vi.fn(),
  plantManor: vi.fn(),
}))

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("loads soul state from api and upgrades soul", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8303,
    nickname: "魔魂旅人",
  })

  vi.mocked(getSoulState)
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 9,
    })
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 10,
    })
  vi.mocked(upgradeSoul).mockResolvedValue({
    name: "魔魂",
    power: 10,
  })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 8303,
      total_power: 120,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 120,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 0,
            spirit: 0,
            soul: 0,
            total: 120,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
    .mockResolvedValueOnce({
      player_id: 8303,
      total_power: 136,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 136,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 0,
            spirit: 0,
            soul: 16,
            total: 136,
          },
          is_active: true,
        },
      ],
      roster: [],
    })

  const wrapper = mount(SoulPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getSoulState).toHaveBeenCalledWith(8303)
  expect(wrapper.text()).toContain("魔魂")
  expect(wrapper.text()).toContain("9")
  expect(wrapper.text()).toContain("阵容战力 120")
  expect(wrapper.text()).toContain("当前魔魂加成 +0")
  expect(wrapper.text()).toContain("升级后预计阵容战力 128")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(upgradeSoul).toHaveBeenCalledWith(8303)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("10")
  expect(wrapper.text()).toContain("阵容战力 136")
  expect(wrapper.text()).toContain("当前魔魂加成 +16")
  expect(wrapper.text()).toContain("升级后预计阵容战力 144")
})

test("refreshes soul when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 8304,
    nickname: "魔魂联动",
  })

  vi.mocked(getSoulState)
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 0,
    })
    .mockResolvedValueOnce({
      name: "魔魂",
      power: 1,
    })
  vi.mocked(getPetCollection)
    .mockResolvedValueOnce({
      player_id: 8304,
      total_power: 120,
      team_size: 1,
      active_team: [
        {
          pet_id: 10011,
          slot: 1,
          name: "初始灵狐",
          level: 1,
          exp: 0,
          next_level_exp: 100,
          power: 120,
          power_breakdown: {
            base: 120,
            level: 0,
            bone: 0,
            spirit: 0,
            soul: 0,
            total: 120,
          },
          is_active: true,
        },
      ],
      roster: [],
    })
    .mockResolvedValueOnce({
      player_id: 8304,
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
            spirit: 0,
            soul: 8,
            total: 128,
          },
          is_active: true,
        },
      ],
      roster: [],
    })

  const wrapper = mount(SoulPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  syncStore.touch()
  await flushPromises()

  expect(getSoulState).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("1")
  expect(wrapper.text()).toContain("当前魔魂加成 +8")
  expect(wrapper.text()).toContain("升级后预计阵容战力 136")
})
