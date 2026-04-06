import { createPinia, setActivePinia } from "pinia"
import { flushPromises, mount } from "@vue/test-utils"

import LoginPage from "../LoginPage.vue"
import { apiRequest } from "@/api/http"
import { useSessionStore } from "@/stores/session"

const push = vi.fn()
const route = {
  query: {
    redirect: "/arena",
    channel: undefined as string | undefined,
    token: undefined as string | undefined,
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
  route.query.channel = undefined
  route.query.token = undefined
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

test("auto bootstraps wxmini session when channel and token are present", async () => {
  const pinia = createPinia()
  setActivePinia(pinia)
  route.query.redirect = "/home"
  route.query.channel = "wxmini"
  route.query.token = "wxmini-unified-token"

  vi.mocked(apiRequest).mockResolvedValue({
    player_id: 4001,
    token: "wxmini-session-token",
    nickname: "小程序玩家",
    channel: "wxmini",
  })

  mount(LoginPage, {
    global: {
      plugins: [pinia],
    },
  })
  await flushPromises()

  const sessionStore = useSessionStore()
  expect(apiRequest).toHaveBeenCalledWith("bridge/wxmini/session/bootstrap", {
    method: "POST",
    body: JSON.stringify({
      unified_token: "wxmini-unified-token",
    }),
  })
  expect(sessionStore.token).toBe("wxmini-session-token")
  expect(sessionStore.playerId).toBe(4001)
  expect(sessionStore.nickname).toBe("小程序玩家")
  expect(push).toHaveBeenCalledWith("/home")
})
