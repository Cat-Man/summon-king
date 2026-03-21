<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'

const currentPrivileges = [
  '每日宝箱可领 1 次',
  '修行队列 +1',
  '庄园收获一键完成',
  '免费洗炼次数 +2'
]

const nextPrivileges = [
  '战骨背包容量 +20',
  '副本重置次数 +1',
  '火能修行房间上限 +1'
]
</script>

<template>
  <section class="vip-page">
    <UiPageHero
      eyebrow="升级价值感"
      title="当前 VIP 4"
      description="今日权益与下一档成长收益同屏展示，突出效率型特权。"
      tone="amber"
      meta-label="距离下一档"
      meta-value="还差 320 充值点"
    />

    <section class="status-grid">
      <UiPanelCard title="每日宝箱">
        <p>今日状态：未领取</p>
        <button type="button">立即领取</button>
      </UiPanelCard>

      <UiPanelCard title="见面礼包">
        <p>VIP 4 礼包已解锁，可查看铜钱、焚火晶与洗炼石内容。</p>
      </UiPanelCard>
    </section>

    <section class="benefit-grid">
      <UiPanelCard title="当前权益">
        <ul>
          <li v-for="item in currentPrivileges" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="下一档新增权益">
        <ul>
          <li v-for="item in nextPrivileges" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.vip-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.status-grid,
.benefit-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

ul {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #374151;
}

button {
  margin-top: 12px;
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #fef3c7;
  color: #92400e;
  cursor: pointer;
}

@media (max-width: 768px) {
  .status-grid,
  .benefit-grid {
    grid-template-columns: 1fr;
  }
}
</style>
