import { computed, ref } from "vue"
import { defineStore } from "pinia"

export type SessionSnapshot = {
  token: string
  playerId: number | null
  nickname: string
}

export const sessionStorageKey = "zhzw-session"

export function readSessionSnapshot(): SessionSnapshot {
  if (typeof window === "undefined") {
    return { token: "", playerId: null, nickname: "" }
  }

  const raw = window.sessionStorage.getItem(sessionStorageKey)
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
  const initial = readSessionSnapshot()
  const token = ref(initial.token)
  const playerId = ref<number | null>(initial.playerId)
  const nickname = ref(initial.nickname)

  const isLoggedIn = computed(() => token.value.length > 0 && playerId.value !== null)

  function persist() {
    if (typeof window === "undefined") {
      return
    }
    window.sessionStorage.setItem(
      sessionStorageKey,
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
      window.sessionStorage.removeItem(sessionStorageKey)
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
