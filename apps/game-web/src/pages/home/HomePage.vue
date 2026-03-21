<script setup lang="ts">
import UiChipGroup from '@/components/ui/UiChipGroup.vue'
import UiPageHero from '@/components/ui/UiPageHero.vue'
import UiPanelCard from '@/components/ui/UiPanelCard.vue'
import UiStatGrid from '@/components/ui/UiStatGrid.vue'

const dailyTodos = [
  { title: '签到状态', value: '今日未签', action: '前往签到' },
  { title: '当前修行状态', value: '还有 18 分钟可领取', action: '查看修行' },
  { title: '当前推荐副本', value: '青木林地 · Boss 可挑战', action: '进入副本' }
]

const resources = [
  { label: '等级', value: 'Lv.36' },
  { label: '战力', value: '18,620' },
  { label: '活力', value: '86 / 120' },
  { label: '铜钱', value: '286,400' },
  { label: '元宝', value: '1,280' },
  { label: '声望', value: '540' }
]

const entries = ['世界地图', '联盟', '幻兽', '背包', '竞技场', '庄园', '修行', '排行']

const messages = [
  '世界消息：青木林地今日双倍经验已开启',
  '联盟消息：今晚 20:00 盟战锁定名单',
  '系统消息：VIP 每日宝箱可领取'
]
</script>

<template>
  <section class="home-page">
    <UiPageHero
      eyebrow="今日工作台"
      title="召唤之王"
      description="第一屏直接告诉玩家今天能做什么、有哪些收益可领、当前主推进目标是什么。"
      tone="navy"
      meta-label="收益聚合"
      meta-value="3 项待处理"
    />

    <UiPanelCard title="每日必做">
      <div class="todo-list">
        <article v-for="item in dailyTodos" :key="item.title" class="todo-item">
          <strong>{{ item.title }}</strong>
          <span>{{ item.value }}</span>
          <em>{{ item.action }}</em>
        </article>
      </div>
    </UiPanelCard>

    <section class="content-grid">
      <UiPanelCard title="资源总览">
        <UiStatGrid :items="resources" min-width="110px" />
      </UiPanelCard>

      <UiPanelCard title="消息流入口">
        <ul class="message-list">
          <li v-for="item in messages" :key="item">{{ item }}</li>
        </ul>
      </UiPanelCard>

      <UiPanelCard title="功能矩阵" wide>
        <UiChipGroup :items="entries" tone="blue" min-width="100px" />
      </UiPanelCard>
    </section>
  </section>
</template>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  color: #1f2937;
}

.todo-list,
.content-grid {
  display: grid;
  gap: 14px;
}

.todo-list {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.todo-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
  border-radius: 16px;
  background: #f8fafc;
}

.todo-item span {
  color: #4b5563;
}

.todo-item em {
  font-style: normal;
  color: #2563eb;
  font-size: 13px;
}

.content-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.message-list {
  margin: 0;
  padding-left: 18px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: #4b5563;
}

@media (max-width: 768px) {
  .content-grid {
    grid-template-columns: 1fr;
  }
}
</style>
