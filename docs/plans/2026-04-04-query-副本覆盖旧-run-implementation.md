# Query 副本覆盖旧 Run Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让副本页在存在旧 run 时，仍然优先遵循地图或路由 query 传入的 `dungeon_id`，避免页面展示和实际副本状态脱节。

**Architecture:** 保持后端不变，只在前端副本页增加一个轻量的目标副本判定。页面加载状态后，如果发现当前 run 的 `dungeon_id` 与 query 指定目标不一致，就立刻调用 `enterDungeon` 重开目标副本；如果一致，则沿用已有 run。这样可以补上“地图入口决定当前副本”的最后一个行为缺口。

**Tech Stack:** Vue 3 + Vue Router + Pinia + Vitest。

### Task 1: 副本页在旧 run 与 query 冲突时自动切副本

**Files:**
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Test: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 新增测试：
  - `getDungeonStatus` 返回 `dungeon_id=1` 的已有 run
  - 路由 query 传入 `dungeon_id=2`
  - 页面初始化后应自动调用 `enterDungeon(playerId, 2)`
  - 页面展示切换为 `寒渊裂隙`

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: FAIL，因为当前逻辑只会采用已有 run 并覆盖 `selectedDungeonId`，不会重开目标副本。

**Step 3: Write minimal implementation**
- 提炼一个目标副本解析函数
- `refreshRun` 成功拿到已有 run 后，比较 `run.dungeon_id` 与目标副本
- 不一致时，直接调用 `enterDungeon` 重开目标副本，并回填页面状态
- 保留现有 404 兜底创建逻辑

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 2: 回归验证与提交

**Files:**
- Verify only

**Step 1: Run frontend verification**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: PASS

**Step 2: Run broader frontend verification**
Run: `make test-frontend`
Expected: PASS

**Step 3: Run lint/typecheck**
Run: `cd apps/game-web && pnpm lint`
Expected: PASS

**Step 4: Commit**
```bash
git add apps/game-web/src/pages/dungeons docs/plans/2026-04-04-query-副本覆盖旧-run-implementation.md
git commit -m "feat: honor route dungeon over stale run"
```
