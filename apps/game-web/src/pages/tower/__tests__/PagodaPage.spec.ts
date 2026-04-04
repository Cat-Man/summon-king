import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { getTowerStatus, startTowerChallenge } from "@/api/modules/tower"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import PagodaPage from "../PagodaPage.vue"

vi.mock("@/api/modules/tower", () => ({
  getTowerStatus: vi.fn(),
  startTowerChallenge: vi.fn(),
}))

test("loads pagoda status and refreshes after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6101,
    nickname: "通天试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "pagoda",
      label: "通天塔",
      current_floor: 0,
      max_floor: 10,
      remaining_challenges: 5,
      reward_preview: "战骨锻造",
    })
    .mockResolvedValueOnce({
      tower: "pagoda",
      label: "通天塔",
      current_floor: 1,
      max_floor: 10,
      remaining_challenges: 4,
      reward_preview: "战骨锻造",
    })
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6101,
    tower: "pagoda",
    floor: 1,
    reward: "战骨锻造",
    remaining_challenges: 4,
    reward_delta: {
      bone_level: 1,
      spirit_power: 5,
    },
    battle: {
      battle_type: "tower",
      result: "success",
      rounds: 1,
      attacker_power: 180,
      defender_power: 120,
    },
    wallet_snapshot: {
      spirit_power: 125,
      bone_level: 3,
      soul_pieces: 0,
    },
  })

  const wrapper = mount(PagodaPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(getTowerStatus).toHaveBeenCalledWith("pagoda", 6101)
  expect(wrapper.text()).toContain("第 0 层")
  expect(wrapper.text()).toContain("5/5")
  expect(wrapper.text()).toContain("战骨锻造")

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(startTowerChallenge).toHaveBeenCalledWith("pagoda", 6101)
  expect(wrapper.text()).toContain("第 1 层")
  expect(wrapper.text()).toContain("4/5")
  expect(wrapper.text()).toContain("战骨 +1")
  expect(wrapper.text()).toContain("灵力 +5")
  expect(wrapper.text()).toContain("战斗摘要")
  expect(wrapper.text()).toContain("战斗结果")
  expect(wrapper.text()).toContain("回合数")
  expect(wrapper.text()).toContain("我方战力")
  expect(wrapper.text()).toContain("敌方战力")
  expect(resourceSyncStore.version).toBe(1)
})

test("still syncs resources when pagoda status refresh fails after challenge", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const resourceSyncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6102,
    nickname: "通天试炼者",
  })

  vi.mocked(getTowerStatus)
    .mockResolvedValueOnce({
      tower: "pagoda",
      label: "通天塔",
      current_floor: 0,
      max_floor: 10,
      remaining_challenges: 5,
      reward_preview: "战骨锻造",
    })
    .mockRejectedValueOnce(new Error("refresh failed"))
  vi.mocked(startTowerChallenge).mockResolvedValue({
    player_id: 6102,
    tower: "pagoda",
    floor: 1,
    reward: "战骨锻造",
    remaining_challenges: 4,
    reward_delta: {
      bone_level: 1,
    },
    battle: {
      battle_type: "tower",
      result: "success",
      rounds: 1,
      attacker_power: 170,
      defender_power: 130,
    },
    wallet_snapshot: {
      spirit_power: 100,
      bone_level: 2,
      soul_pieces: 0,
    },
  })

  const wrapper = mount(PagodaPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(resourceSyncStore.version).toBe(1)
  expect(wrapper.text()).toContain("战骨 +1")
})
