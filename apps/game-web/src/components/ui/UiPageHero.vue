<script setup lang="ts">
withDefaults(
  defineProps<{
    eyebrow: string
    title: string
    description: string
    tone?: 'navy' | 'blue' | 'violet' | 'amber' | 'cyan' | 'green' | 'teal' | 'orange'
    metaLabel?: string
    metaValue?: string
  }>(),
  {
    tone: 'blue',
    metaLabel: '',
    metaValue: ''
  }
)
</script>

<template>
  <header class="ui-page-hero" :class="`ui-page-hero--${tone}`">
    <div class="ui-page-hero__body">
      <p class="ui-page-hero__eyebrow">{{ eyebrow }}</p>
      <h1 class="ui-page-hero__title">{{ title }}</h1>
      <p class="ui-page-hero__description">{{ description }}</p>
    </div>

    <div v-if="metaValue" class="ui-page-hero__meta">
      <span v-if="metaLabel">{{ metaLabel }}</span>
      <strong>{{ metaValue }}</strong>
    </div>
  </header>
</template>

<style scoped>
.ui-page-hero {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 20px;
  border-radius: 20px;
  color: #fff;
  box-shadow: 0 10px 30px rgb(15 23 42 / 8%);
}

.ui-page-hero--navy { background: linear-gradient(135deg, #0f172a, #1d4ed8); }
.ui-page-hero--blue { background: linear-gradient(135deg, #1d4ed8, #4338ca); }
.ui-page-hero--violet { background: linear-gradient(135deg, #7c3aed, #db2777); }
.ui-page-hero--amber { background: linear-gradient(135deg, #f59e0b, #f97316); }
.ui-page-hero--cyan { background: linear-gradient(135deg, #0891b2, #2563eb); }
.ui-page-hero--green { background: linear-gradient(135deg, #16a34a, #22c55e); }
.ui-page-hero--teal { background: linear-gradient(135deg, #065f46, #0f766e); }
.ui-page-hero--orange { background: linear-gradient(135deg, #7c2d12, #ea580c); }

.ui-page-hero__eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  letter-spacing: 0.12em;
  opacity: 0.82;
}

.ui-page-hero__title {
  margin: 0 0 10px;
  font-size: 28px;
}

.ui-page-hero__description {
  margin: 0;
  line-height: 1.6;
}

.ui-page-hero__meta {
  display: flex;
  flex-direction: column;
  align-self: flex-start;
  gap: 6px;
  padding: 10px 14px;
  border-radius: 16px;
  background: rgb(255 255 255 / 16%);
  white-space: nowrap;
}

.ui-page-hero__meta strong {
  font-size: 28px;
}

@media (max-width: 768px) {
  .ui-page-hero {
    flex-direction: column;
  }
}
</style>
