<template>
  <section class="growth-page growth-page--spirit">
    <header>
      <p class="tag">战灵</p>
      <h1>Spirit Forge</h1>
      <p>灵石数量 8 | 洗炼剩余次数 3</p>
    </header>
    <div class="banner">
      <span>免费洗炼优先消耗</span>
      <span>保底属性提升 +5%</span>
    </div>
    <button class="primary">立刻洗炼</button>
  </section>
</template>

<style scoped>
.growth-page--spirit {
  padding: 32px;
  border-radius: 24px;
  background: radial-gradient(circle at top right, #f7ba38, #dd2e6d 60%);
  color: #171717;
  min-height: 52vh;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.tag {
  font-size: 12px;
  letter-spacing: 0.25em;
  text-transform: uppercase;
}
.banner {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.4);
  display: flex;
  justify-content: space-between;
}
.primary {
  padding: 14px 24px;
  border: none;
  border-radius: 12px;
  background: #171717;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
}
</style>
