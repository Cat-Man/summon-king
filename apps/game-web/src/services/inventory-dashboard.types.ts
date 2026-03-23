export interface InventoryDashboardHero {
  eyebrow: string
  title: string
  description: string
  tone: 'amber' | 'blue' | 'navy' | 'violet' | 'cyan' | 'green' | 'teal' | 'orange'
  metaLabel: string
  metaValue: string
}

export interface InventoryWalletStat {
  label: string
  value: string
}

export interface InventoryWalletResource {
  label: string
  value: string
  icon: string
}

export interface InventoryItemCard {
  itemId: string
  name: string
  count: string
  action: string
  sellPrice: number
  isHighValue: boolean
  icon?: string
}

export interface InventoryLogItem {
  id: string
  createdAt: string
  changeType: string
  changeTypeText: string
  reason: string
  coinsDelta: number
  diamondsDelta: number
  title: string
  delta: string
  createdAtLabel: string
}

export type InventoryLogEntry = InventoryLogItem

export interface InventoryDashboardData {
  hero: InventoryDashboardHero
  wallet: InventoryWalletStat[]
  walletResources: InventoryWalletResource[]
  items: InventoryItemCard[]
  rules: string[]
}
