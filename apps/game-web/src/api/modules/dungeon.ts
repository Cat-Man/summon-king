import { apiRequest } from "@/api/http"

export type WorldMap = {
  name: string
  cities: Array<{
    city_id: number
    name: string
    region: string
    loc_x: number
    loc_y: number
    dungeon_id: number
    dungeon_name: string
  }>
}

export type DungeonRun = {
  player_id: number
  dungeon_id: number
  remain_dice: number
  current_floor: number
  status: string
  started_at: string
  last_reward: {
    label: string
    spirit_power: number
    soul_pieces: number
  }
  wallet_snapshot: {
    player_id: number
    spirit_power: number
    spirit_free_wash: number
    bone_level: number
    soul_pieces: number
    manor_plots: number
  }
}

export type CultivationStatus = {
  player_id: number
  spirit_power: number
  state: string
  start_at?: string
  claimable_at?: string
}

function toFormBody(payload: Record<string, number>) {
  const params = new URLSearchParams()
  Object.entries(payload).forEach(([key, value]) => {
    params.set(key, String(value))
  })
  return params.toString()
}

function postForm<T>(path: string, payload: Record<string, number>) {
  return apiRequest<T>(path, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: toFormBody(payload),
  })
}

export function getWorldMap() {
  return apiRequest<WorldMap>("dungeon/overview")
}

export function getDungeonStatus(playerId: number) {
  return apiRequest<DungeonRun>(`dungeon/status?player_id=${playerId}`)
}

export function enterDungeon(playerId: number, dungeonId: number) {
  return postForm<DungeonRun>("dungeon/enter", {
    player_id: playerId,
    dungeon_id: dungeonId,
  })
}

export function rollDungeonDice(playerId: number) {
  return postForm<DungeonRun>("dungeon/roll", {
    player_id: playerId,
  })
}

export function startCultivation(playerId: number) {
  return postForm<CultivationStatus>("dungeon/cultivation/start", {
    player_id: playerId,
  })
}

export function claimCultivation(playerId: number) {
  return postForm<CultivationStatus>("dungeon/cultivation/claim", {
    player_id: playerId,
  })
}
