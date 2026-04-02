<template>
  <section class="growth-page growth-page--manor">
    <header>
      <p class="tag">庄园</p>
      <h1>Elder Grove</h1>
      <p>两个地块运行中，预计收获 3h 后</p>
    </header>
    <div class="plots">
      <article v-for="plot in plots" :key="plot.plot_id">
        <strong>地块 {{ plot.plot_id }}</strong>
        <span>{{ plot.state }}</span>
      </article>
    </div>
    <div class="panel">
      <p>庄园扩建可解锁更多地块，当前可领地 2 个，收获后可直接兑换灵晶。</p>
    </div>
  </section>
</template>

<script setup>
const plots = [
  { plot_id: 1, state: '智慧草苗 · 成长中' },
  { plot_id: 2, state: '灵根· 休眠中' }
]
</script>

<style scoped>
.growth-page--manor {
  padding: 32px;
  border-radius: 24px;
  background: linear-gradient(150deg, #1d4028, #0a1b10);
  color: #f0ffec;
  min-height: 52vh;
}
.plots {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 24px 0;
}
.plots article {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  justify-content: space-between;
}
.panel {
  padding: 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
}
.tag {
  font-size: 12px;
  letter-spacing: 0.3em;
  text-transform: uppercase;
}
</style>
