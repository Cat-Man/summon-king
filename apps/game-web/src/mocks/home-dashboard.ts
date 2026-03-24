import { legacyActivityAssets, legacyResourceIconMap } from '@/assets/legacy'

import type { HomeDashboardData } from '@/services/home-dashboard.types'

export function createMockHomeDashboard(): HomeDashboardData {
  return {
    hero: {
      eyebrow: '今日工作台',
      title: '召唤之王',
      description: '欢迎回来，游客1001。第一屏直接告诉玩家今天能做什么、有哪些收益可领、当前主推进目标是什么。',
      tone: 'navy',
      metaLabel: '当前角色',
      metaValue: '游客1001'
    },
    dailyTodos: [
      { title: '签到状态', value: '今日未签', action: '前往签到', actionKey: 'signin' },
      { title: '当前修行状态', value: '还有 18 分钟可领取', action: '查看修行', actionKey: 'cultivation' },
      { title: '当前推荐副本', value: '青木林地 · Boss 可挑战', action: '进入副本', actionKey: 'dungeon_run' }
    ],
    resources: [
      { label: '等级', value: 'Lv.1' },
      { label: '战力', value: '1,080' },
      { label: '活力', value: '120 / 120' },
      { label: '铜钱', value: '0' },
      { label: '元宝', value: '0' },
      { label: '声望', value: '0' }
    ],
    resourceIcons: [
      { label: '铜钱', value: '0', icon: legacyResourceIconMap['铜钱'] },
      { label: '元宝', value: '0', icon: legacyResourceIconMap['元宝'] }
    ],
    activityEntry: {
      title: '今日活动',
      description: '夺宝双倍',
      note: '21:00 开始',
      icon: legacyActivityAssets.activitySparkIcon,
      actionKey: 'dungeon_run'
    },
    cultivationSummary: {
      action: '查看修行',
      actionKey: 'cultivation',
      items: [
        { label: '修行地图', value: '青木林地', subtext: '进行中' },
        { label: '结束时间', value: '13:40', subtext: '剩余 18 分钟' },
        { label: '铜钱收益', value: '240' },
        { label: '幻兽经验', value: '160' }
      ]
    },
    entries: [
      { label: '世界地图', actionKey: 'world_map' },
      { label: '联盟', actionKey: 'alliance' },
      { label: '幻兽', actionKey: 'pet_catalog' },
      { label: '背包', actionKey: 'assets' },
      { label: '竞技场', actionKey: 'arena' },
      { label: '庄园', actionKey: 'growth_manor' },
      { label: '修行', actionKey: 'cultivation' },
      { label: '排行', actionKey: 'ranking' }
    ],
    messages: [
      '世界消息：青木林地今日双倍经验已开启',
      '联盟消息：今晚 20:00 盟战锁定名单',
      '系统消息：VIP 每日宝箱可领取'
    ]
  }
}
