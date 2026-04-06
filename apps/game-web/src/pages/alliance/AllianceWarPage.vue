<template>
  <section class="war-page">
    <header class="hero-panel">
      <p class="hero-tag">Alliance War</p>
      <h1>盟战前线</h1>
      <p class="hero-copy">
        这一版先把联盟入口打通，盟战还保持占位状态。后续会在这里补报名、签到、对阵和战功兑换。
      </p>
    </header>

    <article class="card">
      <h2>当前阶段</h2>
      <p>筹备中</p>
      <p class="copy-text">等联盟基础稳定后，再把盟战节奏和战功结算接进来。</p>
      <RouterLink class="back-link" to="/alliance">返回联盟首页</RouterLink>
    </article>
  </section>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router"
</script>

<style scoped>
.war-page {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.hero-panel,
.card {
  padding: 24px;
  border-radius: 24px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background: rgba(10, 13, 22, 0.72);
}

.hero-panel {
  background:
    radial-gradient(circle at 20% 20%, rgba(248, 113, 113, 0.15), transparent 34%),
    radial-gradient(circle at 82% 18%, rgba(56, 189, 248, 0.16), transparent 32%),
    linear-gradient(180deg, rgba(15, 18, 30, 0.95), rgba(10, 12, 20, 0.95));
}

.hero-tag,
.hero-copy,
.copy-text {
  margin: 0;
  color: rgba(247, 239, 225, 0.76);
}

.hero-tag {
  letter-spacing: 0.24em;
  text-transform: uppercase;
  font-size: 12px;
}

.hero-panel h1,
.card h2,
.card p {
  margin: 0;
}

.hero-panel h1 {
  margin: 8px 0;
  font-size: clamp(28px, 4.8vw, 44px);
}

.card {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.back-link {
  display: inline-flex;
  align-self: flex-start;
  margin-top: 8px;
  padding: 10px 14px;
  border-radius: 14px;
  border: 1px solid rgba(247, 189, 120, 0.28);
  background: rgba(247, 189, 120, 0.12);
  color: #fff8f0;
  text-decoration: none;
}
</style>
