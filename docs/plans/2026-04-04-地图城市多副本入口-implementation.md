# 地图城市多副本入口 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让地图页中的每座城市不再只映射一个默认副本，而是展示多个可进入副本入口，用户点击哪个入口就进入哪个副本。

**Architecture:** 后端把 `MapCity` 从单个 `dungeon_id/dungeon_name` 升级成 `dungeons[]`，每个元素只包含最小的副本 id 与名称。前端地图页改为在城市卡片中渲染副本入口列表，直接按点击的副本按钮跳到 `/dungeon?dungeon_id=<id>`，继续复用已经完成的副本页 query 联动逻辑。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Vue Router + Vitest。

### Task 1: 后端地图城市返回多个副本入口

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/repository.go`
- Test: `apps/backend/internal/modules/dungeon/service_test.go`

**Step 1: Write the failing test**
- 调整 `TestGetWorldMap_IncludesPrimaryDungeonBinding`
- 断言：
  - `晨曦城` 返回两个副本入口
  - 第二个入口为 `dungeon_id=2`
  - `霞光堡` 也返回多个副本入口

**Step 2: Run test to verify it fails**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run TestGetWorldMap_IncludesPrimaryDungeonBinding -v`
Expected: FAIL，因为当前 `MapCity` 仍然是单个副本字段。

**Step 3: Write minimal implementation**
- 增加 `MapDungeon` 结构
- `MapCity` 改成 `dungeons []`
- 内存世界地图给两座城各配两个可进入副本

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 2: 地图页渲染多个副本按钮并按按钮跳转

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/pages/maps/WorldMapPage.vue`
- Test: `apps/game-web/src/pages/maps/__tests__/WorldMapPage.spec.ts`

**Step 1: Write the failing test**
- 调整地图页测试：
  - mock 城市返回 `dungeons[]`
  - 页面展示两条副本名
  - 点击城市中的第二个副本按钮
  - 断言路由跳到 `{ name: "dungeon", query: { dungeon_id: "2" } }`

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/maps/__tests__/WorldMapPage.spec.ts`
Expected: FAIL，因为当前页面只支持单个副本按钮。

**Step 3: Write minimal implementation**
- 更新 `WorldMap` 前端类型
- 地图页展示 `city.dungeons`
- 给副本按钮加稳定选择器
- 点击按钮按所选副本跳转

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
git add apps/backend/internal/modules/dungeon apps/game-web/src/api/modules/dungeon.ts apps/game-web/src/pages/maps docs/plans/2026-04-04-地图城市多副本入口-implementation.md
git commit -m "feat: show multiple dungeon entries on map"
```
