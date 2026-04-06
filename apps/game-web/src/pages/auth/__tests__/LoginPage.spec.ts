import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import LoginPage from "../LoginPage.vue"
import { apiRequest } from "@/api/http"
import { useSessionStore } from "@/stores/session"

const push = vi.fn()
const route = {
  query: {
    redirect: "/arena",
  },
}

vi.mock("@/api/http", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/api/http")>()
  return {
    ...actual,
    apiRequest: vi.fn(),
  }
})

vi.mock("vue-router", async () => {
  const actual = await vi.importActual<typeof import("vue-router")>("vue-router")

  return {
    ...actual,
    useRouter: () => ({
      push,
    }),
    useRoute: () => route,
  }
})

beforeEach(() => {
  vi.clearAllMocks()
  route.query.redirect = "/arena"
})

test("redirects to query redirect after guest login succeeds", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)

  vi.mocked(apiRequest).mockResolvedValue({
    player_id: 3001,
    token: "guest-token",
    nickname: "深链玩家",
  })

  const wrapper = mount(LoginPage, {
    global: {
      plugins: [pinia],
    },
  })

  await wrapper.get("#nickname").setValue("深链玩家")
  await wrapper.get("form").trigger("submit.prevent")
  await flushPromises()

  const sessionStore = useSessionStore()
  expect(apiRequest).toHaveBeenCalledWith("auth/guest-login", {
    method: "POST",
    body: JSON.stringify({
      nickname: "深链玩家",
    }),
  })
  expect(sessionStore.token).toBe("guest-token")
  expect(sessionStore.playerId).toBe(3001)
  expect(sessionStore.nickname).toBe("深链玩家")
  expect(push).toHaveBeenCalledWith("/arena")
})
