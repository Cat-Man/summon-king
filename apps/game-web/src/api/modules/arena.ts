import { apiRequest } from "@/api/http"

export type ArenaRecord = {
  player_id: number
  current_streak: number
  last_win: boolean
}

export type ArenaRewardDelta = {
  spirit_power: number
  soul_pieces: number
}

export type ArenaWalletSnapshot = {
  spirit_power: number
  bone_level: number
  soul_pieces: number
}

export type ArenaBattleResult = {
  record: ArenaRecord
  reward_delta: ArenaRewardDelta
  wallet_snapshot: ArenaWalletSnapshot
}

export function getArenaStatus(playerId: number) {
  return apiRequest<ArenaRecord>(`arena/status?player_id=${playerId}`)
}

export function challengeArena(playerId: number, won: boolean) {
  return apiRequest<ArenaBattleResult>("arena/challenge", {
    method: "POST",
    body: JSON.stringify({ player_id: playerId, won }),
  })
}
