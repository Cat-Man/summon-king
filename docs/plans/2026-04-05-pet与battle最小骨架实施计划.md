# pet / battle 最小骨架 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 新增最小 `pet` / `battle` 内部模块，让 `dungeon`、`tower`、`arena` 不再只返回各自散落的伪战斗结果，而是统一输出阵容快照驱动的战斗摘要。

**Architecture:** 第一批不做完整幻兽图鉴、上阵编辑、战斗回放接口，只做后端内部可复用骨架。`pet` 提供统一玩家战斗队伍快照，默认给新玩家一只 starter 幻兽；`battle` 提供统一战斗摘要计算，允许当前玩法先用“简化数值判定 + 预设胜负”过渡；`dungeon / tower / arena` 接入同一套 `pet` 与 `battle` 服务，并保持原有奖励返回结构不破坏前端闭环。

**Tech Stack:** Go 1.22, Gin, 内存仓储, Vue 3, TypeScript, Vitest

### Task 1: 新建最小 pet 模块与测试

**Files:**
- Create: `apps/backend/internal/modules/pet/domain.go`
- Create: `apps/backend/internal/modules/pet/repository.go`
- Create: `apps/backend/internal/modules/pet/service.go`
- Create: `apps/backend/internal/modules/pet/service_test.go`

**Step 1: Write the failing tests**

在 `apps/backend/internal/modules/pet/service_test.go` 写：

```go
func TestService_GetBattleTeamReturnsStarterPet(t *testing.T) {
    ctx := context.Background()
    svc := NewService(NewMemoryRepository())

    team, err := svc.GetBattleTeam(ctx, 1001)

    if err != nil {
        t.Fatalf("expected battle team success, got %v", err)
    }
    if team.PlayerID != 1001 {
        t.Fatalf("expected player 1001, got %d", team.PlayerID)
    }
    if len(team.Pets) != 1 {
        t.Fatalf("expected 1 starter pet, got %d", len(team.Pets))
    }
    if team.Pets[0].Slot != 1 {
        t.Fatalf("expected starter slot 1, got %d", team.Pets[0].Slot)
    }
    if team.Pets[0].Power <= 0 {
        t.Fatalf("expected starter power > 0, got %d", team.Pets[0].Power)
    }
}

func TestService_GetBattleTeamAggregatesTeamPower(t *testing.T) {
    ctx := context.Background()
    svc := NewService(NewMemoryRepository())

    team, err := svc.GetBattleTeam(ctx, 1002)

    if err != nil {
        t.Fatalf("expected battle team success, got %v", err)
    }
    if team.TotalPower != team.Pets[0].Power {
        t.Fatalf("expected total power %d, got %d", team.Pets[0].Power, team.TotalPower)
    }
}
```

**Step 2: Run tests to verify they fail**

Run:
```bash
cd apps/backend && go test ./internal/modules/pet -v
```

Expected: FAIL because `pet` 模块还不存在

**Step 3: Write minimal implementation**

- `domain.go`：
  - `BattlePet`
  - `TeamSnapshot`
- `repository.go`：
  - `Repository`
  - `MemoryRepository`
  - 默认 starter 幻兽一只，固定在 1 号位
- `service.go`：
  - `NewService(repo)`
  - `GetBattleTeam(ctx, playerID)`
  - 负责汇总 `TotalPower`

**Step 4: Run tests to verify they pass**

Run:
```bash
cd apps/backend && go test ./internal/modules/pet -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/backend/internal/modules/pet
git commit -m "feat: add minimal pet team snapshot service"
```

### Task 2: 新建最小 battle 模块与测试

**Files:**
- Create: `apps/backend/internal/modules/battle/domain.go`
- Create: `apps/backend/internal/modules/battle/service.go`
- Create: `apps/backend/internal/modules/battle/service_test.go`

**Step 1: Write the failing tests**

在 `apps/backend/internal/modules/battle/service_test.go` 写：

```go
func TestService_ResolveChoosesHigherPowerWithoutPreset(t *testing.T) {
    ctx := context.Background()
    svc := NewService()

    result, err := svc.Resolve(ctx, Request{
        Type: "tower",
        Attacker: TeamSnapshot{PlayerID: 1001, TotalPower: 180},
        Defender: TeamSnapshot{PlayerID: 0, TotalPower: 120},
    })

    if err != nil {
        t.Fatalf("expected resolve success, got %v", err)
    }
    if result.Result != "success" {
        t.Fatalf("expected success, got %s", result.Result)
    }
}

func TestService_ResolveUsesPresetOutcomeWhenProvided(t *testing.T) {
    ctx := context.Background()
    svc := NewService()

    result, err := svc.Resolve(ctx, Request{
        Type: "arena",
        PresetResult: "fail",
        Attacker: TeamSnapshot{PlayerID: 1001, TotalPower: 180},
        Defender: TeamSnapshot{PlayerID: 0, TotalPower: 120},
    })

    if err != nil {
        t.Fatalf("expected resolve success, got %v", err)
    }
    if result.Result != "fail" {
        t.Fatalf("expected fail, got %s", result.Result)
    }
    if result.BattleType != "arena" {
        t.Fatalf("expected battle type arena, got %s", result.BattleType)
    }
}
```

**Step 2: Run tests to verify they fail**

Run:
```bash
cd apps/backend && go test ./internal/modules/battle -v
```

Expected: FAIL because `battle` 模块还不存在

**Step 3: Write minimal implementation**

- `domain.go`：
  - `Request`
  - `Summary`
  - `Result`
- `service.go`：
  - `NewService()`
  - `Resolve(ctx, Request)`
  - 规则：
    - 若 `PresetResult` 非空，直接采用
    - 否则按 `Attacker.TotalPower >= Defender.TotalPower` 判定 `success/fail`
    - 返回统一战斗摘要：`battle_type / result / rounds / attacker_power / defender_power`

**Step 4: Run tests to verify they pass**

Run:
```bash
cd apps/backend && go test ./internal/modules/battle -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/backend/internal/modules/battle
git commit -m "feat: add minimal battle summary service"
```

### Task 3: 让 dungeon / tower / arena 接入统一阵容与战斗摘要

**Files:**
- Modify: `apps/backend/internal/modules/dungeon/domain.go`
- Modify: `apps/backend/internal/modules/dungeon/service.go`
- Modify: `apps/backend/internal/modules/dungeon/service_test.go`
- Modify: `apps/backend/internal/modules/tower/domain.go`
- Modify: `apps/backend/internal/modules/tower/service.go`
- Modify: `apps/backend/internal/modules/tower/service_test.go`
- Modify: `apps/backend/internal/modules/arena/domain.go`
- Modify: `apps/backend/internal/modules/arena/service.go`
- Modify: `apps/backend/internal/modules/arena/service_test.go`
- Modify: `apps/backend/internal/modules/home/service_test.go`
- Modify: `apps/backend/internal/modules/home/http_test.go`
- Modify: `apps/backend/internal/modules/ranking/service_test.go`
- Modify: `apps/backend/internal/modules/ranking/http_test.go`
- Modify: `apps/backend/internal/bootstrap/router.go`

**Step 1: Write the failing tests**

在三类玩法现有测试里补统一断言：

```go
if result.Battle.BattleType != "tower" {
    t.Fatalf("expected tower battle summary, got %s", result.Battle.BattleType)
}
if result.Battle.Result == "" {
    t.Fatal("expected battle result")
}
if result.Battle.AttackerPower <= 0 {
    t.Fatalf("expected attacker power > 0, got %d", result.Battle.AttackerPower)
}
```

`dungeon` 重点验证：
- `RollDice()` 在非 `exhausted` 状态下返回 `LastBattle`

`tower` 重点验证：
- `StartChallenge()` 返回 `Battle`

`arena` 重点验证：
- `RecordBattleResult(true/false)` 的 `Battle.Result` 与传入胜负一致

并让 `home/ranking` 的测试装配更新为新构造签名。

**Step 2: Run tests to verify they fail**

Run:
```bash
cd apps/backend && go test ./internal/modules/dungeon ./internal/modules/tower ./internal/modules/arena ./internal/modules/home ./internal/modules/ranking -v
```

Expected:
- FAIL because结果结构还没有统一战斗摘要
- FAIL because路由和测试装配还没注入 `pet` / `battle`

**Step 3: Write minimal implementation**

- 每个玩法服务新增最小依赖接口：
  - `GetBattleTeam(ctx, playerID)`
  - `Resolve(ctx, battle.Request)`
- `dungeon/domain.go`：
  - `DungeonRun` 新增 `LastBattle`
- `tower/domain.go`：
  - `TowerResult` 新增 `Battle`
- `arena/domain.go`：
  - `BattleResult` 新增 `Battle`
- `service.go` 规则：
  - `dungeon`：`exhausted` 不生成战斗；其他楼层生成系统怪快照并写 `LastBattle`
  - `tower`：每次挑战生成系统怪快照并写 `Battle`
  - `arena`：继续尊重当前传入 `won`，但改由 `battle.Resolve()` 产出统一摘要
- `router.go`：
  - 新建一个 `petService := pet.NewService(pet.NewMemoryRepository())`
  - 新建一个 `battleService := battle.NewService()`
  - 注入 `dungeonService / towerService / arenaService`

**Step 4: Run tests to verify they pass**

Run:
```bash
cd apps/backend && go test ./internal/modules/pet ./internal/modules/battle ./internal/modules/dungeon ./internal/modules/tower ./internal/modules/arena ./internal/modules/home ./internal/modules/ranking -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/backend/internal/modules/{pet,battle,dungeon,tower,arena,home,ranking} apps/backend/internal/bootstrap/router.go
git commit -m "feat: route gameplay through minimal pet and battle services"
```

### Task 4: 前端展示统一战斗摘要

**Files:**
- Modify: `apps/game-web/src/api/modules/dungeon.ts`
- Modify: `apps/game-web/src/api/modules/tower.ts`
- Modify: `apps/game-web/src/api/modules/arena.ts`
- Modify: `apps/game-web/src/pages/dungeons/DungeonRunPage.vue`
- Modify: `apps/game-web/src/pages/dungeons/__tests__/DungeonRunPage.spec.ts`
- Modify: `apps/game-web/src/pages/tower/PagodaPage.vue`
- Modify: `apps/game-web/src/pages/tower/SpiritTowerPage.vue`
- Modify: `apps/game-web/src/pages/tower/__tests__/PagodaPage.spec.ts`
- Modify: `apps/game-web/src/pages/tower/__tests__/SpiritTowerPage.spec.ts`
- Modify: `apps/game-web/src/pages/arena/ArenaPage.vue`
- Modify: `apps/game-web/src/pages/arena/__tests__/ArenaPage.spec.ts`

**Step 1: Write the failing tests**

增加断言：

```ts
expect(wrapper.text()).toContain("战斗摘要")
expect(wrapper.text()).toContain("战斗结果")
expect(wrapper.text()).toContain("我方战力")
expect(wrapper.text()).toContain("敌方战力")
```

mock 数据新增：

```ts
battle: {
  battle_type: "tower",
  result: "success",
  rounds: 1,
  attacker_power: 180,
  defender_power: 120,
}
```

`dungeon` 用 `last_battle`，`tower / arena` 用 `battle`。

**Step 2: Run tests to verify they fail**

Run:
```bash
pnpm --dir apps/game-web test -- DungeonRunPage ArenaPage PagodaPage SpiritTowerPage
```

Expected:
- FAIL because前端类型和页面还没有统一战斗摘要渲染

**Step 3: Write minimal implementation**

- 扩三类 API 类型
- 在三类页面增加统一“战斗摘要”卡片
- 文案只展示：
  - 战斗类型
  - 战斗结果
  - 回合数
  - 我方战力
  - 敌方战力

**Step 4: Run tests to verify they pass**

Run:
```bash
pnpm --dir apps/game-web test -- DungeonRunPage ArenaPage PagodaPage SpiritTowerPage
pnpm --dir apps/game-web build
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web
git commit -m "feat: show unified battle summaries in gameplay pages"
```

### Task 5: 全量验证与收口

**Files:**
- No code changes expected

**Step 1: Focused verification**

Run:
```bash
cd apps/backend && go test ./internal/modules/pet ./internal/modules/battle ./internal/modules/dungeon ./internal/modules/tower ./internal/modules/arena ./internal/modules/home ./internal/modules/ranking -v
pnpm --dir apps/game-web test -- DungeonRunPage ArenaPage PagodaPage SpiritTowerPage
```

**Step 2: Full verification**

Run:
```bash
make test-backend
make test-frontend
make smoke-backend
```

Expected: PASS

**Step 3: Review workspace**

Run:
```bash
git status --short
```

Expected:
- 本任务只新增 `pet`、`battle` 和相关玩法接入改动
- 不提交旧未跟踪文件 `docs/plans/2026-04-02-底座与首条闭环实施计划.md`
