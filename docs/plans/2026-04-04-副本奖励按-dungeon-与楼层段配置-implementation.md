# 副本奖励按 Dungeon 与楼层段配置 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让副本奖励不再只有一套全局常量，而是能按 `dungeon_id` 与楼层段返回不同掉落，并在副本页直接切换查看差异。

**Architecture:** 后端把奖励规则从简单 `status -> reward` 收敛成 `dungeon_id + floor range + status` 的最小规则表。`Service.RollDice` 继续只关心解析结果、写回 growth 钱包与回填 `last_reward`。前端副本页增加一个很轻的副本选择器，允许切换两档副本并重新进入，这样配置差异可以直接在页面展示，不需要额外管理后台或配置中心。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Pinia + Vitest。

### Task 1: 后端奖励规则按 dungeon 与楼层段解析

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/service.go`
- Test: `apps/backend/internal/modules/dungeon/service_test.go`

**Step 1: Write the failing test**
- 新增测试：
  - `TestResolveRollReward_DungeonTwoOngoingStatus`
  - `TestResolveRollReward_DungeonOneDeepBossStatus`
  - `TestRollDice_DungeonTwoUsesConfiguredReward`
- 断言：
  - `dungeon_id=2` 普通层掉落高于 `dungeon_id=1`
  - 深层 Boss 奖励高于浅层 Boss
  - 实际掷骰写回的钱包遵守配置

**Step 2: Run test to verify it fails**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run 'Test(ResolveRollReward_DungeonTwoOngoingStatus|ResolveRollReward_DungeonOneDeepBossStatus|RollDice_DungeonTwoUsesConfiguredReward)' -v`
Expected: FAIL，因为当前规则只看 `status`，不看 `dungeon_id` 和楼层段。

**Step 3: Write minimal implementation**
- 设计最小 reward rule 结构：
  - `dungeon_id`
  - `min_floor`
  - `max_floor`
  - `status`
  - `reward`
- 保留当前 `dungeon_id=1` 浅层奖励，避免已有行为回退
- 新增 `dungeon_id=2` 与深层 Boss 规则

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 2: 前端副本页增加副本切换入口

**Files:**
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Test: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 新增或扩展测试：
  - 选择 `dungeon_id=2`
  - 重新进入副本时调用 `enterDungeon(playerId, 2)`
  - 掷骰后页面显示 `dungeon 2` 的奖励标签

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: FAIL，因为当前页面没有副本切换，也始终进入 `dungeon_id=1`。

**Step 3: Write minimal implementation**
- 增加本地副本选项：
  - `1: 妖窟试炼`
  - `2: 寒渊裂隙`
- 页面增加切换按钮
- `restartRun` 与首次兜底进入都使用当前选择的 dungeon id

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 3: 全量验证与提交

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
git add apps/backend/internal/modules/dungeon apps/game-web/src/pages/dungeons docs/plans/2026-04-04-副本奖励按-dungeon-与楼层段配置-implementation.md
git commit -m "feat: configure dungeon rewards by tier"
```
