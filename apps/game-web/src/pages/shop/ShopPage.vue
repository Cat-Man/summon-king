<script setup lang="ts">
const goods = [
  { name: '首充礼包', price: '¥6', desc: '主线推进、火能修行与洗炼资源的起步礼包' },
  { name: '焚火补给包', price: '¥30', desc: '适合联盟玩法日常补充，覆盖原石与联盟资金' },
  { name: '养成冲刺包', price: '¥68', desc: '集中补充修行、庄园与副本效率道具' }
]

const paymentSteps = ['创建订单', '拉起支付', '支付回调', '幂等发货']
</script>

<template>
  <section class="shop-page">
    <header class="hero-card">
      <div>
        <p class="eyebrow">商城中心</p>
        <h1>推荐商品</h1>
        <p>突出首购、效率包与常驻礼包，保持下单到发货的路径透明。</p>
      </div>
      <div class="hero-side">今日热销 3 款</div>
    </header>

    <section class="goods-grid">
      <article v-for="item in goods" :key="item.name" class="goods-card">
        <div>
          <h2>{{ item.name }}</h2>
          <p>{{ item.desc }}</p>
        </div>
        <strong>{{ item.price }}</strong>
      </article>
    </section>

    <section class="panel-card">
      <h2>支付流程</h2>
      <div class="step-list">
        <span v-for="step in paymentSteps" :key="step" class="step-item">{{ step }}</span>
      </div>
      <p class="tip">后端将以订单号 + 幂等键保证重复回调不重复发货，页面仅展示用户可感知流程。</p>
    </section>
  </section>
</template>

<style scoped>
.shop-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.hero-card,
.goods-card,
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
  background: linear-gradient(135deg, #0891b2, #2563eb);
  color: #fff;
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.12em;
  opacity: 0.8;
}

.hero-card h1,
.goods-card h2,
.panel-card h2 {
  margin: 0 0 10px;
}

.hero-side {
  align-self: flex-start;
  padding: 10px 14px;
  border-radius: 999px;
  background: rgb(255 255 255 / 18%);
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.goods-card {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  background: linear-gradient(180deg, #fff, #f8fafc);
}

.goods-card p,
.tip {
  margin: 0;
  color: #6b7280;
  line-height: 1.6;
}

.goods-card strong {
  align-self: flex-start;
  font-size: 24px;
  color: #2563eb;
}

.step-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.step-item {
  padding: 12px 14px;
  border-radius: 14px;
  background: #eff6ff;
  color: #1d4ed8;
  text-align: center;
}

@media (max-width: 768px) {
  .hero-card,
  .goods-card {
    flex-direction: column;
  }
}
</style>
