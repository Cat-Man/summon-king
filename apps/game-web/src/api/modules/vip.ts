import { apiRequest } from "@/api/http"

type RewardDelta = {
  spirit_power?: number
  bone_level?: number
  soul_pieces?: number
}

type WalletSnapshot = {
  spirit_power?: number
  bone_level?: number
  soul_pieces?: number
}

type BackendVIPIndex = {
  player_id: number
  vip_level: number
  next_vip_level: number
  daily_claimed: boolean
  daily_reward: RewardDelta
  current_benefits: string[]
  next_benefits: string[]
  wallet_snapshot?: WalletSnapshot
}

type BackendVIPClaim = BackendVIPIndex & {
  reward_delta: RewardDelta
}

export type VIPReward = {
  spirit_power: number
  bone_level: number
  soul_pieces: number
}

export type VIPIndex = {
  player_id: number
  vip_level: number
  total_gem_spent: number
  daily_claimed: boolean
  daily_reward: VIPReward
  current_benefits: string[]
  next_level: {
    level: number
    required_gem_spent: number
    benefits: string[]
  }
  wallet_snapshot: {
    spirit_power: number
    bone_level: number
    soul_pieces: number
  }
}

export type VIPDailyClaimResult = {
  index: VIPIndex
  reward_delta: VIPReward
}

function normalizeReward(delta?: RewardDelta): VIPReward {
  return {
    spirit_power: delta?.spirit_power ?? 0,
    bone_level: delta?.bone_level ?? 0,
    soul_pieces: delta?.soul_pieces ?? 0,
  }
}

function normalizeWallet(wallet?: WalletSnapshot) {
  return {
    spirit_power: wallet?.spirit_power ?? 0,
    bone_level: wallet?.bone_level ?? 0,
    soul_pieces: wallet?.soul_pieces ?? 0,
  }
}

function normalizeIndex(raw: BackendVIPIndex): VIPIndex {
  return {
    player_id: raw.player_id,
    vip_level: raw.vip_level,
    total_gem_spent: raw.vip_level * 100,
    daily_claimed: raw.daily_claimed,
    daily_reward: normalizeReward(raw.daily_reward),
    current_benefits: raw.current_benefits ?? [],
    next_level: {
      level: raw.next_vip_level,
      required_gem_spent: raw.next_vip_level * 100,
      benefits: raw.next_benefits ?? [],
    },
    wallet_snapshot: normalizeWallet(raw.wallet_snapshot),
  }
}

export function getVIPIndex(playerId: number) {
  return apiRequest<BackendVIPIndex>(`vip/index?player_id=${playerId}`).then(normalizeIndex)
}

export function claimVIPDaily(playerId: number) {
  return apiRequest<BackendVIPClaim>("vip/claim-daily", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
    }),
  }).then((raw) => ({
    index: normalizeIndex(raw),
    reward_delta: normalizeReward(raw.reward_delta),
  }))
}
