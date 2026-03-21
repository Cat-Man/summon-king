<script setup lang="ts">
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

const overview = [
  { label: '我的职位', value: '副盟主' },
  { label: '成员数', value: '32 / 40' },
  { label: '本周贡献', value: '1,280' },
  { label: '盟战阶段', value: '报名阶段' }
]

const buildings = [
  { name: '议事厅', level: 6, desc: '提升成员上限与公告栏容量' },
  { name: '焚火塔', level: 5, desc: '提升火能修行收益与房间数量' },
  { name: '宝库', level: 4, desc: '提升联盟资金保护与兑换上限' }
]

const activities = [
  '09:30 盟主发布了今晚盟战签到提醒',
  '10:10 成员「烈焰狼王」完成火能修行领取',
  '11:20 联盟建筑「焚火塔」升级完成'
]
</script>

<template>
  <section class="alliance-page">
    <UiPageHero
      eyebrow="次级主城"
      title="青焰盟"
      description="聚合公告、贡献、盟战状态、火能修行、建筑与联盟动态，方便玩家上线后一眼判断今天还能做什么。"
      tone="blue"
      meta-value="联盟资金 286,000"
    />

    <UiStatGrid :items="overview" min-width="140px" />

    <section class="content-grid">
      <UiPanelCard title="联盟公告">
        <p>今晚 20:00 锁定盟战名单，请先完成捐献并安排火能修行空房。</p>
      </UiPanelCard>

      <UiPanelCard title="火能修行">
        <div class="mini-row">
          <span>空闲房间：2 / 6</span>
          <span>可领取原石：480</span>
        </div>
        <p>推荐优先安排主力幻兽修行，领取后可直达养成页继续消耗。</p>
      </UiPanelCard>

      <UiPanelCard title="联盟建筑" wide>
        <div class="building-list">
          <div v-for="building in buildings" :key="building.name" class="building-item">
            <div>
              <strong>{{ building.name }} Lv.{{ building.level }}</strong>
              <p>{{ building.desc }}</p>
            </div>
            <span class="building-action">可升级</span>
          </div>
        </div>
      </UiPanelCard>

      <UiPanelCard title="联盟动态">
        <ul class="activity-list">
          <li v-for="item in activities" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="聊天室入口">
        <p>支持快速进入联盟频道，查看战前安排、成员求助与活动通知。</p>
        <button type="button" class="ghost-button">进入联盟频道</button>
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.alliance-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.mini-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  color: #374151;
  font-size: 14px;
}

.building-list,
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.building-item {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 14px;
  border-radius: 14px;
  background: #f8fafc;
}

.building-item p {
  margin: 6px 0 0;
  color: #6b7280;
  font-size: 14px;
}

.building-action {
  align-self: center;
  padding: 6px 10px;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 12px;
}

.activity-list li {
  padding: 12px 14px;
  border-radius: 14px;
  background: #f8fafc;
  line-height: 1.5;
}

.ghost-button {
  margin-top: 12px;
  padding: 10px 14px;
  border: 1px solid #c7d2fe;
  border-radius: 12px;
  background: #eef2ff;
  color: #3730a3;
  cursor: pointer;
}

@media (max-width: 768px) {
  .mini-row,
  .building-item {
    flex-direction: column;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
