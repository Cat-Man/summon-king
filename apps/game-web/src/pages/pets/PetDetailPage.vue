<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'
import { legacyInventoryIconMap, legacyPetIconMap } from '@/assets/legacy'

const overview = [
  { label: '等级', value: 'Lv.36' },
  { label: '种族', value: '火系 · 狼族' },
  { label: '成长率', value: '1588' },
  { label: '境界', value: '天界' },
  { label: '攻资', value: '3420' },
  { label: '速度资质', value: '1280' }
]

const skills = [
  '烈焰撕咬：单体火系爆发，适合作为主战收割技',
  '狼王威吓：开场降低敌方前排防御',
  '焚风追击：当目标残血时追加一次追击'
]

const bones = ['头骨 Lv.8', '胸骨 Lv.7', '臂骨 Lv.7', '手骨 Lv.6', '腿骨 Lv.6', '尾骨 Lv.6', '元魂 Lv.5']
const spirits = ['火灵·暴怒', '木灵·护心', '神灵·迅捷']
const souls = ['龙魂·灼炎', '天魂·破军', '玄魂·护体']
const sources = ['初始赠送可得同系基础幻兽', '青木林地相关掉落', '修行 2 小时以上概率获得召唤球', '礼包与商城活动补充']
const heroIcon = legacyPetIconMap['烈焰狼王']
const evolutionStoneIcon = legacyInventoryIconMap['火系进化石']
</script>

<template>
  <section class="pet-detail-page">
    <UiPageHero
      eyebrow="养成中枢"
      title="烈焰狼王"
      description="把技能、战骨、战灵、魔魂、进化和重生入口集中展示，作为当前主养成幻兽的总控制台。"
      tone="orange"
      meta-label="综合战力"
      meta-value="6,820"
    />

    <UiStatGrid :items="overview" tone="warm" min-width="130px" />

    <div v-if="heroIcon" class="pet-detail-hero" data-testid="pet-detail-hero-art">
      <img :src="heroIcon" alt="烈焰狼王 头像" />
      <strong>烈焰狼王</strong>
    </div>

    <section class="grid">
      <UiPanelCard title="技能">
        <ul class="text-list">
          <li v-for="item in skills" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="战骨">
        <UiChipGroup :items="bones" tone="orange" min-width="120px" />
      </UiPanelCard>

      <UiPanelCard title="战灵">
        <UiChipGroup :items="spirits" tone="blue" min-width="140px" />
      </UiPanelCard>

      <UiPanelCard title="魔魂">
        <UiChipGroup :items="souls" tone="green" min-width="140px" />
      </UiPanelCard>

      <UiPanelCard title="进化/升境">
        <div class="evolution-cost" data-testid="pet-detail-evolution-cost">
          <img v-if="evolutionStoneIcon" :src="evolutionStoneIcon" alt="火系进化石" />
          <div>
            <strong>火系进化石 6 / 8</strong>
            <span>铜钱 180,000 / 200,000</span>
          </div>
        </div>
        <p>当前升境条件：火系进化石 6 / 8、铜钱 180,000 / 200,000。满足后可直接提交升境并刷新最终属性。</p>
      </UiPanelCard>

      <UiPanelCard title="重生/放生">
        <p>重生将返还大部分培养材料；放生将永久失去该幻兽实例，提交前必须二次确认。</p>
      </UiPanelCard>

      <UiPanelCard title="来源说明" wide>
        <ul class="text-list">
          <li v-for="item in sources" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-detail-page {
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

.text-list {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
  line-height: 1.6;
}

.grid p {
  margin: 0;
  color: #4b5563;
  line-height: 1.7;
}

.pet-detail-hero {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px;
  border-radius: 20px;
  background: #fff7ed;
  color: #b45309;
}

.pet-detail-hero img {
  width: 80px;
  height: 80px;
  border-radius: 999px;
  object-fit: cover;
  border: 2px solid #fb923c;
}

.evolution-cost {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 10px;
}

.evolution-cost img {
  width: 40px;
  height: 40px;
  object-fit: contain;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
