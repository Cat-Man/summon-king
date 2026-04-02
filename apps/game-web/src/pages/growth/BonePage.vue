<template>
  <section class="growth-page growth-page--bone">
    <header>
      <p class="tag">战骨</p>
      <h1>Pyramid Bone</h1>
      <p>铸魂之骨为成长基石，当前等级：5 级</p>
    </header>
    <div class="metrics">
      <article>
        <strong>力量</strong>
        <span>+18%</span>
      </article>
      <article>
        <strong>坚韧</strong>
        <span>+12%</span>
      </article>
    </div>
    <div class="panel">
      <p>当前战骨正在吸收雷火精华，下一层成就需要 12 个战骨石。</p>
    </div>
  </section>
</template>

<style scoped>
.growth-page {
  min-height: 60vh;
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(160deg, #1f1b4b, #2b2a7b);
  color: #fff;
  box-shadow: 0 20px 40px rgba(22, 15, 70, 0.5);
}
.growth-page header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.tag {
  background: rgba(255, 255, 255, 0.15);
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 12px;
  letter-spacing: 0.2em;
  text-transform: uppercase;
}
.metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin: 24px 0;
}
.metrics article {
  padding: 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.08);
}
.metrics strong {
  display: block;
  text-transform: uppercase;
  letter-spacing: 0.15em;
  font-size: 12px;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.06);
  line-height: 1.6;
}
</style>
