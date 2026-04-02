<template>
  <section class="spirit-tower-page">
    <div class="spirits-header">
      <p>战灵塔试炼</p>
      <h1>灵息共鸣，魂魄觉醒</h1>
      <p class="hint">每一次挑战都能刷新防守者组合，灵魂碎片与灵石随层级上升。</p>
      <div class="spirits-status">
        <span>当前层级：幽渊 7层</span>
        <span>可挑战次数：2/5</span>
      </div>
    </div>
    <div class="spirits-list">
      <article v-for="guard in guards" :key="guard.name">
        <h3>{{ guard.name }}</h3>
        <p>{{ guard.detail }}</p>
        <div class="guard-meta">
          <span>抗性：{{ guard.resistance }}</span>
          <span>技能：{{ guard.skill }}</span>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
const guards = [
  { name: '玄冰灵灵', detail: '冰封之心，群体减速附加破盾。', resistance: '霜/水', skill: '寒风裂地' },
  { name: '炎纹魂甲', detail: '生命越低伤害越高，附带反震。', resistance: '火/物', skill: '烈焰反噬' },
  { name: '雷影追魂', detail: '快速突袭目标、累积印记后爆发。', resistance: '雷/速', skill: '电光幻影' }
]
</script>

<style scoped>
:global(:root) {
  --spirit-gradient: linear-gradient(180deg, #0c0d2c, #26156a);
}

.spirit-tower-page {
  padding: 2rem;
  min-height: 100vh;
  background: radial-gradient(circle at 20% 20%, rgba(193, 198, 255, 0.2), transparent 60%), var(--spirit-gradient);
  color: #f6f5ff;
  font-family: 'Cinzel', 'STSong', serif;
}

.spirits-header {
  max-width: 640px;
  margin-bottom: 2rem;
}

.spirits-header p {
  margin: 0;
}

.spirits-header h1 {
  font-size: 2.8rem;
  margin-bottom: 0.5rem;
}

.hint {
  color: rgba(246, 245, 255, 0.75);
  line-height: 1.5;
}

.spirits-status {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
  font-family: 'Source Sans Pro', sans-serif;
  color: #d3d4ff;
}

.spirits-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.spirits-list article {
  background: rgba(255, 255, 255, 0.04);
  padding: 1.25rem;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.spirits-list h3 {
  margin-bottom: 0.4rem;
}

.guard-meta {
  margin-top: 0.8rem;
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
  color: #bbb8ff;
}
</style>
