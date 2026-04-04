import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getTowerStatus, startTowerChallenge } from "@/api/modules/tower"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import SpiritTowerPage from "../SpiritTowerPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

test("loads spirit tower status from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6202,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus).mockResolvedValue({
    tower: "spirit",
    label: "战灵塔",
    current_floor: 3,
    max_floor: 12,
    remaining_challenges: 2,
    reward_preview: "灵魂碎片",
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getTowerStatus).toHaveBeenCalledWith("spirit", 6202)
  expect(wrapper.text()).toContain("第 3 层")
  expect(wrapper.text()).toContain("2/5")
  expect(wrapper.text()).toContain("灵魂碎片")
})

test("challenges spirit tower and refreshes resources", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6205,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 3,
      max_floor: 12,
      remaining_challenges: 2,
      reward_preview: "灵魂碎片",
    })
    .mockResolvedValue({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 4,
      max_floor: 12,
      remaining_challenges: 1,
      reward_preview: "灵魂碎片",
    })
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6205,
    tower: "spirit",
    floor: 4,
    reward: "灵魂碎片",
    remaining_challenges: 1,
    reward_delta: {
      spirit_power: 12,
      soul_pieces: 1,
    },
    wallet_snapshot: {
      spirit_power: 132,
      bone_level: 2,
      soul_pieces: 5,
    },
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startTowerChallenge).toHaveBeenCalledWith("spirit", 6205)
  expect(wrapper.text()).toContain("灵力 +12")
  expect(wrapper.text()).toContain("魔魂碎片 +1")
  expect(resourceSyncStore.version).toBe(1)
})

test("still syncs resources when spirit tower status refresh fails after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6206,
    nickname: "战灵试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "spirit",
      label: "战灵塔",
      current_floor: 3,
      max_floor: 12,
      remaining_challenges: 2,
      reward_preview: "灵魂碎片",
    })
    .mockRejectedValueOnce(new Error("refresh failed"))
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6206,
    tower: "spirit",
    floor: 4,
    reward: "灵魂碎片",
    remaining_challenges: 1,
    reward_delta: {
      spirit_power: 12,
      soul_pieces: 1,
    },
    wallet_snapshot: {
      spirit_power: 132,
      bone_level: 2,
      soul_pieces: 5,
    },
  })

  const wrapper = mount(SpiritTowerPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(resourceSyncStore.version).toBe(1)
  expect(wrapper.text()).toContain("灵力 +12")
  expect(wrapper.text()).toContain("魔魂碎片 +1")
})
