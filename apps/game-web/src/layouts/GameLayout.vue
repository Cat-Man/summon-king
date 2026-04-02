<template>
  <div class="game-shell">
    <header class="shell-header">
      <div class="brand-block">
        <p class="brand-mark">ZHZW</p>
        <div>
          <h1>召唤之王</h1>
          <p>正式应用壳</p>
        </div>
      </div>

      <nav class="shell-nav" aria-label="主导航">
        <RouterLink to="/home">首页</RouterLink>
      </nav>
    </header>

    <main class="shell-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
:global(:root) {
  color-scheme: dark;
  font-family: "Noto Serif SC", "Hiragino Mincho ProN", "Songti SC", serif;
  background:
    radial-gradient(circle at top left, rgba(247, 189, 120, 0.18), transparent 28%),
    radial-gradient(circle at top right, rgba(65, 122, 220, 0.22), transparent 34%),
    linear-gradient(180deg, #120e1d, #07090f 60%);
}

:global(body) {
  margin: 0;
  min-width: 320px;
  background: transparent;
  color: #f7efe1;
}

:global(*) {
  box-sizing: border-box;
}

:global(a) {
  color: inherit;
  text-decoration: none;
}

.game-shell {
  min-height: 100vh;
  padding: 24px;
}

.shell-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin: 0 auto 24px;
  max-width: 1180px;
  padding: 18px 22px;
  border: 1px solid rgba(247, 239, 225, 0.12);
  border-radius: 24px;
  background: rgba(11, 13, 20, 0.68);
  backdrop-filter: blur(16px);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
}

.brand-block {
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand-mark {
  margin: 0;
  padding: 10px 12px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(244, 164, 96, 0.24), rgba(217, 119, 6, 0.48));
  font-family: "Avenir Next Condensed", "DIN Alternate", sans-serif;
  font-size: 13px;
  letter-spacing: 0.32em;
}

.brand-block h1,
.brand-block p {
  margin: 0;
}

.brand-block h1 {
  font-size: clamp(24px, 4vw, 34px);
}

.brand-block p {
  margin-top: 4px;
  color: rgba(247, 239, 225, 0.68);
  font-size: 14px;
}

.shell-nav {
  display: flex;
  align-items: center;
  gap: 12px;
}

.shell-nav a {
  padding: 10px 16px;
  border-radius: 999px;
  border: 1px solid rgba(247, 239, 225, 0.14);
  color: rgba(247, 239, 225, 0.78);
  transition: transform 160ms ease, border-color 160ms ease, color 160ms ease;
}

.shell-nav a.router-link-active {
  color: #fff8f0;
  border-color: rgba(247, 189, 120, 0.55);
  background: rgba(247, 189, 120, 0.12);
}

.shell-nav a:hover {
  transform: translateY(-1px);
  border-color: rgba(247, 239, 225, 0.35);
}

.shell-main {
  margin: 0 auto;
  max-width: 1180px;
}

@media (max-width: 720px) {
  .game-shell {
    padding: 16px;
  }

  .shell-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
