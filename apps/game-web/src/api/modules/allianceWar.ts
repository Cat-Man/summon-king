import { apiRequest } from "@/api/http"

export type AllianceWarIndex = {
  has_alliance: boolean
  current_role: string
  phase: string
  target_label: string
  can_register: boolean
  available_targets: string[]
}

export function getAllianceWarIndex(playerId: number) {
  return apiRequest<AllianceWarIndex>(`alliance-war/index?player_id=${playerId}`)
}

export function registerAllianceWarTarget(playerId: number, targetLabel: string) {
  return apiRequest<AllianceWarIndex>("alliance-war/register-target", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      target: targetLabel,
    }),
  })
}
