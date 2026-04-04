# 地图入口决定当前副本 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 让地图页中的城市入口不再只是静态按钮，而是能决定进入哪个副本，并把该副本透传到副本页作为默认目标。

**Architecture:** 后端先把每个城市绑定一个最小主副本信息，直接挂在 `WorldMap.Cities` 中返回给前端，避免额外请求或复杂配置。前端地图页点击城市后直接跳转到 `/dungeon?dungeon_id=<id>`，副本页从路由 query 初始化当前副本，并在 query 变化时刷新当前选择，这样地图和副本形成真实联动。

**Tech Stack:** Go + Gin + 内存仓储；Vue 3 + Vue Router + Pinia + Vitest。

### Task 1: 后端地图数据带出主副本绑定

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/repository.go`
- Test: `apps/backend/internal/modules/dungeon/service_test.go`

**Step 1: Write the failing test**
- 新增测试 `TestGetWorldMap_IncludesPrimaryDungeonBinding`
- 断言：
  - 地图至少返回城市数据
  - `晨曦城` 绑定 `dungeon_id=1`
  - `霞光堡` 绑定 `dungeon_id=2`
  - 同时带出副本名，避免前端再硬编码映射

**Step 2: Run test to verify it fails**
Run: `cd apps/backend && go test ./internal/modules/dungeon -run TestGetWorldMap_IncludesPrimaryDungeonBinding -v`
Expected: FAIL，因为当前 `MapCity` 还没有副本绑定字段。

**Step 3: Write minimal implementation**
- 在 `MapCity` 增加：
  - `dungeon_id`
  - `dungeon_name`
- 在内存地图仓储里补两座城的主副本绑定

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 2: 地图页点击城市跳转到带 query 的副本页

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/pages/maps/WorldMapPage.vue`
- Test: `apps/game-web/src/pages/maps/__tests__/WorldMapPage.spec.ts`

**Step 1: Write the failing test**
- 扩展地图页测试：
  - mock `getWorldMap` 返回 `dungeon_id` 与 `dungeon_name`
  - 点击 `传送到 晨曦城`
  - 断言调用路由跳转到 `{ name: "dungeon", query: { dungeon_id: "1" } }`
  - 页面上能看到副本名提示

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/maps/__tests__/WorldMapPage.spec.ts`
Expected: FAIL，因为当前页面按钮没有点击逻辑，也没有副本绑定字段。

**Step 3: Write minimal implementation**
- 更新 `WorldMap` 类型，补城市副本字段
- 地图页接入 `useRouter`
- 城市卡片展示主副本名
- 按钮点击跳到副本页并附带 `dungeon_id` query

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 3: 副本页从 query 初始化并切换当前副本

**Files:**
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Test: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`

**Step 1: Write the failing test**
- 新增测试：
  - 路由 query 为 `dungeon_id=2`
  - 首次 `getDungeonStatus` 返回 404 时，应调用 `enterDungeon(playerId, 2)`
  - 页面显示当前副本为 `寒渊裂隙`

**Step 2: Run test to verify it fails**
Run: `cd apps/game-web && pnpm test -- src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
Expected: FAIL，因为当前副本页不会从路由 query 读取默认副本。

**Step 3: Write minimal implementation**
- 接入 `useRoute`
- 增加解析 query 的轻量函数，只接受已知副本 id
- 初始化时优先从 query 设定 `selectedDungeonId`
- 监听 query 变化，必要时同步当前副本选择

**Step 4: Run test to verify it passes**
Run the same command again.
Expected: PASS

### Task 4: 回归验证与提交

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
git add apps/backend/internal/modules/dungeon apps/game-web/src/api/modules/dungeon.ts apps/game-web/src/pages/maps apps/game-web/src/pages/dungeons docs/plans/2026-04-04-地图入口决定当前副本-implementation.md
git commit -m "feat: link map entries to dungeon selection"
```
