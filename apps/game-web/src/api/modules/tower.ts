import { apiRequest } from "@/api/http"

export type TowerKind = "pagoda" | "spirit"

export type TowerStatus = {
  tower: TowerKind
  label: string
  current_floor: number
  max_floor: number
  remaining_challenges: number
  reward_preview: string
  last_reward?: string
  last_reward_delta?: TowerRewardDelta
}

export type TowerRewardDelta = {
  bone_level?: number
  spirit_power?: number
  soul_pieces?: number
}

export type WalletSnapshot = {
  spirit_power: number
  bone_level: number
  soul_pieces: number
}

export type TowerBattleSummary = {
  battle_no: string
  battle_type: string
  result: string
  winner_side: string
  rounds: number
  attacker_power: number
  defender_power: number
}

export type TowerChallengeResult = {
  player_id: number
  tower: TowerKind
  floor: number
  reward: string
  remaining_challenges: number
  reward_delta: TowerRewardDelta
  battle: TowerBattleSummary
  wallet_snapshot: WalletSnapshot
}

function towerPath(tower: TowerKind) {
  return tower === "spirit" ? "tower/spirit-tower" : "tower/pagoda"
}

export function getTowerStatus(tower: TowerKind, playerId: number) {
  return apiRequest<TowerStatus>(`${towerPath(tower)}/status?player_id=${playerId}`)
}

export function startTowerChallenge(tower: TowerKind, playerId: number) {
  return apiRequest<TowerChallengeResult>(`${towerPath(tower)}/start`, {
    method: "POST",
    body: JSON.stringify({ player_id: playerId }),
  })
}
