import { defineStore } from 'pinia'

export const useSessionStore = defineStore('session', {
  state: () => ({
    token: '',
    playerId: 0
  }),
  actions: {
    setToken(token: string) {
      this.token = token
    },
    setPlayerId(playerId: number) {
      this.playerId = playerId
    }
  }
})
