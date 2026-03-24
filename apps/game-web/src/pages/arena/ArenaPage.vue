<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import { legacyHeadAssets, legacyItemAssets } from '@/assets/legacy'

const opponents = [
  {
    id: 1,
    name: '烈焰狼王',
    rank: 'S 赛区 · 当前第一',
    status: '战力 6,820',
    icon: legacyHeadAssets.head1Icon
  },
  {
    id: 2,
    name: '寒枝鹿灵',
    rank: 'S 赛区 · 当前第二',
    status: '战力 5,980',
    icon: legacyHeadAssets.head2Icon
  }
]

const reward = {
  title: '赛季通行证奖励',
  description: '累计胜利 5 场可领取',
  icon: legacyItemAssets.expPillIcon,
  note: '火能经验增益 +12%'
}
</script>

<template>
  <section class="arena-page">
    <UiPageHero
      eyebrow="格斗擂台"
      title="神域竞技场"
      description="每日刷新对手与赛季奖励清晰展示，明确当前战斗目标。"
      tone="violet"
      meta-label="推荐时段"
      meta-value="18:00 - 23:00"
    />

    <UiPanelCard title="今日对手">
      <div class="opponent-list">
        <article
          v-for="opponent in opponents"
          :key="opponent.id"
          class="opponent-card"
          :data-testid="`arena-opponent-${opponent.id}`"
        >
          <img v-if="opponent.icon" :src="opponent.icon" :alt="`${opponent.name} 头像`" />
          <div>
            <strong>{{ opponent.name }}</strong>
            <span>{{ opponent.rank }}</span>
            <p>{{ opponent.status }}</p>
          </div>
        </article>
      </div>
    </UiPanelCard>

    <section class="grid">
      <UiPanelCard title="挑战提示">
        <p>挑战提醒：刷新的对手中 1 号位自动配备最新战装，建议先击败他获取排行榜增益。</p>
        <button type="button" class="ghost-button">立即挑战</button>
      </UiPanelCard>

      <UiPanelCard title="奖励目标" wide>
        <div class="reward-card" data-testid="arena-reward-card">
          <img v-if="reward.icon" :src="reward.icon" :alt="`${reward.title} 图标`" />
          <div>
            <strong>{{ reward.title }}</strong>
            <p>{{ reward.description }}</p>
            <small>{{ reward.note }}</small>
          </div>
        </div>
        <p class="hint">挑战每日对手并达成连胜，即可解锁火能经验倍增。</p>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.arena-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #111827;
}

.opponent-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.opponent-card {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px;
  border-radius: 16px;
  background: #f5f3ff;
}

.opponent-card img {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #ddd6fe;
}

.opponent-card span,
.opponent-card p {
  margin: 0;
  color: #4b5563;
  font-size: 13px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.reward-card {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 14px;
  background: #fff7ed;
  border: 1px solid #fde3b6;
  align-items: center;
}

.reward-card img {
  width: 46px;
  height: 46px;
  object-fit: cover;
  border-radius: 12px;
}

.hint {
  margin-top: 8px;
  color: #92400e;
}

.ghost-button {
  margin-top: 12px;
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid #ddd6fe;
  background: #f5f3ff;
  color: #4338ca;
  cursor: pointer;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
