<template>
  <section class="growth-page growth-page--soul">
    <header>
      <p class="tag">魔魂</p>
      <h1>Soul Catcher</h1>
      <p>当前附魔残片 4 | 能量池 320</p>
    </header>
    <div class="grid">
      <article>
        <strong>搜魂速度</strong>
        <span>+22%</span>
      </article>
      <article>
        <strong>灵魂护盾</strong>
        <span>+18%</span>
      </article>
    </div>
    <p class="note">魔魂猎魂时间越短，掉落越丰厚。灵魂池刷新：12min</p>
  </section>
</template>

<style scoped>
.growth-page--soul {
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(180deg, #0c2c4a, #07233d);
  color: #e3f2ff;
  min-height: 52vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.grid article {
  padding: 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
}
.tag {
  font-size: 12px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
}
.note {
  max-width: 420px;
  line-height: 1.6;
}
</style>
