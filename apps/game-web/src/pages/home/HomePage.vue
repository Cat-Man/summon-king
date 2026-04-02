<template>
  <section class="home-page">
    <div class="hero-panel">
      <p class="hero-tag">召唤之王</p>
      <h2>首页总览</h2>
      <p class="hero-copy">
        正式前端应用壳已经接通。当前阶段先保留主循环入口和今日概览占位，下一步把地图、副本、修行接到真实 API。
      </p>
      <div class="hero-actions">
        <span>版本阶段：底座与首条闭环</span>
        <span>后端接口：/api/v1</span>
      </div>
    </div>

    <div class="overview-grid">
      <article v-for="card in cards" :key="card.title" class="overview-card">
        <p>{{ card.eyebrow }}</p>
        <h3>{{ card.title }}</h3>
        <strong>{{ card.value }}</strong>
        <span>{{ card.description }}</span>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
const cards = [
  {
    eyebrow: "今日概览",
    title: "主城状态",
    value: "待接 API",
    description: "登录、首页聚合与钱包摘要将在下一阶段接入。",
  },
  {
    eyebrow: "冒险链路",
    title: "地图 / 副本",
    value: "骨架就绪",
    description: "已有页面原型，待路由挂载与真实数据联调。",
  },
  {
    eyebrow: "成长链路",
    title: "修行 / 战灵",
    value: "准备中",
    description: "保留视觉稿，后续统一接入会话和 API 客户端。",
  },
]
</script>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hero-panel {
  padding: 32px;
  border-radius: 30px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  background:
    radial-gradient(circle at top right, rgba(237, 164, 93, 0.16), transparent 24%),
    linear-gradient(180deg, rgba(22, 25, 37, 0.92), rgba(12, 15, 24, 0.9));
  box-shadow: 0 30px 70px rgba(0, 0, 0, 0.35);
}

.hero-tag {
  display: inline-flex;
  margin: 0 0 10px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(247, 189, 120, 0.16);
  color: #f7d8a8;
  letter-spacing: 0.24em;
  font-size: 12px;
  text-transform: uppercase;
}

.hero-panel h2,
.hero-copy {
  margin: 0;
}

.hero-panel h2 {
  font-size: clamp(32px, 5vw, 52px);
  line-height: 1.05;
}

.hero-copy {
  margin-top: 14px;
  max-width: 720px;
  color: rgba(247, 239, 225, 0.72);
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
}

.hero-actions span {
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(247, 239, 225, 0.06);
  color: rgba(247, 239, 225, 0.8);
  font-family: "Avenir Next", "PingFang SC", sans-serif;
  font-size: 14px;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 18px;
}

.overview-card {
  padding: 22px;
  border-radius: 22px;
  border: 1px solid rgba(247, 239, 225, 0.1);
  background: rgba(12, 14, 22, 0.8);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.25);
}

.overview-card p,
.overview-card h3,
.overview-card strong,
.overview-card span {
  display: block;
}

.overview-card p {
  margin: 0 0 10px;
  color: rgba(247, 189, 120, 0.78);
  letter-spacing: 0.16em;
  font-size: 12px;
  text-transform: uppercase;
}

.overview-card h3 {
  margin: 0;
  font-size: 22px;
}

.overview-card strong {
  margin-top: 18px;
  font-size: 28px;
}

.overview-card span {
  margin-top: 12px;
  color: rgba(247, 239, 225, 0.68);
  line-height: 1.6;
}

@media (max-width: 720px) {
  .hero-panel {
    padding: 24px;
  }
}
</style>
