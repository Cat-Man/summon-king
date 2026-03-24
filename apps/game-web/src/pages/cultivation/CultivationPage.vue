<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { legacyItemAssets, legacyPetAssets } from '@/assets/legacy'

const currentStatus = [
  { label: '可修行地图', value: '定老城 · 青木林地' },
  { label: '开始时间', value: '09:40' },
  { label: '结束时间', value: '13:40' },
  { label: '当前状态', value: '进行中' }
]

const durationPlans = ['2 小时自由开放', '4 小时自由开放', '8 小时自由开放', '12 小时 VIP2 开放', '24 小时 VIP5 开放']
const rewardPreview = [
  { label: '玩家声望 x 160' },
  { label: '幻兽经验 x 2,400', icon: legacyItemAssets.expPillIcon },
  { label: '强化石 x 18' },
  { label: '召唤球概率掉落' }
]
const teamSnapshot = [
  {
    name: '烈焰狼王 · Lv.36',
    icon: legacyPetAssets.flameWolfKingIcon,
    dataId: 'cultivation-team-烈焰狼王'
  },
  {
    name: '寒枝鹿灵 · Lv.34',
    icon: legacyPetAssets.coldBranchDeerIcon,
    dataId: 'cultivation-team-寒枝鹿灵'
  },
  { name: '玄甲龟 · Lv.35' },
  { name: '流云狐 · Lv.33' },
  { name: '焚羽雀 · Lv.32' }
]
const rules = [
  '修行至少 5 分钟才有收益',
  '30 秒内舍去，30-60 秒进 1 分钟',
  '2 小时及以上概率掉对应地图召唤球',
  '修行队伍 = 战斗队伍',
  '到期领取，不按超时时长追加收益'
]
</script>

<template>
  <section class="cultivation-page">
    <UiPageHero
      eyebrow="离线收益规划"
      title="修行总览"
      description="集中展示当前修行状态、可选时长、收益规则与队伍快照，方便玩家判断现在是否继续挂机。"
      tone="green"
      meta-value="还剩 2 小时 18 分钟"
    />

    <section class="grid">
      <UiPanelCard title="当前修行状态">
        <template #actions>
          <button type="button" class="claim-button">领取收益</button>
        </template>
        <UiStatGrid :items="currentStatus" tone="mint" min-width="150px" />
      </UiPanelCard>

      <UiPanelCard title="修行时长">
        <UiChipGroup :items="durationPlans" tone="green" min-width="150px" />
      </UiPanelCard>

      <UiPanelCard title="收益规则">
        <ul class="text-list">
          <li v-for="item in rules" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="队伍快照">
        <ul class="text-list">
          <li v-for="item in teamSnapshot" :key="item.name" :data-testid="item.dataId || null">
            <img
              v-if="item.icon"
              :src="item.icon"
              :alt="`${item.name} 头像`"
              class="list-icon"
            />
            {{ item.name }}
          </li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="奖励预览" wide>
        <ul class="reward-list">
          <li
            v-for="reward in rewardPreview"
            :key="reward.label"
            :data-testid="reward.label === '幻兽经验 x 2,400' ? 'cultivation-reward-幻兽经验 x 2,400' : null"
          >
            <img
              v-if="reward.icon"
              :src="reward.icon"
              :alt="`${reward.label} 图标`"
              class="list-icon"
            />
            {{ reward.label }}
          </li>
        </ul>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.cultivation-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.claim-button {
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #dcfce7;
  color: #166534;
  cursor: pointer;
}

.text-list {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
  line-height: 1.6;
}

.text-list img,
.reward-list img {
  width: 34px;
  height: 34px;
  object-fit: cover;
  border-radius: 10px;
  margin-right: 8px;
  background: #fff;
  border: 1px solid #d1fae5;
}

.reward-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
}

.reward-list li {
  display: flex;
  align-items: center;
  gap: 8px;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
