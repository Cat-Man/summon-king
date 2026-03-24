export type GameDataSource = 'mock' | 'api'

function normalizeApiBaseUrl(value?: string): string {
  const trimmed = value?.trim()
  if (!trimmed) {
    return '/api/v1'
  }

  return trimmed.replace(/\/$/, '')
}

function normalizeGameDataSource(value?: string): GameDataSource {
  return value === 'api' ? 'api' : 'mock'
}

export const runtimeConfig = {
  apiBaseUrl: normalizeApiBaseUrl(import.meta.env.VITE_API_BASE_URL),
  gameDataSource: normalizeGameDataSource(import.meta.env.VITE_GAME_DATA_SOURCE)
}
