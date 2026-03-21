import { defineStore } from 'pinia'

import { request } from '@/api/http'

interface APIResponse<T> {
  code: number
  message: string
  data: T
  trace_id: string
}

interface LoginResponse {
  player_id: number
  token: string
  channel: string
}

interface HomeIndexResponse {
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
      const response = await request<APIResponse<LoginResponse>>('/api/v1/player/auth/login', {
        method: 'POST',
        body: JSON.stringify({ channel })
      })
      this.setSession(response.data.player_id, response.data.token)
      return response.data
    },
    async fetchHomeIndex() {
      if (!this.playerId) {
        throw new Error('player id is required')
      }
      const response = await request<APIResponse<HomeIndexResponse>>(
        `/api/v1/player/home/index?player_id=${this.playerId}`,
        {
          headers: this.token ? { Authorization: `Bearer ${this.token}` } : undefined
        }
      )
      return response.data
    }
  }
})
