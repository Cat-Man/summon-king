import { apiRequest } from "@/api/http"

export type LeaderboardEntry = {
  rank: number
  player_id: number
  name: string
  score: number
  updated: number
  is_self: boolean
}

export function getLeaderboard(playerId: number) {
  return apiRequest<LeaderboardEntry[]>(`ranking/leaderboard?player_id=${playerId}`)
}
