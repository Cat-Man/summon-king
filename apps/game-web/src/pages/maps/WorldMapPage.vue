<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'

const currentCityDungeons = [
  { name: '青木林地', level: 'Lv.30-35', status: 'Boss 可挑战' },
  { name: '灰岩洞窟', level: 'Lv.34-38', status: '今日剩余 2 次' }
]

const nearbyCities = [
  { name: '流火城', level: 'Lv.36-42', action: '可传送' },
  { name: '白露港', level: 'Lv.40-46', action: '需先解锁主线' },
  { name: '霜风峡', level: 'Lv.45-52', action: '未开放' }
]
</script>

<template>
  <section class="world-map-page">
    <UiPageHero
      eyebrow="成长路径页"
      title="世界地图"
      description="让玩家一眼看懂我现在在哪、下一步去哪里、当前是否具备挑战条件。"
      tone="teal"
      meta-label="当前位置"
      meta-value="青木城"
    />

    <section class="grid">
      <UiPanelCard title="当前城市副本">
        <div class="list">
          <div v-for="item in currentCityDungeons" :key="item.name" class="list-item">
            <strong>{{ item.name }}</strong>
            <span>{{ item.level }}</span>
            <em>{{ item.status }}</em>
          </div>
        </div>
      </UiPanelCard>

      <UiPanelCard title="当前推荐目标">
        <p>推荐优先挑战「青木林地」Boss，拿到火系幻兽经验与木属性进化材料。</p>
      </UiPanelCard>

      <UiPanelCard title="临近城市" wide>
        <div class="city-grid">
          <div v-for="city in nearbyCities" :key="city.name" class="city-item">
            <strong>{{ city.name }}</strong>
            <span>等级段提示：{{ city.level }}</span>
            <em>移动/传送：{{ city.action }}</em>
          </div>
        </div>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.world-map-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.grid,
.city-grid {
  display: grid;
  gap: 16px;
}

.grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-item,
.city-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px;
  border-radius: 14px;
  background: #f0fdfa;
}

.list-item span,
.city-item span {
  color: #4b5563;
}

.list-item em,
.city-item em {
  font-style: normal;
  color: #0f766e;
  font-size: 13px;
}

.city-grid {
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
