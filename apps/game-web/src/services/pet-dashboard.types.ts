import type { HomeDashboardHero } from './home-dashboard.types'

export interface PetCatalogCard {
  petId: number
  name: string
  status: string
  map: string
  skillPool: string
  aptitude: string
  source: string
  level: string
  power: string
  icon?: string
}

export interface PetCatalogDashboardData {
  hero: HomeDashboardHero
  overview: Array<{ label: string; value: string }>
  sources: string[]
  pets: PetCatalogCard[]
}

export interface PetDetailDashboardData {
  hero: HomeDashboardHero
  overview: Array<{ label: string; value: string }>
  skills: string[]
  bones: string[]
  spirits: string[]
  souls: string[]
  sources: string[]
  heroIcon?: string
  evolutionCost: {
    itemName: string
    itemProgress: string
    coinProgress: string
  }
}

export interface PetTeamMember {
  petId: number
  slot: string
  name: string
  role: string
  level: string
  power: string
  icon?: string
}

export interface PetTeamRosterItem {
  petId: number
  name: string
  role: string
  level: string
  status: string
  power: string
  icon?: string
}

export interface PetTeamDashboardData {
  hero: HomeDashboardHero
  overview: Array<{ label: string; value: string }>
  team: PetTeamMember[]
  roster: PetTeamRosterItem[]
  strategies: string[]
  focus: {
    name: string
    description: string
    nextStep: string
  }
}

export interface PetTeamSaveResult {
  message: string
}
