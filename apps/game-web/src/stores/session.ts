import { defineStore } from 'pinia'

import { request } from '@/api/http'

export interface LoginResponse {
  player_id: number
  token: string
  channel: string
}

export interface HomeIndexResponse {
  player_id: number
  nickname: string
  level: number
  coin: number
  diamond: number
  last_login_at: string
}

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
    },
    setSession(playerId: number, token: string) {
      this.playerId = playerId
      this.token = token
    },
    async loginGuest(channel = 'web') {
      const response = await request<LoginResponse>('/player/auth/login', {
        method: 'POST',
        body: JSON.stringify({ channel })
      })
      this.setSession(response.player_id, response.token)
      return response
    },
    async ensureGuestSession(channel = 'web') {
      if (this.playerId && this.token) {
        return {
          player_id: this.playerId,
          token: this.token,
          channel
        }
      }

      return this.loginGuest(channel)
    },
    async fetchHomeIndex() {
      if (!this.playerId) {
        throw new Error('player id is required')
      }
      const response = await request<HomeIndexResponse>(
        `/player/home/index?player_id=${this.playerId}`,
        {
          headers: this.token ? { Authorization: `Bearer ${this.token}` } : undefined
        }
      )
      return response
    }
  }
})
