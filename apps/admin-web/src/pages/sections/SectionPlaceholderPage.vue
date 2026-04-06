<script setup lang="ts">
defineProps<{
  title: string
  summary: string
  status: string
}>()
</script>

<template>
  <section class="placeholder-page">
    <header class="placeholder-hero">
      <p class="placeholder-eyebrow">Module Placeholder</p>
      <h1>{{ title }}</h1>
      <p>{{ summary }}</p>
    </header>

    <article class="placeholder-card">
      <h2>接入状态</h2>
      <p>{{ status }}</p>
    </article>
  </section>
</template>

<style scoped>
.placeholder-page {
  display: grid;
  gap: 18px;
}

.placeholder-hero,
.placeholder-card {
  padding: 24px;
  border-radius: 24px;
  background: rgba(255, 252, 247, 0.92);
  border: 1px solid rgba(120, 91, 59, 0.12);
  box-shadow: 0 18px 36px rgba(76, 55, 33, 0.08);
}

.placeholder-eyebrow,
.placeholder-hero h1,
.placeholder-hero p,
.placeholder-card h2,
.placeholder-card p {
  margin: 0;
}

.placeholder-eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.18em;
  font-size: 12px;
  color: #8b6d4c;
}

.placeholder-hero h1 {
  margin-top: 10px;
  font-size: 36px;
  color: #24190f;
}

.placeholder-hero p:last-child,
.placeholder-card p {
  margin-top: 12px;
  line-height: 1.7;
  color: #5a4532;
}

.placeholder-card h2 {
  font-size: 20px;
  color: #352618;
}
</style>
