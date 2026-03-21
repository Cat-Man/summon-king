<script setup lang="ts">
const timeline = [
  { phase: '报名阶段', time: '周四 00:00 - 周六 20:00', status: '进行中' },
  { phase: '成员签到', time: '锁定前完成主力签到', status: '待完成' },
  { phase: '对阵信息', time: '20:05 自动匹配并展示占领点', status: '待开启' },
  { phase: '结果结算', time: '战斗结束后发放联盟战功', status: '待开启' }
]

const redeemOptions = [
  { name: '焚火晶礼包', cost: 120 },
  { name: '盟战增援令', cost: 80 },
  { name: '高阶洗炼石', cost: 160 }
]
</script>

<template>
  <section class="war-page">
    <header class="hero-card">
      <div>
        <p class="eyebrow">联盟赛事</p>
        <h1>盟战作战台</h1>
        <p>集中展示阶段进度、签到状态、匹配结果与战功兑换，方便战前指挥与战后结算查看。</p>
      </div>
      <div class="score-card">
        <span>本周联盟战功</span>
        <strong>368</strong>
      </div>
    </header>

    <section class="timeline-card">
      <h2>本轮流程</h2>
      <div class="timeline-list">
        <article v-for="item in timeline" :key="item.phase" class="timeline-item">
          <strong>{{ item.phase }}</strong>
          <span>{{ item.time }}</span>
          <em>{{ item.status }}</em>
        </article>
      </div>
    </section>

    <section class="grid">
      <article class="panel-card">
        <h2>成员签到</h2>
        <p>已签到 14 / 20，推荐优先确认主力队与替补队已全部完成签到。</p>
      </article>

      <article class="panel-card">
        <h2>对阵信息</h2>
        <p>预计对手：赤霄盟，战区为 3 号焚火谷，占领点数 5 个。</p>
      </article>

      <article class="panel-card panel-card--wide">
        <h2>战功兑换</h2>
        <div class="redeem-list">
          <div v-for="item in redeemOptions" :key="item.name" class="redeem-item">
            <div>
              <strong>{{ item.name }}</strong>
              <p>消耗 {{ item.cost }} 点联盟战功，盟主 / 副盟主可发起兑换。</p>
            </div>
            <button type="button">查看兑换</button>
          </div>
        </div>
      </article>
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

.hero-card,
.timeline-card,
.panel-card {
  padding: 20px;
  border-radius: 20px;
  background: #fff;
  box-shadow: 0 10px 30px rgb(15 23 42 / 8%);
}

.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  background: linear-gradient(135deg, #7c3aed, #db2777);
  color: #fff;
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.12em;
  opacity: 0.85;
}

.hero-card h1,
.timeline-card h2,
.panel-card h2 {
  margin: 0 0 10px;
}

.score-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  align-self: flex-start;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgb(255 255 255 / 18%);
}

.score-card strong {
  font-size: 28px;
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

.panel-card--wide {
  grid-column: 1 / -1;
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
  .hero-card,
  .redeem-item {
    flex-direction: column;
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
