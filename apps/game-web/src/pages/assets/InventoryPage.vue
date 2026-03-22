<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

import { legacyInventoryIconMap } from '@/assets/legacy'

const wallet = [
  { label: '铜钱', value: '286,400' },
  { label: '元宝', value: '1,280' },
  { label: '声望', value: '540' },
  { label: '活力', value: '86 / 120' }
]

const walletResources = [
  { label: '铜钱', value: '286,400', icon: legacyInventoryIconMap['铜钱'] },
  { label: '元宝', value: '1,280', icon: legacyInventoryIconMap['元宝'] }
]

const items = [
  { name: '中级经验丹', count: 'x12', action: '可直接喂给主力幻兽', icon: legacyInventoryIconMap['中级经验丹'] },
  { name: '火系进化石', count: 'x6', action: '烈焰狼王升境材料', icon: legacyInventoryIconMap['火系进化石'] },
  { name: '强化石', count: 'x48', action: '用于战骨与装备强化' },
  { name: '召唤球碎片', count: 'x18', action: '可合成低阶召唤球' }
]

const rules = [
  '使用道具前先校验目标幻兽与当前玩法是否匹配',
  '出售前需展示预计获得的铜钱并二次确认',
  '高价值道具默认锁定，避免误操作批量出售'
]
</script>

<template>
  <section class="inventory-page">
    <UiPageHero
      eyebrow="资源与道具管理"
      title="背包"
      description="把资源钱包、普通背包与资源流水入口放在同一页，方便玩家快速判断当前是否缺钱、缺材料、缺活力。"
      tone="amber"
      meta-label="背包容量"
      meta-value="18 / 30"
    />

    <section class="grid">
      <UiPanelCard title="资源钱包">
        <div class="wallet-icons">
          <article
            v-for="res in walletResources"
            :key="res.label"
            class="wallet-icon"
            :data-testid="`wallet-icon-${res.label}`"
          >
            <img :src="res.icon" :alt="`${res.label} 图标`" />
            <div>
              <strong>{{ res.label }}</strong>
              <span>{{ res.value }}</span>
            </div>
          </article>
        </div>
        <UiStatGrid :items="wallet" tone="warm" min-width="130px" />
      </UiPanelCard>

      <UiPanelCard title="普通背包">
        <div class="item-list">
          <article
            v-for="item in items"
            :key="item.name"
            class="item-card"
            :data-testid="`inventory-item-${item.name}`"
          >
            <div class="item-card__top">
              <img v-if="item.icon" :src="item.icon" :alt="`${item.name} 图标`" />
              <div>
                <strong>{{ item.name }}</strong>
                <span>{{ item.count }}</span>
              </div>
            </div>
            <p>{{ item.action }}</p>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="使用/出售说明">
        <ul class="text-list">
          <li v-for="rule in rules" :key="rule">{{ rule }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="资源流水入口">
        <p>后续这里直接挂接资产流水查询，可按铜钱、元宝、道具获得与消耗来源回溯最近变更。</p>
        <button type="button" class="ghost-button">查看最近流水</button>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.inventory-page {
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

.item-list,
.text-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-list {
  margin: 0;
}

.item-card {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

  .wallet-icons {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    gap: 12px;
    margin-bottom: 12px;
  }

  .wallet-icon {
    display: flex;
    gap: 10px;
    align-items: center;
    padding: 10px 12px;
    border-radius: 14px;
    background: #faf5ff;
  }

  .wallet-icon img {
    width: 32px;
    height: 32px;
    object-fit: cover;
    border-radius: 8px;
  }

  .item-card__top {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 6px;
  }

  .item-card__top img {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 12px;
    box-shadow: 0 6px 16px rgb(15 23 42 / 15%);
  }

.item-card p,
.text-list,
.grid p {
  margin: 0;
  color: #4b5563;
  line-height: 1.6;
}

.text-list {
  padding-left: 18px;
}

.ghost-button {
  margin-top: 12px;
  padding: 10px 14px;
  border: 1px solid #fcd34d;
  border-radius: 12px;
  background: #fffbeb;
  color: #b45309;
  cursor: pointer;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
