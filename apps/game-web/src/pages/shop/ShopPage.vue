<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import { legacyItemAssets } from '@/assets/legacy'

const goods = [
  {
    name: '首充礼包',
    price: '¥6',
    desc: '主线推进、火能修行与洗炼资源的起步礼包',
    icon: legacyItemAssets.expPillIcon
  },
  {
    name: '焚火补给包',
    price: '¥30',
    desc: '适合联盟玩法日常补充，覆盖原石与联盟资金',
    icon: legacyItemAssets.fireEvolutionStoneIcon
  },
  { name: '养成冲刺包', price: '¥68', desc: '集中补充修行、庄园与副本效率道具' }
]

const paymentSteps = ['创建订单', '拉起支付', '支付回调', '幂等发货']
</script>

<template>
  <section class="shop-page">
    <UiPageHero
      eyebrow="商城中心"
      title="推荐商品"
      description="突出首购、效率包与常驻礼包，保持下单到发货的路径透明。"
      tone="cyan"
      meta-value="今日热销 3 款"
    />

    <section class="goods-grid">
      <UiPanelCard v-for="item in goods" :key="item.name" :data-testid="`shop-item-${item.name}`">
        <div class="goods-card">
          <div class="goods-card__icon" v-if="item.icon">
            <img :src="item.icon" :alt="`${item.name} 图标`" />
          </div>
          <div class="goods-card__meta">
            <h2>{{ item.name }}</h2>
            <p>{{ item.desc }}</p>
          </div>
          <strong>{{ item.price }}</strong>
        </div>
      </UiPanelCard>
    </section>

    <UiPanelCard title="支付流程">
      <UiChipGroup :items="paymentSteps" tone="blue" />
      <p class="tip">后端将以订单号 + 幂等键保证重复回调不重复发货，页面仅展示用户可感知流程。</p>
    </UiPanelCard>
  </section>
</template>

<style scoped>
.shop-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.goods-card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 12px;
  align-items: center;
  background: linear-gradient(180deg, #fff, #f8fafc);
}

.goods-card h2 {
  margin: 0 0 10px;
}

.goods-card__icon img {
  width: 48px;
  height: 48px;
  object-fit: contain;
  background: #fff;
  border-radius: 12px;
  padding: 6px;
  border: 1px solid #ede9fe;
}

.goods-card__meta {
  display: flex;
  flex-direction: column;
  gap: 6px;
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

.tip {
  margin-top: 12px;
}

@media (max-width: 768px) {
  .goods-card {
    flex-direction: column;
  }
}
</style>
