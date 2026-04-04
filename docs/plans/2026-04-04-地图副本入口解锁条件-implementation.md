# 地图副本入口解锁条件 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让地图中的副本入口真正带有解锁条件，未达到条件时前端显示禁用态，后端也拒绝直接进入。

**Architecture:** 继续复用当前地图入口模型，在 `MapDungeon` 上增加最小解锁字段 `unlock_spirit_power`。后端 `EnterDungeon` 在有钱包仓储时校验玩家灵力是否满足要求，不满足则返回明确的锁定错误；前端地图页读取当前钱包灵力，对未解锁入口展示条件并禁用按钮，避免用户误点。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Pinia + Vitest。

### Task 1: 后端为副本入口增加灵力解锁条件并在进入时校验

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/repository.go`
- Modify: `apps/backend/internal/modules/dungeon/service.go`
- Modify: `apps/backend/internal/modules/dungeon/http.go`
- Test: `apps/backend/internal/modules/dungeon/service_test.go`
- Test: `apps/backend/internal/modules/dungeon/http_test.go`

**Step 1: Write the failing test**
- 扩展地图测试，断言 `dungeons[]` 里返回 `unlock_spirit_power`
- 新增服务测试：
  - 默认钱包灵力 `100`
  - 进入要求 `120` 灵力的副本
  - 断言返回 `ErrDungeonLocked`
- 新增 handler 测试：
  - POST `/dungeon/enter`
  - 断言返回 `403`
  - 业务码不是通用 `5002`

**Step 2: Run test to verify it fails**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run 'Test(GetWorldMap_IncludesPrimaryDungeonBinding|EnterDungeon_RejectsLockedDungeonWithoutEnoughSpiritPower|Handler_EnterDungeonRejectsLockedDungeon)' -v`
Expected: FAIL，因为当前地图入口没有解锁字段，也没有进入校验。

**Step 3: Write minimal implementation**
- 给 `MapDungeon` 增加 `unlock_spirit_power`
- 内存地图给高阶副本配 `120` 灵力门槛
- `Service.EnterDungeon` 在有钱包仓储时校验条件
- `Handler.enterDungeon` 对锁定错误返回 `403`

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 2: 地图页显示并禁用未解锁副本入口

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/pages/maps/WorldMapPage.vue`
- Test: `apps/game-web/src/pages/maps/__tests__/WorldMapPage.spec.ts`

**Step 1: Write the failing test**
- mock 地图返回 `unlock_spirit_power`
- mock 玩家钱包灵力为 `100`
- 断言：
  - 页面显示 `需灵力 120`
  - 对应副本按钮为禁用态
  - 点击未锁定副本可以跳转
  - 点击锁定副本不会跳转

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/maps/__tests__/WorldMapPage.spec.ts`
Expected: FAIL，因为当前地图页没有加载钱包，也没有禁用逻辑。

**Step 3: Write minimal implementation**
- 更新前端 `WorldMap` 类型
- 地图页加载当前玩家钱包灵力
- 增加 `canEnterDungeon`
- 未解锁副本按钮添加 `disabled`
- 页面展示解锁条件文案

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 3: 回归验证与提交

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
git add apps/backend/internal/modules/dungeon apps/game-web/src/api/modules/dungeon.ts apps/game-web/src/pages/maps docs/plans/2026-04-04-地图副本入口解锁条件-implementation.md
git commit -m "feat: gate dungeon map entries by spirit power"
```
