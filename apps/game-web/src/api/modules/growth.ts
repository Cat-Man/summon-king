import { apiRequest } from "@/api/http"

export type GrowthWallet = {
  player_id: number
  spirit_power: number
  spirit_free_wash: number
  bone_level: number
  soul_pieces: number
  manor_plots: number
}

export type BoneState = {
  name: string
  level: number
}

export type SoulState = {
  name: string
  power: number
}

export type ManorPlot = {
  plot_id: number
  state: string
}

export type WashSpiritResult = {
  washed: boolean
}

export type ManorHarvestResult = {
  message: string
  plots: ManorPlot[]
}

export function getGrowthWallet(playerId: number) {
  return apiRequest<GrowthWallet>(`growth/wallet?player_id=${playerId}`)
}

export function washSpirit(playerId: number) {
  return apiRequest<WashSpiritResult>(`growth/spirit/wash?player_id=${playerId}`, {
    method: "POST",
    body: JSON.stringify({}),
  })
}

export function getBoneState(playerId: number) {
  return apiRequest<BoneState>(`growth/bone?player_id=${playerId}`)
}

export function upgradeBone(playerId: number) {
  return apiRequest<BoneState>(`growth/bone/upgrade?player_id=${playerId}`, {
    method: "POST",
    body: JSON.stringify({}),
  })
}

export function getSoulState(playerId: number) {
  return apiRequest<SoulState>(`growth/soul?player_id=${playerId}`)
}

export function getManorPlots(playerId: number) {
  return apiRequest<ManorPlot[]>(`growth/manor?player_id=${playerId}`)
}

export function harvestManor(playerId: number) {
  return apiRequest<ManorHarvestResult>(`growth/manor/harvest?player_id=${playerId}`, {
    method: "POST",
    body: JSON.stringify({}),
  })
}
