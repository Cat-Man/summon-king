<script setup lang="ts">
const entryCards = [
  {
    title: "配置中心",
    description: "幻兽、技能、地图副本、修行与成长配置",
    status: "当前只提供后台壳入口，配置读写、校验、发布与回滚链路待后续接入。",
    route: { name: "config" as const },
    enabled: true,
    accent: "accent-config",
  },
  {
    title: "GM 操作台",
    description: "发奖、补偿、封禁、解卡与运营处置",
    status: "真实 GM 执行与审计闭环待后续接入，本轮只保留主入口与能力边界。",
    route: { name: "gm" as const },
    enabled: true,
    accent: "accent-gm",
  },
  {
    title: "玩家查询",
    description: "玩家信息、资源、背包与战斗记录查询",
    status: "本轮不开放玩家检索与详情查询。",
    enabled: false,
    accent: "accent-muted",
  },
  {
    title: "运营报表",
    description: "新增、留存、付费、VIP 与资源产消",
    status: "报表与图表数据接入不在本轮范围。",
    enabled: false,
    accent: "accent-muted",
  },
  {
    title: "调度与审计",
    description: "任务补跑、执行状态、GM 日志与资源流水",
    status: "任务与审计后台待后续模块落地后接入。",
    enabled: false,
    accent: "accent-muted",
  },
] as const
</script>

<template>
  <section class="dashboard-page">
    <header class="hero-card">
      <div>
        <p class="hero-eyebrow">Admin Foundation Loop</p>
        <h1>运营后台</h1>
        <p class="hero-copy">
          本轮只完成后台工作台与核心入口壳，明确配置、GM、查询、报表与调度审计的能力边界。
        </p>
      </div>
      <div class="hero-metric">
        <strong>2 / 5</strong>
        <span>当前开放主入口</span>
      </div>
    </header>

    <div class="entry-grid">
      <article
        v-for="card in entryCards"
        :key="card.title"
        class="entry-card"
        :class="card.accent"
      >
        <div class="entry-head">
          <h2>{{ card.title }}</h2>
          <span class="entry-badge">{{ card.enabled ? "可进入" : "规划中" }}</span>
        </div>
        <p class="entry-description">{{ card.description }}</p>
        <p class="entry-status">{{ card.status }}</p>

        <RouterLink
          v-if="card.enabled"
          :to="card.route"
          class="entry-action"
        >
          进入模块
        </RouterLink>
        <span
          v-else
          class="entry-action entry-action-disabled"
        >
          待后续接入
        </span>
      </article>
    </div>
  </section>
</template>

<style scoped>
.dashboard-page {
  display: grid;
  gap: 20px;
}

.hero-card {
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 20px;
  padding: 28px;
  border-radius: 28px;
  background:
    linear-gradient(135deg, rgba(36, 24, 14, 0.92), rgba(83, 54, 29, 0.84)),
    linear-gradient(135deg, #1b1511, #6c4727);
  color: #f7f0e2;
  box-shadow: 0 28px 60px rgba(44, 28, 16, 0.16);
}

.hero-eyebrow,
.hero-card h1,
.hero-copy,
.hero-metric strong,
.hero-metric span {
  margin: 0;
}

.hero-eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 12px;
  opacity: 0.72;
}

.hero-card h1 {
  margin-top: 10px;
  font-size: 42px;
  line-height: 1.05;
}

.hero-copy {
  margin-top: 12px;
  max-width: 640px;
  color: rgba(247, 240, 226, 0.8);
  line-height: 1.65;
}

.hero-metric {
  min-width: 136px;
  padding: 18px 20px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.08);
  text-align: center;
}

.hero-metric strong {
  display: block;
  font-size: 34px;
}

.hero-metric span {
  display: block;
  margin-top: 8px;
  font-size: 13px;
  color: rgba(247, 240, 226, 0.74);
}

.entry-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.entry-card {
  padding: 22px;
  border-radius: 24px;
  background: rgba(255, 252, 247, 0.92);
  border: 1px solid rgba(130, 100, 62, 0.12);
  box-shadow: 0 18px 40px rgba(76, 55, 33, 0.08);
  display: grid;
  gap: 14px;
}

.entry-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.entry-head h2,
.entry-description,
.entry-status {
  margin: 0;
}

.entry-head h2 {
  font-size: 24px;
  color: #24190f;
}

.entry-badge {
  padding: 6px 10px;
  border-radius: 999px;
  font-size: 12px;
  background: rgba(54, 37, 19, 0.08);
  color: #6f563d;
}

.entry-description {
  color: #433223;
  line-height: 1.65;
}

.entry-status {
  color: #84644a;
  line-height: 1.6;
}

.entry-action {
  justify-self: start;
  padding: 10px 16px;
  border-radius: 999px;
  text-decoration: none;
  color: #fffdf8;
  background: linear-gradient(135deg, #6f4b2a, #a06d34);
}

.entry-action-disabled {
  background: rgba(67, 50, 35, 0.08);
  color: #7b6148;
}

.accent-config {
  background:
    linear-gradient(180deg, rgba(248, 241, 228, 0.96), rgba(255, 252, 247, 0.96)),
    rgba(255, 252, 247, 0.92);
}

.accent-gm {
  background:
    linear-gradient(180deg, rgba(244, 231, 213, 0.96), rgba(255, 252, 247, 0.96)),
    rgba(255, 252, 247, 0.92);
}

.accent-muted {
  opacity: 0.88;
}

@media (max-width: 960px) {
  .hero-card {
    flex-direction: column;
    align-items: start;
  }

  .entry-grid {
    grid-template-columns: 1fr;
  }
}
</style>
