import { createPinia, setActivePinia } from "pinia"
import { mount, RouterLinkStub } from "@vue/test-utils"

import { useSessionStore } from "@/stores/session"

import GameLayout from "../GameLayout.vue"

const push = vi.fn()

vi.mock("vue-router", async () => {
  const actual = await vi.importActual<typeof import("vue-router")>("vue-router")

  return {
    ...actual,
    useRouter: () => ({
      push,
    }),
  }
})

beforeEach(() => {
  vi.clearAllMocks()
})

test("shows current player and supports logout", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)

  const sessionStore = useSessionStore()
  sessionStore.setSession({
    token: "guest-token",
    playerId: 9301,
    nickname: "壳层旅人",
  })

  const wrapper = mount(GameLayout, {
    global: {
      plugins: [pinia],
      stubs: {
        RouterLink: RouterLinkStub,
        RouterView: true,
      },
    },
  })

  expect(wrapper.text()).toContain("壳层旅人")
  expect(wrapper.text()).toContain("9301")

  await wrapper.get('[data-testid="logout-button"]').trigger("click")

  expect(sessionStore.token).toBe("")
  expect(sessionStore.playerId).toBeNull()
  expect(sessionStore.nickname).toBe("")
  expect(push).toHaveBeenCalledWith("/login")
})
