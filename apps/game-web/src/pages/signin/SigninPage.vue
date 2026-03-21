<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'

const signinDays = [
  { day: 'Day 1', reward: '铜钱 x 20k', claimed: true },
  { day: 'Day 2', reward: '洗炼石 x 5', claimed: true },
  { day: 'Day 3', reward: '焚火晶 x 60', claimed: false },
  { day: 'Day 4', reward: '幻兽经验丹 x 8', claimed: false }
]

const gifts = [
  '新服见面礼：已解锁，待领取',
  '七日成长礼：进行中，需完成主线',
  '联盟联动礼：加入联盟后开放'
]
</script>

<template>
  <section class="signin-page">
    <UiPageHero
      eyebrow="每日必做"
      title="今日签到"
      description="把签到、连续签到、礼包中心与兑换码入口聚合在一个页面，便于首页工作台直接跳转。"
      tone="green"
      meta-value="连续签到 9 天"
    />

    <section class="grid">
      <UiPanelCard title="连续签到">
        <div class="signin-list">
          <div v-for="item in signinDays" :key="item.day" class="signin-item">
            <strong>{{ item.day }}</strong>
            <span>{{ item.reward }}</span>
            <em>{{ item.claimed ? '已领取' : '待领取' }}</em>
          </div>
        </div>
      </UiPanelCard>

      <UiPanelCard title="礼包中心">
        <ul>
          <li v-for="gift in gifts" :key="gift">{{ gift }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="兑换码" wide>
        <div class="redeem-box">
          <span>输入礼包码后校验可用状态，并由后端完成幂等兑换。</span>
          <button type="button">前往兑换</button>
        </div>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.signin-page {
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

.signin-list {
  display: grid;
  gap: 10px;
}

.signin-item {
  display: grid;
  gap: 4px;
  padding: 12px 14px;
  border-radius: 14px;
  background: #f0fdf4;
}

.signin-item span,
ul {
  color: #4b5563;
}

.signin-item em {
  font-style: normal;
  color: #15803d;
  font-size: 13px;
}

ul {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.redeem-box {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  padding: 14px;
  border-radius: 14px;
  background: #f8fafc;
}

.redeem-box button {
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #dcfce7;
  color: #166534;
  cursor: pointer;
}

@media (max-width: 768px) {
  .redeem-box {
    flex-direction: column;
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
