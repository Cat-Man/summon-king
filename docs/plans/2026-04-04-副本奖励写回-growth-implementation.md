# 副本奖励写回 Growth Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让玩家在副本掷骰推进后，奖励能真实写回 growth 钱包，并在副本页直接看到最新掉落与资源变化。

**Architecture:** 继续复用 backend 内存仓储共享状态：dungeon service 在 `RollDice` 成功后，把奖励写入 growth repository；奖励规则保持最小化，避免引入背包/掉落表系统。前端不额外加新页面，只增强 `DungeonRunPage`，直接展示本次奖励与资源快照，首页继续通过 `home/overview` 读取共享钱包数据。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Pinia + Vitest。

### Task 1: 后端副本奖励写回 growth

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/service.go`
- Modify: `apps/backend/internal/modules/dungeon/service_test.go`
- Modify: `apps/backend/internal/modules/dungeon/http_test.go`（如需要 API 断言）

**Step 1: Write the failing test**
- 在 `service_test.go` 新增测试：`TestRollDice_UpdatesGrowthWalletRewards`
- 断言一次掷骰后：
  - 钱包灵力增加固定值
  - 到 Boss 层时魂力增加固定值
  - `DungeonRun` 返回 reward/wallet snapshot 字段

**Step 2: Run test to verify it fails**
- Run: `cd apps/backend && go test ./internal/modules/dungeon -run TestRollDice_UpdatesGrowthWalletRewards -v`
- Expected: FAIL，因为 `DungeonRun` 还没有奖励字段，`RollDice` 也没有写回 growth。

**Step 3: Write minimal implementation**
- 给 `DungeonRun` 增加最小奖励信息：`last_reward`、`wallet_snapshot`
- `Service.RollDice` 在仓储推进后：
  - 每次推进奖励固定灵力
  - Boss 层额外奖励固定魂力
  - 写入 growth repo
  - 读取最新 wallet 填回响应

**Step 4: Run test to verify it passes**
- Run: `cd apps/backend && go test ./internal/modules/dungeon -run TestRollDice_UpdatesGrowthWalletRewards -v`
- Expected: PASS

**Step 5: Commit**
- 先不单独 commit，和前端一起验证后合并成一个功能提交。

### Task 2: 前端副本页展示奖励闭环

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Modify: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 让 `DungeonRunPage.spec.ts` 断言：
  - 掷骰后页面显示本次奖励
  - 页面显示最新灵力/魂力快照

**Step 2: Run test to verify it fails**
- Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
- Expected: FAIL，因为当前页面没有奖励 UI，也没有对应字段。

**Step 3: Write minimal implementation**
- 扩展 `DungeonRun` 前端类型
- 在页面增加“本次掉落 / 当前资源”卡片
- 保持现有按钮与状态逻辑不变，只展示 API 返回数据

**Step 4: Run test to verify it passes**
- Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
- Expected: PASS

**Step 5: Commit**
- 与 Task 1 一起提交。

### Task 3: 全量验证与收口

**Files:**
- Verify only: backend/frontend related suites

**Step 1: Run backend verification**
- Run: `make test-backend`
- Expected: PASS

**Step 2: Run frontend verification**
- Run: `make test-frontend`
- Expected: PASS

**Step 3: Run type/lint verification**
- Run: `cd apps/game-web && pnpm lint`
- Expected: PASS

**Step 4: Commit**
```bash
git add apps/backend/internal/modules/dungeon apps/game-web/src/api/modules/dungeon.ts apps/game-web/src/pages/dungeons docs/plans/2026-04-04-副本奖励写回-growth-implementation.md
git commit -m "feat: connect dungeon rewards to growth loop"
```
