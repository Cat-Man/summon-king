<template>
  <section class="ranking-page">
    <div class="ranking-header">
      <p>排行榜</p>
      <h1>巅峰荣耀榜单</h1>
      <p>当前排行榜展示最近上榜的战力和更新时间，榜内玩家可获取额外奖励。</p>
    </div>

    <div class="ranking-table">
      <article v-for="entry in leaderboard" :key="entry.player_id">
        <span class="rank">NO. {{ entry.rank }}</span>
        <div>
          <h3>{{ entry.name }}</h3>
          <p>战力 {{ entry.score }} · {{ entry.updated }}</p>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
const leaderboard = [
  { player_id: 1001, name: '星痕丶苍穹', score: 13872, updated: '刚刚', rank: 1 },
  { player_id: 1002, name: '雨落荒原', score: 13700, updated: '2分钟前', rank: 2 },
  { player_id: 1003, name: '剑舞倾城', score: 13450, updated: '5分钟前', rank: 3 },
  { player_id: 1004, name: '千魂之主', score: 13310, updated: '8分钟前', rank: 4 }
]
</script>

<style scoped>
:global(:root) {
  --rank-bg: linear-gradient(180deg, #0d111b, #1b2035);
}

.ranking-page {
  min-height: 100vh;
  padding: 2.5rem;
  background: var(--rank-bg);
  color: #f8fbff;
  font-family: 'Sahitya', 'PingFang HK', serif;
}

.ranking-header h1 {
  font-size: 2.8rem;
  margin: 0.3rem 0;
}

.ranking-header p {
  max-width: 640px;
  color: rgba(248, 251, 255, 0.7);
}

.ranking-table {
  margin-top: 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.ranking-table article {
  border: 1px solid rgba(248, 251, 255, 0.07);
  border-radius: 18px;
  padding: 1rem 1.5rem;
  background: rgba(248, 251, 255, 0.03);
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rank {
  font-size: 1rem;
  letter-spacing: 0.2em;
  color: #8ae3ff;
  min-width: 80px;
}

.ranking-table h3 {
  margin: 0;
  font-size: 1.4rem;
}

.ranking-table p {
  margin: 0.2rem 0 0;
  color: rgba(248, 251, 255, 0.6);
}
</style>
