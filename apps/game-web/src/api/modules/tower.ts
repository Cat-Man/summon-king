import { apiRequest } from "@/api/http"

export type TowerKind = "pagoda" | "spirit"

export type TowerStatus = {
  tower: TowerKind
  label: string
  current_floor: number
  max_floor: number
  remaining_challenges: number
  reward_preview: string
}

export type TowerChallengeResult = {
  player_id: number
  tower: TowerKind
  floor: number
  reward: string
  remaining_challenges: number
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
