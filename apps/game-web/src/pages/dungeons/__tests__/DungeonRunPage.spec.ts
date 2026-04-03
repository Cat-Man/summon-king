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
  })
  vi.mocked(rollDungeonDice).mockResolvedValue({
    player_id: 3003,
    dungeon_id: 1,
    remain_dice: 14,
    current_floor: 2,
    status: "ongoing",
    started_at: "2026-04-04T00:00:00Z",
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
  expect(wrapper.text()).toContain("2")
})
