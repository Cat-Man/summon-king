/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_GAME_DATA_SOURCE?: 'mock' | 'api'
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
