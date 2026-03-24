<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import { legacyAllianceAvatarMap, legacyItemAssets } from '@/assets/legacy'

const timeline = [
  { phase: '报名阶段', time: '周四 00:00 - 周六 20:00', status: '进行中' },
  { phase: '成员签到', time: '锁定前完成主力签到', status: '待完成' },
  { phase: '对阵信息', time: '20:05 自动匹配并展示占领点', status: '待开启' },
  { phase: '结果结算', time: '战斗结束后发放联盟战功', status: '待开启' }
]

const redeemOptions = [
  { name: '焚火晶礼包', cost: 120, icon: legacyItemAssets.fireEvolutionStoneIcon },
  { name: '盟战增援令', cost: 80 },
  { name: '高阶洗炼石', cost: 160 }
]

const commanders = [
  {
    role: '盟主',
    status: '指挥中 · 在线',
    icon: legacyAllianceAvatarMap.盟主,
    dataId: 'war-commander-盟主'
  },
  {
    role: '副盟主',
    status: '战术支援 · 待命',
    icon: legacyAllianceAvatarMap.副盟主,
    dataId: 'war-commander-副盟主'
  }
]
</script>

<template>
  <section class="war-page">
    <UiPageHero
      eyebrow="联盟赛事"
      title="盟战作战台"
      description="集中展示阶段进度、签到状态、匹配结果与战功兑换，方便战前指挥与战后结算查看。"
      tone="violet"
      meta-label="本周联盟战功"
      meta-value="368"
    />

    <UiPanelCard title="本轮流程">
      <div class="timeline-list">
        <article v-for="item in timeline" :key="item.phase" class="timeline-item">
          <strong>{{ item.phase }}</strong>
          <span>{{ item.time }}</span>
          <em>{{ item.status }}</em>
        </article>
      </div>
    </UiPanelCard>

<section class="grid">
  <UiPanelCard title="成员签到">
    <p>已签到 14 / 20，推荐优先确认主力队与替补队已全部完成签到。</p>
  </UiPanelCard>

  <UiPanelCard title="指挥席">
    <div class="commander-list">
      <article v-for="commander in commanders" :key="commander.role" class="commander-card" :data-testid="commander.dataId">
        <img v-if="commander.icon" :src="commander.icon" :alt="`${commander.role} 头像`" />
        <div>
          <strong>{{ commander.role }}</strong>
          <span>{{ commander.status }}</span>
        </div>
      </article>
    </div>
  </UiPanelCard>

  <UiPanelCard title="战功兑换" wide>
    <div class="redeem-list">
      <div v-for="item in redeemOptions" :key="item.name" class="redeem-item" :data-testid="`war-redeem-${item.name}`">
        <img v-if="item.icon" :src="item.icon" :alt="`${item.name} 图标`" />
        <div>
          <strong>{{ item.name }}</strong>
          <p>消耗 {{ item.cost }} 点联盟战功，盟主 / 副盟主可发起兑换。</p>
        </div>
        <button type="button">查看兑换</button>
      </div>
    </div>
  </UiPanelCard>
</section>
  </section>
</template>

<style scoped>
.war-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.timeline-list,
.redeem-list,
.grid {
  display: grid;
  gap: 14px;
}

.timeline-list {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.timeline-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.timeline-item span {
  color: #6b7280;
  font-size: 14px;
}

.timeline-item em {
  font-style: normal;
  color: #7c3aed;
  font-size: 13px;
}

.grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.redeem-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.redeem-item p {
  margin: 6px 0 0;
  color: #6b7280;
}

.redeem-item button {
  align-self: center;
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #ede9fe;
  color: #6d28d9;
  cursor: pointer;
}

@media (max-width: 768px) {
  .redeem-item {
    flex-direction: column;
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
