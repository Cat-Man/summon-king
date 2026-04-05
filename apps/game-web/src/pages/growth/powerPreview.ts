import type { BattlePet, PetCollection } from "@/api/modules/pet"

export const boneBonusPerLevel = 24
export const soulBonusPerPiece = 8

type BreakdownKey = "bone" | "spirit" | "soul"

export const emptyPetCollection: PetCollection = {
  player_id: 0,
  total_power: 0,
  team_size: 0,
  active_team: [],
  roster: [],
}

export function sumActiveTeamBonus(activeTeam: BattlePet[], key: BreakdownKey) {
  return activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.[key] ?? 0), 0)
}

export function predictNextTotalPower(collection: PetCollection, delta: number) {
  if (collection.active_team.length === 0) {
    return collection.total_power
  }
  return collection.total_power + delta
}
