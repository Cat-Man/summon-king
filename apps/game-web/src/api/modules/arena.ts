import { apiRequest } from "@/api/http"

export type ArenaRecord = {
  player_id: number
  current_streak: number
  last_win: boolean
}

export function getArenaStatus(playerId: number) {
  return apiRequest<ArenaRecord>(`arena/status?player_id=${playerId}`)
}

export function challengeArena(playerId: number, won: boolean) {
  return apiRequest<ArenaRecord>("arena/challenge", {
    method: "POST",
    body: JSON.stringify({ player_id: playerId, won }),
  })
}
