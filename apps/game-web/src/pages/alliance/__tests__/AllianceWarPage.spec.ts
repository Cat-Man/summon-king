import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount, RouterLinkStub } from "@vue/test-utils"

import { getAllianceWarIndex, registerAllianceWarTarget } from "@/api/modules/allianceWar"
import { useSessionStore } from "@/stores/session"
import AllianceWarPage from "../AllianceWarPage.vue"

vi.mock("@/api/modules/allianceWar", () => ({
  getAllianceWarIndex: vi.fn(),
  registerAllianceWarTarget: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
})

test("renders alliance war index and updates after registering target", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)

  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 6101,
    nickname: "盟战旅人",
  })

  vi.mocked(getAllianceWarIndex).mockResolvedValue({
    has_alliance: true,
    current_role: "leader",
    phase: "preparing",
    target_label: "待开放目标",
    can_register: true,
    available_targets: ["赤焰谷", "寒月岭"],
  } as any)
  vi.mocked(registerAllianceWarTarget).mockResolvedValue({
    has_alliance: true,
    current_role: "leader",
    phase: "preparing",
    target_label: "赤焰谷",
    can_register: true,
    available_targets: ["赤焰谷", "寒月岭"],
  } as any)

  const wrapper = mount(AllianceWarPage, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })
  await flushPromises()

  expect(wrapper.text()).toContain("待开放目标")
  expect(wrapper.text()).toContain("赤焰谷")

  await wrapper.get('[data-testid="register-target-赤焰谷"]').trigger("click")
  await flushPromises()

  expect(wrapper.text()).toContain("当前目标")
  expect(wrapper.text()).toContain("赤焰谷")
})
