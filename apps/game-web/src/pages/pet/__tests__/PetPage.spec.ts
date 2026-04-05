import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getPetCollection } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import PetPage from "../PetPage.vue"

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders active team and roster from api", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3001,
    nickname: "幻兽旅人",
  })

  vi.mocked(getPetCollection).mockResolvedValue({
    player_id: 3001,
    total_power: 120,
    team_size: 1,
    active_team: [
      {
        pet_id: 30011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
    roster: [
      {
        pet_id: 30011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
  })

  const wrapper = mount(PetPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(getPetCollection).toHaveBeenCalledWith(3001)
  expect(wrapper.text()).toContain("战斗队")
  expect(wrapper.text()).toContain("幻兽栏")
  expect(wrapper.text()).toContain("初始灵狐")
  expect(wrapper.text()).toContain("120")
  expect(wrapper.findAllComponents(RouterLinkStub).length).toBeGreaterThanOrEqual(3)
})

test("refreshes pet collection when resource sync changes", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3002,
    nickname: "幻兽联动",
  })

  const first = {
    player_id: 3002,
    total_power: 120,
    team_size: 1,
    active_team: [
      {
        pet_id: 30021,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
    roster: [
      {
        pet_id: 30021,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
  }

  const second = {
    ...first,
    total_power: 144,
    active_team: [
      {
        ...first.active_team[0],
        level: 2,
        power: 144,
      },
    ],
    roster: [
      {
        ...first.roster[0],
        level: 2,
        power: 144,
      },
    ],
  }

  vi.mocked(getPetCollection).mockResolvedValueOnce(first).mockResolvedValueOnce(second)

  const wrapper = mount(PetPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  syncStore.touch()
  await flushPromises()

  expect(getPetCollection).toHaveBeenCalledTimes(2)
  expect(wrapper.text()).toContain("144")
  expect(wrapper.text()).toContain("Lv.2")
})
