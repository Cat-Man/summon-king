import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getPetCollection, savePetTeam, setMainPet } from "@/api/modules/pet"
import { useResourceSyncStore } from "@/stores/resourceSync"
import { useSessionStore } from "@/stores/session"

import PetPage from "../PetPage.vue"

vi.mock("@/api/modules/pet", () => ({
  getPetCollection: vi.fn(),
  savePetTeam: vi.fn(),
  setMainPet: vi.fn(),
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
        exp: 0,
        next_level_exp: 100,
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
        exp: 0,
        next_level_exp: 100,
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
  expect(wrapper.text()).toContain("EXP 0/100")
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
        exp: 40,
        next_level_exp: 100,
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
        exp: 40,
        next_level_exp: 100,
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
        exp: 25,
        next_level_exp: 200,
        power: 144,
      },
    ],
    roster: [
      {
        ...first.roster[0],
        level: 2,
        exp: 25,
        next_level_exp: 200,
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
  expect(wrapper.text()).toContain("EXP 25/200")
})

test("adds pet to active team when slot available", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3001,
    nickname: "主战切换",
  })

  const first = {
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
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: false,
      },
    ],
  }

  const second = {
    player_id: 3001,
    total_power: 276,
    team_size: 2,
    active_team: [
      {
        pet_id: 30011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
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
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: true,
      },
      {
        pet_id: 30011,
        slot: 2,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
  }

  vi.mocked(getPetCollection).mockResolvedValueOnce(first).mockResolvedValueOnce(second)
  vi.mocked(savePetTeam).mockResolvedValueOnce(second)

  const wrapper = mount(PetPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  await wrapper.get('[data-testid="join-team-30012"]').trigger("click")
  await flushPromises()

  expect(savePetTeam).toHaveBeenCalledWith(3001, [30011, 30012])
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("玄甲龟")
  expect(wrapper.text()).toContain("当前主战")
  expect(wrapper.text()).toContain("槽位 2")
})

test("switches main pet without losing secondary slot", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3001,
    nickname: "双槽切换",
  })

  const first = {
    player_id: 3001,
    total_power: 276,
    team_size: 2,
    active_team: [
      {
        pet_id: 30011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
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
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: true,
      },
    ],
  }

  const second = {
    player_id: 3001,
    total_power: 276,
    team_size: 2,
    active_team: [
      {
        pet_id: 30012,
        slot: 1,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: true,
      },
      {
        pet_id: 30011,
        slot: 2,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
    ],
    roster: first.roster,
  }

  vi.mocked(getPetCollection).mockResolvedValueOnce(first).mockResolvedValueOnce(second)
  vi.mocked(setMainPet).mockResolvedValueOnce(second)

  const wrapper = mount(PetPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  await wrapper.get('[data-testid="set-main-30012"]').trigger("click")
  await flushPromises()

  expect(setMainPet).toHaveBeenCalledWith(3001, 30012)
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("玄甲龟")
  expect(wrapper.text()).toContain("槽位 2")
})

test("removes secondary pet from active team", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  const syncStore = useResourceSyncStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 3001,
    nickname: "移出副位",
  })

  const first = {
    player_id: 3001,
    total_power: 276,
    team_size: 2,
    active_team: [
      {
        pet_id: 30011,
        slot: 1,
        name: "初始灵狐",
        level: 1,
        power: 120,
        is_active: true,
      },
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
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
      {
        pet_id: 30012,
        slot: 2,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: true,
      },
    ],
  }

  const second = {
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
      {
        pet_id: 30012,
        slot: 0,
        name: "玄甲龟",
        level: 1,
        power: 156,
        is_active: false,
      },
    ],
  }

  vi.mocked(getPetCollection).mockResolvedValueOnce(first).mockResolvedValueOnce(second)
  vi.mocked(savePetTeam).mockResolvedValueOnce(second)

  const wrapper = mount(PetPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  await wrapper.get('[data-testid="remove-team-30012"]').trigger("click")
  await flushPromises()

  expect(savePetTeam).toHaveBeenCalledWith(3001, [30011])
  expect(syncStore.version).toBe(1)
  expect(wrapper.text()).toContain("120")
  expect(wrapper.find('[data-testid="join-team-30012"]').exists()).toBe(true)
})
