import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getAllianceHall, getAllianceIndex } from "@/api/modules/alliance"
import AlliancePage from "../AlliancePage.vue"
import { useSessionStore } from "@/stores/session"

vi.mock("@/api/modules/alliance", () => ({
  getAllianceIndex: vi.fn(),
  getAllianceHall: vi.fn(),
  createAlliance: vi.fn(),
  applyToAlliance: vi.fn(),
  approveAllianceApplication: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders alliance hall when player has not joined an alliance", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 5001,
    nickname: "散修旅人",
  })

  vi.mocked(getAllianceIndex).mockResolvedValue({
    player_id: 5001,
    has_alliance: false,
    current_role: "",
    alliance: null,
  } as any)
  vi.mocked(getAllianceHall).mockResolvedValue([
    {
      alliance_id: 1,
      name: "青云盟",
      level: 1,
      member_count: 1,
      member_limit: 20,
      notice: "欢迎新道友",
    },
  ] as any)

  const wrapper = mount(AlliancePage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("联盟大厅")
  expect(wrapper.text()).toContain("青云盟")
  expect(wrapper.text()).toContain("创建联盟")
})
