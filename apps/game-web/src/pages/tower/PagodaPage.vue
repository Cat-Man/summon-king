<template>
  <section class="pagoda-page">
    <div class="pagoda-hero">
      <p class="pagoda-label">通天塔挑战</p>
      <h1>登顶每一层，换取秘宝</h1>
      <p class="pagoda-subtitle">
        快速通关可获得战骨、战灵、魔魂器具。每一层的秘宝由游戏后台配置、可通过排行刷新。
      </p>
      <div class="pagoda-actions">
        <button class="primary">刷新层数</button>
        <button class="ghost">查看奖励</button>
      </div>
    </div>

    <div class="pagoda-grid">
      <article v-for="floor in floors" :key="floor.level" class="pagoda-card">
        <header>
          <strong>第 {{ floor.number }} 层</strong>
          <span>{{ floor.reward }}</span>
        </header>
        <p>{{ floor.description }}</p>
      </article>
    </div>
  </section>
</template>

<script setup>
const floors = [
  { number: 1, level: '入门', reward: '灵石礼包', description: '轻松摩擦，摸清节奏' },
  { number: 2, level: '试炼', reward: '战灵之鉴', description: '秘宝渐丰，小心反震' },
  { number: 3, level: '大荒', reward: '魔魂碎片', description: '怪物攻击力略增' },
  { number: 4, level: '终极', reward: '战骨锻造', description: '挑战极限速度与耐力' }
]
</script>

<style scoped>
:global(:root) {
  --pagoda-bg: radial-gradient(circle at top right, #fdf4f4, #fff5ef 50%, #f1e7d7);
  --pagoda-card: linear-gradient(135deg, rgba(255, 255, 255, 0.8), rgba(255, 230, 210, 0.95));
}

.pagoda-page {
  min-height: 100vh;
  padding: 2rem;
  background: var(--pagoda-bg);
  font-family: 'Playfair Display', 'STKaiti', serif;
  color: #1f1b1b;
}

.pagoda-hero {
  max-width: 720px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 28px;
  padding: 2.5rem;
  box-shadow: 0 40px 80px rgba(31, 27, 27, 0.12);
}

.pagoda-label {
  letter-spacing: 0.3em;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
  color: #c54b00;
  font-size: 0.85rem;
}

.pagoda-hero h1 {
  font-size: 2.6rem;
  margin-bottom: 0.75rem;
}

.pagoda-subtitle {
  font-size: 1rem;
  line-height: 1.8;
  color: #4c4238;
}

.pagoda-actions {
  margin-top: 1.5rem;
  display: flex;
  gap: 0.75rem;
}

.pagoda-actions button {
  border: none;
  border-radius: 999px;
  padding: 0.75rem 1.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.primary {
  background: #e34900;
  color: #fff;
}

.ghost {
  background: transparent;
  border: 1px solid #e34900;
  color: #e34900;
}

.pagoda-actions button:hover {
  transform: translateY(-2px);
}

.pagoda-grid {
  margin-top: 2.5rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
}

.pagoda-card {
  background: var(--pagoda-card);
  padding: 1.4rem;
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(31, 27, 27, 0.1);
  border: 1px solid rgba(197, 75, 0, 0.3);
}

.pagoda-card header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 1rem;
}

.pagoda-card p {
  color: #2b231f;
  line-height: 1.6;
  font-family: 'Source Sans Pro', 'PingFang SC', sans-serif;
}
</style>
