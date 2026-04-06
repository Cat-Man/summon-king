# 首页 Growth 自动刷新联动 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 在副本与成长动作发生后，首页和 growth 资源页能自动重新拉取最新数据，不必依赖手动刷新或重新进页面。

**Architecture:** 前端新增一个最小的 Pinia 同步 store，只维护一个递增版本号。会改变玩家资源/成长状态的动作在成功后 `touch()`；首页与需要联动的 growth 页面监听版本号变化后重新请求自己的 API。这样不引入全局缓存，也不修改后端协议，成本低且可持续扩展到更多页面。

**Tech Stack:** Vue 3 + Pinia + Vitest。

### Task 1: 建立资源同步 store

**Files:**
- Create: `apps/game-web/src/stores/resourceSync.ts`
- Test: 由页面测试间接覆盖

**Step 1: Write the failing test**
- 页面测试中先引入 `useResourceSyncStore`，让测试因为模块不存在而失败。

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/home/__tests__/HomePage.spec.ts src/pages/growth/__tests__/SpiritPage.spec.ts src/pages/growth/__tests__/SoulPage.spec.ts src/pages/growth/__tests__/BonePage.spec.ts src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: FAIL，提示 `@/stores/resourceSync` 不存在。

**Step 3: Write minimal implementation**
- 增加一个 Pinia store：
  - `version`
  - `touch()`

**Step 4: Run test to verify next failure**
- 再次运行同一命令
- Expected: 进入下一层失败，提示页面没有监听/没有触发同步。

### Task 2: 首页与 growth 页监听同步信号

**Files:**
- Modify: `apps/game-web/src/pages/home/HomePage.vue`
- Modify: `apps/game-web/src/pages/growth/SpiritPage.vue`
- Modify: `apps/game-web/src/pages/growth/SoulPage.vue`
- Modify: `apps/game-web/src/pages/growth/BonePage.vue`
- Test: `apps/game-web/src/pages/home/__tests__/HomePage.spec.ts`
- Test: `apps/game-web/src/pages/growth/__tests__/SpiritPage.spec.ts`
- Test: `apps/game-web/src/pages/growth/__tests__/SoulPage.spec.ts`
- Test: `apps/game-web/src/pages/growth/__tests__/BonePage.spec.ts`

**Step 1: Write the failing test**
- 首页测试：触发 `syncStore.touch()` 后重新拉 overview
- 战灵/魔魂页测试：触发 `syncStore.touch()` 后重新拉状态
- 战骨页测试：升级后同步版本增加

**Step 2: Run test to verify it fails**
Run the targeted Vitest command again.
Expected: FAIL，因为页面还没 watch 同步版本，也没在动作后发同步信号。

**Step 3: Write minimal implementation**
- 抽出各页面的 `load...()` 方法
- `watch(syncStore.version)` 后重新拉数据
- 成功动作后调用 `syncStore.touch()`

**Step 4: Run test to verify it passes**
Run the targeted Vitest command again.
Expected: PASS

### Task 3: 副本页发出同步信号

**Files:**
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Test: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 掷骰推进后断言 `syncStore.version` 增加

**Step 2: Run test to verify it fails**
Run the targeted Vitest command again.
Expected: FAIL，因为副本页还没触发同步信号。

**Step 3: Write minimal implementation**
- `rollForward` 成功后 `syncStore.touch()`
- `restartRun` 成功后也触发，保证首页副本摘要能同步

**Step 4: Run test to verify it passes**
Run the targeted Vitest command again.
Expected: PASS

### Task 4: 全量验证与提交

**Files:**
- Verify only

**Step 1: Run frontend verification**
Run: `make test-frontend`
Expected: PASS

**Step 2: Run backend verification**
Run: `make test-backend`
Expected: PASS

**Step 3: Run lint/typecheck**
Run: `cd apps/game-web && pnpm lint`
Expected: PASS

**Step 4: Commit**
```bash
git add apps/game-web/src/stores/resourceSync.ts apps/game-web/src/pages/home apps/game-web/src/pages/growth apps/game-web/src/pages/dungeons docs/plans/2026-04-04-首页-growth-自动刷新联动-implementation.md
git commit -m "feat: sync home and growth resources automatically"
```
