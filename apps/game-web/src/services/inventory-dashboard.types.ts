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
  icon?: string
}

export interface InventoryDashboardData {
  hero: InventoryDashboardHero
  wallet: InventoryWalletStat[]
  walletResources: InventoryWalletResource[]
  items: InventoryItemCard[]
  rules: string[]
}
