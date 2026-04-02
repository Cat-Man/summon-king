import { computed, ref } from "vue"
import { defineStore } from "pinia"

type SessionSnapshot = {
  token: string
  playerId: number | null
  nickname: string
}

const storageKey = "zhzw-session"

function readInitialSession(): SessionSnapshot {
  if (typeof window === "undefined") {
    return { token: "", playerId: null, nickname: "" }
  }

  const raw = window.sessionStorage.getItem(storageKey)
  if (!raw) {
    return { token: "", playerId: null, nickname: "" }
  }

  try {
    const parsed = JSON.parse(raw) as Partial<SessionSnapshot>
    return {
      token: parsed.token ?? "",
      playerId: parsed.playerId ?? null,
      nickname: parsed.nickname ?? "",
    }
  } catch {
    return { token: "", playerId: null, nickname: "" }
  }
}

export const useSessionStore = defineStore("session", () => {
  const initial = readInitialSession()
  const token = ref(initial.token)
  const playerId = ref<number | null>(initial.playerId)
  const nickname = ref(initial.nickname)

  const isLoggedIn = computed(() => token.value.length > 0 && playerId.value !== null)

  function persist() {
    if (typeof window === "undefined") {
      return
    }
    window.sessionStorage.setItem(
      storageKey,
      JSON.stringify({
        token: token.value,
        playerId: playerId.value,
        nickname: nickname.value,
      } satisfies SessionSnapshot),
    )
  }

  function setSession(next: SessionSnapshot) {
    token.value = next.token
    playerId.value = next.playerId
    nickname.value = next.nickname
    persist()
  }

  function clearSession() {
    token.value = ""
    playerId.value = null
    nickname.value = ""
    if (typeof window !== "undefined") {
      window.sessionStorage.removeItem(storageKey)
    }
  }

  return {
    token,
    playerId,
    nickname,
    isLoggedIn,
    setSession,
    clearSession,
  }
})
