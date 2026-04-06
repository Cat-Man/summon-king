import { apiRequest } from "@/api/http"

export type AllianceMember = {
  player_id: number
  role: string
}

export type AllianceBuilding = {
  building_type: string
  level: number
}

export type AllianceApplication = {
  player_id: number
}

export type AllianceSummary = {
  alliance_id: number
  name: string
  level: number
  notice: string
  member_count: number
  member_limit: number
  members: AllianceMember[]
  buildings: AllianceBuilding[]
}

export type AllianceIndex = {
  player_id: number
  has_alliance: boolean
  current_role: string
  alliance: AllianceSummary | null
  pending_applications?: AllianceApplication[]
}

export type AllianceHallEntry = {
  alliance_id: number
  name: string
  level: number
  member_count: number
  member_limit: number
  notice: string
  has_applied: boolean
}

export function getAllianceIndex(playerId: number) {
  return apiRequest<AllianceIndex>(`alliance/index?player_id=${playerId}`)
}

export function getAllianceHall(playerId: number) {
  return apiRequest<AllianceHallEntry[]>(`alliance/hall?player_id=${playerId}`)
}

export function createAlliance(playerId: number, name: string) {
  return apiRequest<AllianceIndex>("alliance/create", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      name,
    }),
  })
}

export function applyToAlliance(playerId: number, allianceId: number) {
  return apiRequest<{ status: string }>("alliance/apply", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      alliance_id: allianceId,
    }),
  })
}

export function approveAllianceApplication(playerId: number, applicantPlayerId: number) {
  return apiRequest<AllianceIndex>("alliance/apply/approve", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      applicant_player_id: applicantPlayerId,
    }),
  })
}
