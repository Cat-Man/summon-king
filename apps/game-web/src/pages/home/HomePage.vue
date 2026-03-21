<template>
  <section class="home-page">
    <h1>召唤之王-首页</h1>
    <p>H5 玩家端初始化完成，游客登录与首页聚合接口已接入。</p>
    <ul>
      <li>登录接口：POST /api/v1/player/auth/login</li>
      <li>首页接口：GET /api/v1/player/home/index</li>
    </ul>
  </section>
</template>

<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
