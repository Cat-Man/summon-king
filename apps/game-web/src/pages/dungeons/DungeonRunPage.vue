<template>
  <section class="dungeon-page">
    <div class="dungeon-meta">
      <h1>地下城跑动</h1>
      <p>掷骰进入楼层，打响各路 Boss 抢夺装备与材料。</p>
      <div class="stat-grid">
        <article>
          <h2>{{ run.current_floor }}</h2>
          <p>当前层</p>
        </article>
        <article>
          <h2>{{ run.remain_dice }}</h2>
          <p>剩余骰子</p>
        </article>
        <article>
          <h2>{{ run.status }}</h2>
          <p>状态</p>
        </article>
      </div>
    </div>
    <div class="action-panel">
      <button class="action-btn primary">掷骰推进</button>
      <button class="action-btn">领取 Boss 奖励</button>
      <button class="action-btn ghost">查看战报</button>
    </div>
    <div class="timeline">
      <p>进度</p>
      <div class="floors">
        <div v-for="n in 6" :key="n" class="floor" :class="{ boss: n % 5 === 0 }">
          <span>{{ n }}</span>
          <small>{{ n % 5 === 0 ? 'Boss' : '怪物' }}</small>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const run = {
  current_floor: 3,
  remain_dice: 12,
  status: 'ongoing'
}
</script>

<style scoped>
.dungeon-page {
  min-height: 100vh;
  padding: 32px;
  background: linear-gradient(180deg, #120c1f, #090515 70%);
  color: #f8fbff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.dungeon-meta {
  border-radius: 24px;
  padding: 28px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 30px 60px rgba(5, 5, 20, 0.6);
}
.stat-grid {
  margin-top: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}
.stat-grid article {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.stat-grid h2 {
  font-size: 32px;
  margin-bottom: 4px;
}
.action-panel {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}
.action-btn {
  flex: 1;
  min-width: 180px;
  border-radius: 15px;
  padding: 14px 20px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  background: transparent;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, background 0.2s ease;
}
.action-btn.primary {
  background: linear-gradient(135deg, #ffb347, #ffcc33);
  color: #1b1000;
  border: none;
}
.action-btn.ghost {
  border-color: rgba(255, 255, 255, 0.3);
}
.action-btn:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.1);
}
.timeline {
  border-radius: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.floors {
  margin-top: 12px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(70px, 1fr));
  gap: 12px;
}
.floor {
  border-radius: 12px;
  padding: 12px;
  text-align: center;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
.floor.boss {
  background: linear-gradient(180deg, rgba(255, 94, 98, 0.1), rgba(255, 94, 98, 0.3));
  border-color: rgba(255, 94, 98, 0.7);
}
.floor small {
  display: block;
  font-size: 12px;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.7);
}
</style>
