<template>
  <section class="cultivate">
    <div class="panel">
      <h1>修行池</h1>
      <p>点燃灵火，汲取天地之力，提升战灵和魔魂底蕴。</p>
      <div class="meter">
        <div class="meter-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <p class="meter-label">神力进度 {{ progress }}%</p>
      <div class="actions">
        <button class="primary">立刻修行</button>
        <button class="ghost">领取灵力</button>
      </div>
    </div>
    <div class="grid">
      <article>
        <h2>灵力储备</h2>
        <p>当前：120</p>
        <small>最多可储存 300 点</small>
      </article>
      <article>
        <h2>炼妖记录</h2>
        <p>常温炼妖完成</p>
        <small>下次可炼妖：12 分钟后</small>
      </article>
      <article>
        <h2>化仙池</h2>
        <p>提升免疫</p>
        <small>已解锁 2 个法阵</small>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
const progress = 62
</script>

<style scoped>
.cultivate {
  min-height: 100vh;
  padding: 36px;
  background: radial-gradient(circle at top, #2c284b, #0a0615 80%);
  color: #fbf8ff;
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.panel {
  padding: 28px;
  border-radius: 26px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.03);
  box-shadow: 0 20px 50px rgba(5, 5, 20, 0.45);
}
.meter {
  width: 100%;
  height: 6px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  margin: 18px 0;
  overflow: hidden;
}
.meter-fill {
  height: 100%;
  background: linear-gradient(90deg, #f4d35d, #f9576f);
  transition: width 0.3s;
}
.meter-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
}
.actions {
  margin-top: 20px;
  display: flex;
  gap: 12px;
}
.primary {
  border: none;
  border-radius: 14px;
  padding: 12px 18px;
  background: linear-gradient(135deg, #f9576f, #ffb347);
  color: #190d00;
  font-weight: 600;
  cursor: pointer;
}
.ghost {
  border-radius: 14px;
  padding: 12px 18px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.04);
  color: #fff;
  cursor: pointer;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 18px;
}
.grid article {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.03), rgba(255, 255, 255, 0.06));
}
.grid h2 {
  margin-bottom: 8px;
}
.grid small {
  color: rgba(255, 255, 255, 0.6);
}
</style>
