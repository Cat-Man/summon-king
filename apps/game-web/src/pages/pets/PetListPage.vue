<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

const overview = [
  { label: '已拥有', value: '12' },
  { label: '未拥有', value: '36' },
  { label: '当前推荐地图', value: '青木林地' },
  { label: '获取来源', value: '副本 / 修行 / 礼包' }
]

const sources = [
  '地图副本：按地图掉落对应幻兽与召唤球',
  '修行：2 小时及以上有概率获得对应地图召唤球',
  '礼包与商城：用于补齐当前养成断档'
]

const pets = [
  {
    name: '烈焰狼王',
    status: '已拥有',
    map: '青木林地',
    skillPool: '烈焰撕咬 / 焚风追击',
    aptitude: '攻击 2200 / 3420',
    source: '地图副本 / 召唤球'
  },
  {
    name: '寒枝鹿灵',
    status: '已拥有',
    map: '白露港',
    skillPool: '寒枝缠绕 / 冰露护体',
    aptitude: '速度 920 / 1280',
    source: '地图副本 / 礼包'
  },
  {
    name: '霜牙虎',
    status: '未拥有',
    map: '霜风峡',
    skillPool: '裂地扑杀 / 寒牙冲击',
    aptitude: '攻击 1880 / 3010',
    source: '后续地图 / 召唤球'
  }
]
</script>

<template>
  <section class="pet-list-page">
    <UiPageHero
      eyebrow="图鉴与来源导航"
      title="幻兽图鉴"
      description="区分已拥有与未拥有幻兽，帮助玩家快速判断下一步去哪里刷、缺什么来源、哪些技能值得追。"
      tone="violet"
      meta-label="已拥有"
      meta-value="12 / 48"
    />

    <UiPanelCard title="图鉴统计">
      <UiStatGrid :items="overview" min-width="150px" />
    </UiPanelCard>

    <section class="grid">
      <UiPanelCard title="获取来源">
        <ul class="text-list">
          <li v-for="item in sources" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="图鉴列表" wide>
        <div class="catalog-grid">
          <article v-for="pet in pets" :key="pet.name" class="catalog-item">
            <header class="catalog-item__header">
              <strong>{{ pet.name }}</strong>
              <span :class="['status', pet.status === '已拥有' ? 'status--owned' : 'status--locked']">
                {{ pet.status }}
              </span>
            </header>
            <p><span>所属地图</span>{{ pet.map }}</p>
            <p><span>技能池</span>{{ pet.skillPool }}</p>
            <p><span>最低/满资质</span>{{ pet.aptitude }}</p>
            <p><span>获取来源</span>{{ pet.source }}</p>
          </article>
        </div>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.pet-list-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.grid {
  display: grid;
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

.catalog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 14px;
}

.catalog-item {
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
  display: grid;
  gap: 8px;
}

.catalog-item__header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.catalog-item p {
  margin: 0;
  display: grid;
  gap: 4px;
  color: #4b5563;
  line-height: 1.5;
}

.catalog-item p span {
  color: #6b7280;
  font-size: 12px;
}

.status {
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.status--owned {
  background: #dcfce7;
  color: #166534;
}

.status--locked {
  background: #ede9fe;
  color: #6d28d9;
}
</style>
