import { request } from '@/api/http'
import { legacyPetIconMap } from '@/assets/legacy'
import { runtimeConfig, type GameDataSource } from '@/config/runtime'

import type {
  PetCatalogCard,
  PetCatalogDashboardData,
  PetDetailDashboardData,
  PetTeamDashboardData,
  PetTeamMember,
  PetTeamRosterItem,
  PetTeamSaveResult
} from './pet-dashboard.types'

interface SessionStoreLike {
  ensureGuestSession(): Promise<{ player_id: number; token: string }>
}

interface PetCatalogResponseItem {
  pet_id: number
  name: string
  rarity: number
  element: string
  locked: boolean
  owned: boolean
  power: number
  portrait: string
  story_text: string
}

interface PlayerPetResponseItem {
  player_id: number
  pet_id: number
  level: number
  star: number
  power: number
}

interface PetTeamResponse {
  player_id: number
  pet_ids: number[]
}

interface LoadPetOptions {
  dataSource?: GameDataSource
  sessionStore: SessionStoreLike
}

interface LoadPetDetailOptions extends LoadPetOptions {
  petId?: number
}

interface SavePetTeamOptions extends LoadPetOptions {
  petIds: number[]
}

interface PetMeta {
  displayName: string
  map: string
  role: string
  skillPool: string
  source: string
  icon?: string
}

const petMetaMap: Record<number, PetMeta> = {
  1: {
    displayName: '烈焰狼王',
    map: '青木林地',
    role: '主战输出',
    skillPool: '烈焰撕咬 / 焚风追击',
    source: '地图副本 / 召唤球',
    icon: legacyPetIconMap['烈焰狼王']
  },
  2: {
    displayName: '寒枝鹿灵',
    map: '北境冰湖',
    role: '控制辅助',
    skillPool: '寒枝缠绕 / 冰露护体',
    source: '地图副本 / 修行',
    icon: legacyPetIconMap['寒枝鹿灵']
  },
  3: {
    displayName: '雷角牛',
    map: '惊雷峡谷',
    role: '前排承伤',
    skillPool: '雷角突袭 / 震地怒吼',
    source: '地图副本 / 召唤球'
  }
}

function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value)
}

function formatElement(element: string): string {
  switch (element) {
    case 'fire':
      return '火系'
    case 'water':
      return '水系'
    case 'thunder':
      return '雷系'
    default:
      return element
  }
}

function formatRarity(rarity: number): string {
  return `${rarity} 星`
}

function getPetMeta(petId: number, fallbackName: string): PetMeta {
  return petMetaMap[petId] ?? {
    displayName: fallbackName,
    map: '未知地图',
    role: '待定定位',
    skillPool: '待补充技能池',
    source: '后续补充',
    icon: legacyPetIconMap[fallbackName]
  }
}

function createMockPetCatalogDashboard(): PetCatalogDashboardData {
  return {
    hero: {
      eyebrow: '图鉴与来源导航',
      title: '幻兽图鉴',
      description: '区分已拥有与未拥有幻兽，帮助玩家快速判断下一步去哪里刷、缺什么来源、哪些技能值得追。',
      tone: 'violet',
      metaLabel: '已拥有',
      metaValue: '3 / 3'
    },
    overview: [
      { label: '已拥有', value: '3' },
      { label: '未拥有', value: '0' },
      { label: '队伍推荐', value: '雷角牛' },
      { label: '获取来源', value: '副本 / 修行 / 礼包' }
    ],
    sources: [
      '地图副本：按地图掉落对应幻兽与召唤球',
      '修行：2 小时及以上有概率获得对应地图召唤球',
      '礼包与商城：用于补齐当前养成断档'
    ],
    pets: [
      {
        petId: 1,
        name: '烈焰狼王',
        status: '已拥有',
        map: '青木林地',
        skillPool: '烈焰撕咬 / 焚风追击',
        aptitude: '战力 120 / 星级 1 星',
        source: '地图副本 / 召唤球',
        level: 'Lv.12',
        power: '120',
        icon: legacyPetIconMap['烈焰狼王']
      },
      {
        petId: 2,
        name: '寒枝鹿灵',
        status: '已拥有',
        map: '北境冰湖',
        skillPool: '寒枝缠绕 / 冰露护体',
        aptitude: '战力 150 / 星级 2 星',
        source: '地图副本 / 修行',
        level: 'Lv.16',
        power: '150',
        icon: legacyPetIconMap['寒枝鹿灵']
      }
    ]
  }
}

export function createInitialPetCatalogDashboard(): PetCatalogDashboardData {
  return {
    hero: {
      eyebrow: '图鉴与来源导航',
      title: '幻兽图鉴',
      description: '',
      tone: 'violet',
      metaLabel: '已拥有',
      metaValue: ''
    },
    overview: [],
    sources: [],
    pets: []
  }
}

function createMockPetDetailDashboard(): PetDetailDashboardData {
  return {
    hero: {
      eyebrow: '养成中枢',
      title: '烈焰狼王',
      description: '把技能、战骨、战灵、魔魂、进化和重生入口集中展示，作为当前主养成幻兽的总控制台。',
      tone: 'orange',
      metaLabel: '综合战力',
      metaValue: '120'
    },
    overview: [
      { label: '等级', value: 'Lv.12' },
      { label: '元素', value: '火系' },
      { label: '星级', value: '1 星' },
      { label: '综合战力', value: '120' }
    ],
    skills: ['烈焰狼王主战输出技能池已准备'],
    bones: ['头骨 Lv.4', '胸骨 Lv.4', '臂骨 Lv.3'],
    spirits: ['火灵·暴怒'],
    souls: ['龙魂·灼炎'],
    sources: ['火山边缘的狩猎者'],
    heroIcon: legacyPetIconMap['烈焰狼王'],
    evolutionCost: {
      itemName: '火系进化石',
      itemProgress: '6 / 8',
      coinProgress: '180,000 / 200,000'
    }
  }
}

export function createInitialPetDetailDashboard(): PetDetailDashboardData {
  return {
    hero: {
      eyebrow: '养成中枢',
      title: '',
      description: '',
      tone: 'orange',
      metaLabel: '综合战力',
      metaValue: ''
    },
    overview: [],
    skills: [],
    bones: [],
    spirits: [],
    souls: [],
    sources: [],
    heroIcon: '',
    evolutionCost: {
      itemName: '',
      itemProgress: '',
      coinProgress: ''
    }
  }
}

function createMockPetTeamDashboard(): PetTeamDashboardData {
  return {
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '先展示当前战斗队，再展示幻兽栏与主养成目标。',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: '270'
    },
    overview: [
      { label: '已上阵', value: '2 / 5' },
      { label: '队伍核心', value: '寒枝鹿灵' },
      { label: '平均等级', value: 'Lv.14' },
      { label: '可替补', value: '1 只' }
    ],
    team: [
      { petId: 2, slot: '1 号位', name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.16', power: '150', icon: legacyPetIconMap['寒枝鹿灵'] },
      { petId: 1, slot: '2 号位', name: '烈焰狼王', role: '主战输出', level: 'Lv.12', power: '120', icon: legacyPetIconMap['烈焰狼王'] }
    ],
    roster: [
      { petId: 2, name: '寒枝鹿灵', level: 'Lv.16', status: '已上阵', power: '150', icon: legacyPetIconMap['寒枝鹿灵'] },
      { petId: 1, name: '烈焰狼王', level: 'Lv.12', status: '已上阵', power: '120', icon: legacyPetIconMap['烈焰狼王'] },
      { petId: 3, name: '雷角牛', level: 'Lv.20', status: '可替补', power: '180' }
    ],
    strategies: ['保存阵容', '支持上下阵', '支持排序调位', '按综合战力筛选'],
    focus: {
      name: '寒枝鹿灵',
      description: '当前控制与辅助能力最稳定，适合作为队伍核心。',
      nextStep: '继续补齐副本与修行材料。'
    }
  }
}

export function createInitialPetTeamDashboard(): PetTeamDashboardData {
  return {
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: ''
    },
    overview: [],
    team: [],
    roster: [],
    strategies: [],
    focus: {
      name: '',
      description: '',
      nextStep: ''
    }
  }
}

function withAuthHeaders(token: string) {
  return token ? { Authorization: `Bearer ${token}` } : undefined
}

async function fetchPetCatalog(session: { player_id: number; token: string }) {
  const headers = withAuthHeaders(session.token)
  const [catalog, playerPets] = await Promise.all([
    request<PetCatalogResponseItem[]>(`/player/pets/catalog?player_id=${session.player_id}`, { headers }),
    request<PlayerPetResponseItem[]>(`/player/pets/list?player_id=${session.player_id}`, { headers })
  ])
  return { catalog, playerPets }
}

function buildCatalogCards(catalog: PetCatalogResponseItem[], playerPets: PlayerPetResponseItem[]): PetCatalogCard[] {
  const playerPetMap = new Map(playerPets.map((item) => [item.pet_id, item]))

  return catalog.map((item) => {
    const playerPet = playerPetMap.get(item.pet_id)
    const meta = getPetMeta(item.pet_id, item.name)
    return {
      petId: item.pet_id,
      name: meta.displayName,
      status: item.owned ? '已拥有' : '未拥有',
      map: meta.map,
      skillPool: meta.skillPool,
      aptitude: `战力 ${formatNumber(playerPet?.power ?? item.power)} / 星级 ${formatRarity(playerPet?.star ?? item.rarity)}`,
      source: meta.source,
      level: playerPet ? `Lv.${playerPet.level}` : '-',
      power: formatNumber(playerPet?.power ?? item.power),
      icon: meta.icon
    }
  })
}

export async function loadPetCatalogDashboard({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore
}: LoadPetOptions): Promise<PetCatalogDashboardData> {
  if (dataSource === 'mock') {
    return createMockPetCatalogDashboard()
  }

  const session = await sessionStore.ensureGuestSession()
  const { catalog, playerPets } = await fetchPetCatalog(session)
  const pets = buildCatalogCards(catalog, playerPets)
  const ownedCount = catalog.filter((item) => item.owned).length
  const recommended = [...pets].sort((left, right) => Number(right.power.replace(/,/g, '')) - Number(left.power.replace(/,/g, '')))[0]

  return {
    hero: {
      eyebrow: '图鉴与来源导航',
      title: '幻兽图鉴',
      description: '区分已拥有与未拥有幻兽，帮助玩家快速判断下一步去哪里刷、缺什么来源、哪些技能值得追。',
      tone: 'violet',
      metaLabel: '已拥有',
      metaValue: `${ownedCount} / ${catalog.length}`
    },
    overview: [
      { label: '已拥有', value: `${ownedCount}` },
      { label: '未拥有', value: `${catalog.length - ownedCount}` },
      { label: '队伍推荐', value: recommended?.name ?? '-' },
      { label: '获取来源', value: '副本 / 修行 / 礼包' }
    ],
    sources: [
      '地图副本：按地图掉落对应幻兽与召唤球',
      '修行：2 小时及以上有概率获得对应地图召唤球',
      '礼包与商城：用于补齐当前养成断档'
    ],
    pets
  }
}

export async function loadPetDetailDashboard({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore,
  petId
}: LoadPetDetailOptions): Promise<PetDetailDashboardData> {
  if (dataSource === 'mock') {
    return createMockPetDetailDashboard()
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = withAuthHeaders(session.token)
  const playerPets = await request<PlayerPetResponseItem[]>(`/player/pets/list?player_id=${session.player_id}`, { headers })
  const targetPetId = petId ?? playerPets[0]?.pet_id ?? 1
  const [detail, playerPet] = await Promise.all([
    request<PetCatalogResponseItem>(`/player/pets/detail?player_id=${session.player_id}&pet_id=${targetPetId}`, { headers }),
    Promise.resolve(playerPets.find((item) => item.pet_id === targetPetId))
  ])

  const meta = getPetMeta(detail.pet_id, detail.name)
  const displayName = meta.displayName
  const level = playerPet?.level ?? 1
  const star = playerPet?.star ?? detail.rarity
  const power = playerPet?.power ?? detail.power

  return {
    hero: {
      eyebrow: '养成中枢',
      title: displayName,
      description: `把技能、战骨、战灵、魔魂、进化和重生入口集中展示，作为 ${displayName} 的总控制台。`,
      tone: detail.element === 'fire' ? 'orange' : detail.element === 'water' ? 'blue' : 'green',
      metaLabel: '综合战力',
      metaValue: formatNumber(power)
    },
    overview: [
      { label: '等级', value: `Lv.${level}` },
      { label: '元素', value: formatElement(detail.element) },
      { label: '星级', value: `${star} 星` },
      { label: '综合战力', value: formatNumber(power) }
    ],
    skills: [
      `${displayName}技能池：${meta.skillPool}`,
      `${displayName}当前故事：${detail.story_text}`
    ],
    bones: [`头骨 Lv.${Math.max(3, star + 2)}`, `胸骨 Lv.${Math.max(3, star + 2)}`, `臂骨 Lv.${Math.max(2, star + 1)}`],
    spirits: [`${formatElement(detail.element)}灵·共鸣`],
    souls: [`${formatElement(detail.element)}魂·加护`],
    sources: [detail.story_text, `获取来源：${meta.source}`],
    heroIcon: meta.icon,
    evolutionCost: {
      itemName: detail.element === 'fire' ? '火系进化石' : '进化石',
      itemProgress: `${Math.min(8, star + 4)} / 8`,
      coinProgress: `${formatNumber(power * 1200)} / ${formatNumber(200000)}`
    }
  }
}

function buildTeamMembers(playerPets: PlayerPetResponseItem[], petIds: number[]): PetTeamMember[] {
  const petMap = new Map(playerPets.map((item) => [item.pet_id, item]))
  return petIds.flatMap((petId, index) => {
    const pet = petMap.get(petId)
    if (!pet) {
      return []
    }
    const meta = getPetMeta(petId, `${petId}`)
    const member: PetTeamMember = {
      petId,
      slot: `${index + 1} 号位`,
      name: meta.displayName,
      role: meta.role,
      level: `Lv.${pet.level}`,
      power: formatNumber(pet.power),
      icon: meta.icon
    }
    return [member]
  })
}

function buildRoster(playerPets: PlayerPetResponseItem[], teamPetIds: number[]): PetTeamRosterItem[] {
  const teamSet = new Set(teamPetIds)
  return [...playerPets]
    .sort((left, right) => right.power - left.power)
    .map((pet) => {
      const meta = getPetMeta(pet.pet_id, `${pet.pet_id}`)
      return {
        petId: pet.pet_id,
        name: meta.displayName,
        level: `Lv.${pet.level}`,
        status: teamSet.has(pet.pet_id) ? '已上阵' : '可替补',
        power: formatNumber(pet.power),
        icon: meta.icon
      }
    })
}

export async function loadPetTeamDashboard({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore
}: LoadPetOptions): Promise<PetTeamDashboardData> {
  if (dataSource === 'mock') {
    return createMockPetTeamDashboard()
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = withAuthHeaders(session.token)
  const [playerPets, team] = await Promise.all([
    request<PlayerPetResponseItem[]>(`/player/pets/list?player_id=${session.player_id}`, { headers }),
    request<PetTeamResponse>(`/player/pets/team?player_id=${session.player_id}`, { headers })
  ])

  const teamMembers = buildTeamMembers(playerPets, team.pet_ids)
  const roster = buildRoster(playerPets, team.pet_ids)
  const totalPower = teamMembers.reduce((sum, item) => sum + Number(item.power.replace(/,/g, '')), 0)
  const averageLevel = teamMembers.length
    ? Math.round(
        teamMembers.reduce((sum, item) => sum + Number(item.level.replace(/[^\d]/g, '')), 0) / teamMembers.length
      )
    : 0
  const focus = teamMembers[0] ?? roster[0]

  return {
    hero: {
      eyebrow: '战斗队与幻兽栏',
      title: '主力阵容',
      description: '先展示当前战斗队，再展示幻兽栏与主养成目标。',
      tone: 'blue',
      metaLabel: '综合战力',
      metaValue: formatNumber(totalPower)
    },
    overview: [
      { label: '已上阵', value: `${teamMembers.length} / 5` },
      { label: '队伍核心', value: focus?.name ?? '-' },
      { label: '平均等级', value: averageLevel > 0 ? `Lv.${averageLevel}` : '-' },
      { label: '可替补', value: `${Math.max(0, roster.length - teamMembers.length)} 只` }
    ],
    team: teamMembers,
    roster,
    strategies: ['保存阵容', '支持上下阵', '支持排序调位', '按综合战力筛选'],
    focus: {
      name: focus?.name ?? '-',
      description: `${focus?.name ?? '当前主力'}当前战力最高，建议优先投入养成资源。`,
      nextStep: '保存阵容后可用于首页战力与战斗编队联动。'
    }
  }
}

export async function savePetTeamSelection({
  dataSource = runtimeConfig.gameDataSource,
  sessionStore,
  petIds
}: SavePetTeamOptions): Promise<PetTeamSaveResult> {
  if (dataSource === 'mock') {
    return { message: '阵容已保存' }
  }

  const session = await sessionStore.ensureGuestSession()
  const headers = withAuthHeaders(session.token)
  await request<{ saved: boolean }>('/player/pets/team/save', {
    method: 'POST',
    headers,
    body: JSON.stringify({
      player_id: session.player_id,
      pet_ids: petIds
    })
  })

  return { message: '阵容已保存' }
}
