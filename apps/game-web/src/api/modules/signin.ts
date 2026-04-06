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

type BackendSigninIndex = {
  player_id: number
  today_claimed: boolean
  streak_days: number
  today_reward: RewardDelta
  wallet_snapshot?: WalletSnapshot
}

type BackendSigninClaim = {
  player_id: number
  today_claimed: boolean
  streak_days: number
  reward_delta: RewardDelta
  wallet_snapshot?: WalletSnapshot
}

export type SigninReward = {
  spirit_power: number
  bone_level: number
  soul_pieces: number
}

export type SigninIndex = {
  player_id: number
  signed_today: boolean
  current_streak: number
  next_reward: SigninReward
  wallet_snapshot: {
    spirit_power: number
    bone_level: number
    soul_pieces: number
  }
}

export type SigninClaimResult = {
  index: SigninIndex
  reward_delta: SigninReward
}

function normalizeReward(delta?: RewardDelta): SigninReward {
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

function normalizeIndex(raw: BackendSigninIndex): SigninIndex {
  return {
    player_id: raw.player_id,
    signed_today: raw.today_claimed,
    current_streak: raw.streak_days,
    next_reward: normalizeReward(raw.today_reward),
    wallet_snapshot: normalizeWallet(raw.wallet_snapshot),
  }
}

export function getSigninIndex(playerId: number) {
  return apiRequest<BackendSigninIndex>(`signin/index?player_id=${playerId}`).then(normalizeIndex)
}

export function claimSignin(playerId: number) {
  return apiRequest<BackendSigninClaim>("signin/claim", {
    method: "POST",
    body: JSON.stringify({
      player_id: playerId,
    }),
  }).then((raw) => ({
    index: {
      player_id: raw.player_id,
      signed_today: raw.today_claimed,
      current_streak: raw.streak_days,
      next_reward: normalizeReward(raw.reward_delta),
      wallet_snapshot: normalizeWallet(raw.wallet_snapshot),
    },
    reward_delta: normalizeReward(raw.reward_delta),
  }))
}
