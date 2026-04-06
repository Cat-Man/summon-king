import { apiRequest } from "@/api/http"
import type { TowerRewardDelta } from "@/api/modules/tower"

export type Wallet = {
  player_id: number
  spirit_power: number
  spirit_free_wash: number
  bone_level: number
  soul_pieces: number
  manor_plots: number
}

export type TowerSummary = {
  current_floor: number
  remaining_challenges: number
  reward_preview: string
  last_reward?: string
  last_reward_delta?: TowerRewardDelta
}

export type TowerOverview = {
  pagoda: TowerSummary
  spirit: TowerSummary
}

export type ArenaOverview = {
  current_streak: number
  last_win: boolean
}

export type RankingOverview = {
  self_rank: number
  self_score: number
  top_name?: string
  top_score?: number
}

export type NextAction = {
  title: string
  description: string
  route: string
  cta: string
}

export type HomeOverview = {
  player_id: number
  nickname: string
  wallet: Wallet
  modules: {
    map_label: string
    map_city_count: number
    dungeon: {
      status: string
      current_floor: number
      remain_dice: number
      dungeon_id: number
    }
    cultivation: {
      state: string
      spirit_power: number
      claimable: boolean
      claimable_at?: string
    }
    pet: {
      total_power: number
      active_count: number
      starter_pet_name: string
    }
    tower: TowerOverview
    arena: ArenaOverview
    ranking: RankingOverview
  }
  next_action: NextAction
}

export function getHomeOverview(playerId: number) {
  return apiRequest<HomeOverview>(`home/overview?player_id=${playerId}`)
}
