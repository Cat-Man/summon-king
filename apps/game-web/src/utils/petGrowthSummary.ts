import type { PetCollection } from "@/api/modules/pet"

export type PetGrowthSummary = {
  teamPower: number
  boneBonus: number
  spiritBonus: number
  soulBonus: number
  totalBonus: number
}

export function summarizePetGrowth(collection: PetCollection | null | undefined): PetGrowthSummary {
  const activeTeam = collection?.active_team ?? []
  const boneBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.bone ?? 0), 0)
  const spiritBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.spirit ?? 0), 0)
  const soulBonus = activeTeam.reduce((total, pet) => total + (pet.power_breakdown?.soul ?? 0), 0)

  return {
    teamPower: collection?.total_power ?? 0,
    boneBonus,
    spiritBonus,
    soulBonus,
    totalBonus: boneBonus + spiritBonus + soulBonus,
  }
}
