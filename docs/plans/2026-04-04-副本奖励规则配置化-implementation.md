# 副本奖励规则配置化 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 把副本奖励从 service 内硬编码常量收敛成可扩展规则，并让前端展示不同类型掉落。

**Architecture:** 在 dungeon 模块内部增加最小奖励规则层，由纯函数根据 `DungeonRun` 的 `status` / 楼层状态解析 `RollReward`。Service 只负责调用规则、写回 growth 钱包并回填 `last_reward` / `wallet_snapshot`。前端不参与规则判断，只直接展示接口返回的奖励标签与数值。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Vitest。

### Task 1: 后端奖励规则层

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/service.go`
- Test: `apps/backend/internal/modules/dungeon/service_test.go`

**Step 1: Write the failing test**
- 新增最小测试：
  - `TestResolveRollReward_OngoingStatus`
  - `TestResolveRollReward_BossStatus`
  - `TestRollDice_ExhaustedRunDoesNotGrantReward`
- 断言：
  - 普通层返回“怪物掉落”，固定灵力奖励
  - Boss 层返回“Boss掉落”，灵力和魂力更高
  - exhausted 不给奖励

**Step 2: Run test to verify it fails**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run 'Test(ResolveRollReward|RollDice_ExhaustedRunDoesNotGrantReward)' -v`
Expected: FAIL，因为当前没有奖励规则函数，也没有 exhausted 特判。

**Step 3: Write minimal implementation**
- 给 `RollReward` 增加 `label`
- 新增最小规则表/解析函数
- `RollDice` 改为走规则函数，再统一写回 wallet

**Step 4: Run test to verify it passes**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run 'Test(ResolveRollReward|RollDice_ExhaustedRunDoesNotGrantReward)' -v`
Expected: PASS

**Step 5: Commit**
- 先和前端一起提交

### Task 2: 前端展示奖励标签

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Test: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 更新副本页测试，断言掷骰后显示 `last_reward.label`
- 用 Boss 奖励样例断言页面可展示“Boss掉落”和 `魂力 +1`

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: FAIL，因为当前 reward 没有 label，页面也没展示。

**Step 3: Write minimal implementation**
- 扩展前端 `DungeonRun.last_reward.label`
- 页面把“本次掉落”标题改为后端返回 label 优先

**Step 4: Run test to verify it passes**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: PASS

**Step 5: Commit**
- 与 Task 1 一起提交

### Task 3: 全量验证

**Files:**
- Verify only

**Step 1: Run backend verification**
Run: `make test-backend`
Expected: PASS

**Step 2: Run frontend verification**
Run: `make test-frontend`
Expected: PASS

**Step 3: Run lint/typecheck**
Run: `cd apps/game-web && pnpm lint`
Expected: PASS

**Step 4: Commit**
```bash
git add apps/backend/internal/modules/dungeon apps/game-web/src/api/modules/dungeon.ts apps/game-web/src/pages/dungeons docs/plans/2026-04-04-副本奖励规则配置化-implementation.md
git commit -m "feat: configure dungeon reward rules"
```
