<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import { legacyHeadAssets } from '@/assets/legacy'

const rankCategories = ['火能榜', '联盟榜', '帮派榜', '跨服榜']
const topPlayers = [
  {
    name: '烈焰狼王',
    title: '永恒战火 · 首席',
    level: 'Lv.48',
    power: '战力 6,820',
    icon: legacyHeadAssets.head1Icon,
    tag: '赛季 TOP1'
  },
  {
    name: '寒枝鹿灵',
    title: '资源调度 · 副盟主',
    level: 'Lv.44',
    power: '战力 5,980',
    icon: legacyHeadAssets.head2Icon,
    tag: '每日助攻'
  },
  {
    name: '青木灵猿',
    title: '侦察前锋',
    level: 'Lv.42',
    power: '战力 5,430',
    tag: '草堂勇士'
  }
]
</script>

<template>
  <section class="ranking-page">
    <UiPageHero
      eyebrow="神域荣耀"
      title="当前排行榜"
      description="实时展示联盟 / 帮派 / 个人战力，榜单锁定火能分段，方便玩家找标杆和炫耀点。"
      tone="violet"
      meta-label="活跃榜"
      meta-value="每日 12:00 刷新"
    />

    <UiPanelCard title="排行榜分类">
      <UiChipGroup :items="rankCategories" tone="blue" min-width="120px" />
    </UiPanelCard>

    <UiPanelCard title="前 3 名亮屏">
      <div class="rank-list">
        <article
          v-for="(player, index) in topPlayers"
          :key="player.name"
          class="rank-card"
          :data-testid="`ranking-player-${index + 1}`"
        >
          <div class="rank-card__avatar">
            <img
              v-if="player.icon"
              :src="player.icon"
              :alt="`${player.name} 头像`"
            />
            <div v-else class="rank-card__avatar-placeholder">{{ index + 1 }}</div>
          </div>
          <div class="rank-card__meta">
            <strong>{{ player.name }}</strong>
            <span class="rank-card__title">{{ player.title }}</span>
            <span class="rank-card__level">{{ player.level }}</span>
            <p>{{ player.power }}</p>
          </div>
          <small class="rank-card__tag">{{ player.tag }}</small>
        </article>
      </div>
    </UiPanelCard>
  </section>
</template>

<style scoped>
.ranking-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #0f172a;
}

.rank-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.rank-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px;
  border-radius: 16px;
  background: linear-gradient(180deg, #fff, #f4f1ff);
  box-shadow: 0 10px 20px rgba(15, 23, 42, 0.08);
}

.rank-card__avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 2px solid #e0d7ff;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fdf6ff;
}

.rank-card__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.rank-card__avatar-placeholder {
  font-size: 20px;
  font-weight: 700;
  color: #7c3aed;
}

.rank-card__meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.rank-card__title {
  color: #6b7280;
  font-size: 13px;
}

.rank-card__level {
  font-size: 12px;
  color: #4b5563;
}

.rank-card__tag {
  color: #4338ca;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  background: #ede9fe;
}

@media (max-width: 768px) {
  .rank-card {
    flex-direction: column;
    align-items: flex-start;
  }

  .rank-card__meta {
    width: 100%;
  }
}
</style>
