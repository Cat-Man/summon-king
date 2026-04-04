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

export type ArenaBattleSummary = {
  battle_no: string
  battle_type: string
  result: string
  winner_side: string
  rounds: number
  attacker_power: number
  defender_power: number
}

export type ArenaBattleResult = {
  record: ArenaRecord
  reward_delta: ArenaRewardDelta
  battle: ArenaBattleSummary
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
