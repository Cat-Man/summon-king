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

export function setMainPet(playerId: number, petId: number) {
  return apiRequest<PetCollection>("pet/team/main", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      pet_id: petId,
    }),
  })
}

export function savePetTeam(playerId: number, petIds: number[]) {
  return apiRequest<PetCollection>("pet/team/save", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
      pet_ids: petIds,
    }),
  })
}
