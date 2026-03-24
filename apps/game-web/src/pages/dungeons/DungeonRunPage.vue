<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { legacyItemAssets } from '@/assets/legacy'

const rewardDetails = [
  {
    label: '火系进化石',
    icon: legacyItemAssets.fireEvolutionStoneIcon,
    note: '进化主战幻兽必备'
  },
  {
    label: '幻兽经验丹',
    icon: legacyItemAssets.expPillIcon,
    note: '副本快速补位'
  }
]

const rewards = [...rewardDetails.map((detail) => detail.label), '铜钱奖励', 'Boss 首通宝箱']

const expectations = [
  { label: '副本名称', value: '青木林地' },
  { label: '所属城市', value: '青木城' },
  { label: '推荐等级段', value: 'Lv.30-35' },
  { label: '经验门槛', value: '主战幻兽经验 ≥ 4,500' },
  { label: '地图主题', value: '林地 / 木属性怪物为主' },
  { label: '剩余次数', value: '今日 2 / 3 次' }
]
</script>

<template>
  <section class="dungeon-run-page">
    <UiPageHero
      eyebrow="进入前预期管理"
      title="地图副本"
      description="先说明副本收益、门槛、成本与主题，帮助玩家决定现在是否值得进入。"
      tone="orange"
      meta-value="进入成本：活力 12"
    />

    <section class="content-grid">
      <UiPanelCard title="副本信息">
        <UiStatGrid :items="expectations" tone="warm" min-width="160px" />
      </UiPanelCard>

      <UiPanelCard title="地图简介">
        <p>青木林地是玩家 30 级后最稳定的养成地图，怪物以木属性为主，适合刷火系培养材料。</p>
      </UiPanelCard>

      <UiPanelCard title="Boss 主要奖励" wide>
        <UiChipGroup :items="rewards" tone="orange" min-width="160px" />
        <div class="reward-grid">
          <article
            v-for="detail in rewardDetails"
            :key="detail.label"
            class="reward-card"
            :data-testid="`dungeon-reward-${detail.label}`"
          >
            <img v-if="detail.icon" :src="detail.icon" :alt="`${detail.label} 图标`" />
            <div>
              <strong>{{ detail.label }}</strong>
              <p>{{ detail.note }}</p>
            </div>
          </article>
        </div>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.dungeon-run-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.reward-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-top: 16px;
}

.reward-card {
  display: flex;
  gap: 10px;
  align-items: center;
  padding: 12px;
  border-radius: 14px;
  background: #fef9f2;
  border: 1px solid #fde9c9;
}

.reward-card img {
  width: 42px;
  height: 42px;
  object-fit: contain;
  padding: 4px;
  border-radius: 12px;
  background: #fff;
}

@media (max-width: 768px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
