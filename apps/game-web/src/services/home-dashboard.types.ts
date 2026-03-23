export interface HomeDashboardHero {
  eyebrow: string
  title: string
  description: string
  tone: 'navy' | 'blue' | 'violet' | 'amber' | 'cyan' | 'green' | 'teal' | 'orange'
  metaLabel: string
  metaValue: string
}

export interface HomeDashboardTodo {
  title: string
  value: string
  action: string
}

export interface HomeDashboardResource {
  label: string
  value: string
}

export interface HomeDashboardResourceIcon {
  label: string
  value: string
  icon: string
}

export interface HomeDashboardActivityEntry {
  title: string
  description: string
  note: string
  icon: string
}

export interface HomeDashboardCultivationSummary {
  action: string
  items: Array<{ label: string; value: string; subtext?: string }>
}

export interface HomeDashboardData {
  hero: HomeDashboardHero
  dailyTodos: HomeDashboardTodo[]
  resources: HomeDashboardResource[]
  resourceIcons: HomeDashboardResourceIcon[]
  activityEntry: HomeDashboardActivityEntry
  cultivationSummary: HomeDashboardCultivationSummary
  entries: string[]
  messages: string[]
}
