import { apiRequest } from "@/api/http"

export type BattlePet = {
  pet_id: number
  slot: number
  name: string
  level: number
  power: number
  is_active: boolean
}

export type PetCollection = {
  player_id: number
  total_power: number
  team_size: number
  active_team: BattlePet[]
  roster: BattlePet[]
}

export function getPetCollection(playerId: number) {
  return apiRequest<PetCollection>(`pet/team?player_id=${playerId}`)
}
