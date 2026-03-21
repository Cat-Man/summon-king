<script setup lang="ts">
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
    <header class="hero-card">
      <div>
        <p class="eyebrow">升级价值感</p>
        <h1>当前 VIP 4</h1>
        <p>今日权益与下一档成长收益同屏展示，突出效率型特权。</p>
      </div>
      <div class="hero-side">
        <span>距离下一档</span>
        <strong>还差 320 充值点</strong>
      </div>
    </header>

    <section class="status-grid">
      <article class="panel-card">
        <h2>每日宝箱</h2>
        <p>今日状态：未领取</p>
        <button type="button">立即领取</button>
      </article>

      <article class="panel-card">
        <h2>见面礼包</h2>
        <p>VIP 4 礼包已解锁，可查看铜钱、焚火晶与洗炼石内容。</p>
      </article>
    </section>

    <section class="benefit-grid">
      <article class="panel-card">
        <h2>当前权益</h2>
        <ul>
          <li v-for="item in currentPrivileges" :key="item">{{ item }}</li>
        </ul>
      </article>

      <article class="panel-card">
        <h2>下一档新增权益</h2>
        <ul>
          <li v-for="item in nextPrivileges" :key="item">{{ item }}</li>
        </ul>
      </article>
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

.hero-card,
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
  background: linear-gradient(135deg, #f59e0b, #f97316);
  color: #fff;
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.12em;
  opacity: 0.8;
}

.hero-card h1,
.panel-card h2 {
  margin: 0 0 10px;
}

.hero-side {
  display: flex;
  flex-direction: column;
  align-self: flex-start;
  gap: 6px;
  padding: 12px 14px;
  border-radius: 16px;
  background: rgb(255 255 255 / 18%);
}

.status-grid,
.benefit-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.panel-card ul {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #374151;
}

.panel-card button {
  margin-top: 12px;
  padding: 10px 14px;
  border: 0;
  border-radius: 12px;
  background: #fef3c7;
  color: #92400e;
  cursor: pointer;
}

@media (max-width: 768px) {
  .hero-card {
    flex-direction: column;
  }

  .status-grid,
  .benefit-grid {
    grid-template-columns: 1fr;
  }
}
</style>
