<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

const overview = [
  { label: '已上阵', value: '5 / 5' },
  { label: '队伍核心', value: '烈焰狼王' },
  { label: '平均等级', value: 'Lv.34' },
  { label: '可替补', value: '4 只' }
]

const team = [
  { slot: '1 号位', name: '烈焰狼王', role: '主战输出', level: 'Lv.36', power: '6,820' },
  { slot: '2 号位', name: '寒枝鹿灵', role: '控制辅助', level: 'Lv.34', power: '4,920' },
  { slot: '3 号位', name: '玄甲龟', role: '前排承伤', level: 'Lv.35', power: '4,610' },
  { slot: '4 号位', name: '流云狐', role: '速度补位', level: 'Lv.33', power: '3,780' },
  { slot: '5 号位', name: '焚羽雀', role: '收割补刀', level: 'Lv.32', power: '3,350' }
]

const pets = [
  { name: '烈焰狼王', level: 'Lv.36', status: '已上阵', power: '6,820' },
  { name: '寒枝鹿灵', level: 'Lv.34', status: '已上阵', power: '4,920' },
  { name: '玄甲龟', level: 'Lv.35', status: '已上阵', power: '4,610' },
  { name: '流云狐', level: 'Lv.33', status: '已上阵', power: '3,780' },
  { name: '焚羽雀', level: 'Lv.32', status: '已上阵', power: '3,350' },
  { name: '霜牙虎', level: 'Lv.31', status: '可替补', power: '3,140' },
  { name: '青木灵猿', level: 'Lv.30', status: '可替补', power: '2,980' }
]

const strategies = ['保存阵容', '支持上下阵', '支持排序调位', '按综合战力筛选']
</script>

<template>
  <section class="pet-team-page">
    <UiPageHero
      eyebrow="战斗队与幻兽栏"
      title="主力阵容"
      description="先展示当前战斗队，再展示幻兽栏与主养成目标，后续可直接接保存阵容、上下阵和排序接口。"
      tone="blue"
      meta-label="综合战力"
      meta-value="23,480"
    />

    <UiStatGrid :items="overview" min-width="140px" />

    <section class="grid">
      <UiPanelCard title="当前战斗队">
        <div class="team-list">
          <article v-for="item in team" :key="item.slot" class="team-item">
            <div>
              <span class="subtle">{{ item.slot }}</span>
              <strong>{{ item.name }}</strong>
              <p>{{ item.role }} · {{ item.level }}</p>
            </div>
            <span class="power">战力 {{ item.power }}</span>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="主养成目标">
        <div class="focus-card">
          <strong>烈焰狼王</strong>
          <p>火系主力输出，当前战力最高，推荐优先投入副本经验、战骨与进化材料。</p>
          <span>下一步：补齐火系进化石后升境，提升 Boss 斩杀效率。</span>
        </div>
      </UiPanelCard>

      <UiPanelCard title="幻兽栏" wide>
        <div class="roster-grid">
          <article v-for="pet in pets" :key="pet.name" class="roster-item">
            <strong>{{ pet.name }}</strong>
            <span>{{ pet.level }}</span>
            <em>{{ pet.status }}</em>
            <small>综合战力 {{ pet.power }}</small>
          </article>
        </div>
      </UiPanelCard>

      <UiPanelCard title="上阵策略" wide>
        <UiChipGroup :items="strategies" tone="blue" min-width="150px" />
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-team-page {
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

.team-list,
.roster-grid {
  display: grid;
  gap: 12px;
}

.team-item,
.roster-item,
.focus-card {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.team-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.roster-grid {
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
}

.roster-item {
  display: grid;
  gap: 6px;
}

.focus-card p,
.team-item p {
  margin: 6px 0 0;
  color: #4b5563;
}

.focus-card span,
.roster-item small,
.subtle {
  color: #6b7280;
  font-size: 13px;
}

.power,
.roster-item em {
  font-style: normal;
  color: #1d4ed8;
  font-size: 13px;
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .team-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
