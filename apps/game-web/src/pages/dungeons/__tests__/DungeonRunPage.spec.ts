import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import { APIError } from "@/api/http"
import { enterDungeon, getDungeonStatus, rollDungeonDice } from "@/api/modules/dungeon"
import { useSessionStore } from "@/stores/session"

import DungeonRunPage from "../DungeonRunPage.vue"

vi.mock("@/api/modules/dungeon", () => ({
  getDungeonStatus: vi.fn(),
  enterDungeon: vi.fn(),
  rollDungeonDice: vi.fn(),
}))

test("enters dungeon when status is missing and rolls forward", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3003,
    nickname: "副本旅人",
  })

  vi.mocked(getDungeonStatus).mockRejectedValue(
    new APIError("dungeon run not found", { status: 404, code: 4041 }),
  )
  vi.mocked(enterDungeon).mockResolvedValue({
    player_id: 3003,
    dungeon_id: 1,
    remain_dice: 15,
    current_floor: 1,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "无掉落",
      spirit_power: 0,
      soul_pieces: 0,
    },
    wallet_snapshot: {
      player_id: 3003,
      spirit_power: 100,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 0,
      manor_plots: 2,
    },
  })
  vi.mocked(rollDungeonDice).mockResolvedValue({
    player_id: 3003,
    dungeon_id: 1,
    remain_dice: 11,
    current_floor: 5,
    status: "boss",
    started_at: "2026-04-04T00:00:00Z",
    last_reward: {
      label: "Boss掉落",
      spirit_power: 8,
      soul_pieces: 1,
    },
    wallet_snapshot: {
      player_id: 3003,
      spirit_power: 108,
      spirit_free_wash: 3,
      bone_level: 1,
      soul_pieces: 1,
      manor_plots: 2,
    },
  })

  const wrapper = mount(DungeonRunPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("1")
  expect(enterDungeon).toHaveBeenCalledWith(3003, 1)

  await wrapper.get("button.primary").trigger("click")
  await flushPromises()

  expect(rollDungeonDice).toHaveBeenCalledWith(3003)
  expect(wrapper.text()).toContain("5")
  expect(wrapper.text()).toContain("Boss掉落")
  expect(wrapper.text()).toContain("灵力 +8")
  expect(wrapper.text()).toContain("魂力 +1")
  expect(wrapper.text()).toContain("当前灵力 108")
})
