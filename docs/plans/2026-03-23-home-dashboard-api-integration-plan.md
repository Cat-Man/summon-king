# Home Dashboard API Integration Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 把 `apps/game-web` 首页从“前端 mock 模板 + 少量真实字段覆盖”的 hybrid 方案，收敛为“后端首页聚合 DTO + 前端纯映射渲染”的接口联调方案。

**Architecture:** 保留现有 `GET /api/v1/player/home/index` 路径，但把返回体从“仅玩家摘要”扩成“首页聚合 DTO”。后端负责输出首页工作台所需的主要结构：玩家摘要、每日必做、资源总览、消息流、功能矩阵、活动入口；前端 `home-dashboard` service 不再以 `createMockHomeDashboard()` 作为 API 模式基底，只保留 mock 模式专用数据和图片资源映射。

**Tech Stack:** Vue 3, TypeScript, Pinia, Vitest, Go, Gin.

**Prerequisites:** 先阅读并冻结以下文件：
- `docs/plans/2026-03-23-game-web-mock-api-integration-plan.md`
- `apps/game-web/src/pages/home/HomePage.vue`
- `apps/game-web/src/services/home-dashboard.ts`
- `apps/game-web/src/mocks/home-dashboard.ts`
- `apps/game-web/src/stores/session.ts`
- `apps/backend/internal/modules/player/service.go`
- `apps/backend/internal/modules/player/http.go`
- `apps/backend/internal/modules/player/service_test.go`

---

### Task 1: 扩展后端首页聚合 DTO

**Files:**
- Modify: `apps/backend/internal/modules/player/service.go`
- Modify: `apps/backend/internal/modules/player/service_test.go`

**Step 1: Write the failing test**

```go
func TestGetHomeIndex_ReturnsDashboardSections(t *testing.T) {
	ctx := context.Background()
	accountRepo := account.NewMemoryRepository()
	assetRepo := asset.NewMemoryRepository()

	accountService := account.NewService(accountRepo)
	playerEntity, err := accountService.CreateGuestPlayer(ctx, "web")
	if err != nil {
		t.Fatalf("expected player init success, got %v", err)
	}

	svc := NewService(NewRepository(accountRepo, assetRepo))
	home, err := svc.GetHomeIndex(ctx, playerEntity.PlayerID)
	if err != nil {
		t.Fatalf("expected home index success, got %v", err)
	}
	if len(home.DailyTodos) == 0 {
		t.Fatal("expected daily todos")
	}
	if len(home.Messages) == 0 {
		t.Fatal("expected messages")
	}
	if len(home.Entries) == 0 {
		t.Fatal("expected entries")
	}
	if home.ActivityEntry.Title == "" {
		t.Fatal("expected activity entry")
	}
}
```

**Step 2: Run test to verify it fails**

Run:
```bash
cd apps/backend && go test ./internal/modules/player -run TestGetHomeIndex_ReturnsDashboardSections -v
```

Expected: FAIL，因为 `HomeIndex` 还没有首页工作台字段。

**Step 3: Write minimal implementation**

```go
type HomeTodo struct {
	Title  string `json:"title"`
	Value  string `json:"value"`
	Action string `json:"action"`
}

type HomeActivityEntry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
}
```

并把 `HomeIndex` 扩展为：
- `daily_todos`
- `resources`
- `resource_icons`
- `messages`
- `entries`
- `activity_entry`

实现原则：
- `nickname/level/coin/diamond/last_login_at` 继续来自真实玩家数据
- `resources` 至少输出：等级、战力、活力、铜钱、元宝、声望
- `daily_todos/messages/entries/activity_entry` 由后端先输出首版工作台默认值，避免继续由前端 mock 模板托底

**Step 4: Run test to verify it passes**

Run:
```bash
cd apps/backend && go test ./internal/modules/player -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/backend/internal/modules/player
git commit -m "feat: enrich home index dashboard payload"
```

---

### Task 2: 前端改为消费后端首页聚合 DTO

**Files:**
- Modify: `apps/game-web/src/stores/session.ts`
- Modify: `apps/game-web/src/services/home-dashboard.ts`
- Modify: `apps/game-web/src/services/__tests__/home-dashboard.spec.ts`

**Step 1: Write the failing test**

```ts
test('loadHomeDashboard should map backend dashboard sections in api mode', async () => {
  const sessionStore = {
    ensureGuestSession: vi.fn().mockResolvedValue(undefined),
    fetchHomeIndex: vi.fn().mockResolvedValue({
      player_id: 2002,
      nickname: '联调游客2002',
      level: 7,
      coin: 1280,
      diamond: 96,
      last_login_at: '2026-03-23T10:00:00Z',
      daily_todos: [{ title: '签到状态', value: '今日未签', action: '前往签到' }],
      resources: [{ label: '等级', value: 'Lv.7' }],
      resource_icons: [{ label: '铜钱', value: '1280' }],
      messages: ['系统消息：VIP 每日宝箱可领取'],
      entries: ['世界地图', '联盟'],
      activity_entry: { title: '今日活动', description: '夺宝双倍', note: '21:00 开始' }
    })
  }

  const result = await loadHomeDashboard({ dataSource: 'api', sessionStore })

  expect(result.dailyTodos[0]?.title).toBe('签到状态')
  expect(result.messages[0]).toContain('VIP')
  expect(result.entries).toContain('联盟')
  expect(result.activityEntry.title).toBe('今日活动')
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/services/__tests__/home-dashboard.spec.ts
```

Expected: FAIL，因为前端 API 模式当前只覆盖 `hero/resources/resourceIcons`，不会读取后端工作台结构。

**Step 3: Write minimal implementation**

```ts
export interface HomeIndexResponse {
  player_id: number
  nickname: string
  level: number
  coin: number
  diamond: number
  last_login_at: string
  daily_todos: Array<{ title: string; value: string; action: string }>
  resources: Array<{ label: string; value: string }>
  resource_icons: Array<{ label: string; value: string }>
  messages: string[]
  entries: string[]
  activity_entry: { title: string; description: string; note: string }
}
```

并在 `createApiHomeDashboard()` 中改为：
- 直接消费后端 `daily_todos/resources/resource_icons/messages/entries/activity_entry`
- 只在前端补 `icon` 资源映射
- 不再以完整 `createMockHomeDashboard()` 作为 API 模式模板

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/services/__tests__/home-dashboard.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/stores/session.ts apps/game-web/src/services/home-dashboard.ts apps/game-web/src/services/__tests__/home-dashboard.spec.ts
git commit -m "feat: map home dashboard from backend payload"
```

---

### Task 3: 页面测试改成验证真实首页工作台渲染

**Files:**
- Modify: `apps/game-web/src/pages/home/__tests__/HomePage.spec.ts`
- Modify: `apps/game-web/src/pages/home/HomePage.vue`

**Step 1: Write the failing test**

```ts
test('renders backend driven dashboard sections', async () => {
  const loadedDashboard = createMockHomeDashboard()
  loadedDashboard.dailyTodos = [{ title: '签到状态', value: '今日已签', action: '查看奖励' }]
  loadedDashboard.messages = ['系统消息：今日活动已开启']
  loadedDashboard.entries = ['世界地图', '背包']

  vi.mocked(loadHomeDashboard).mockResolvedValueOnce(loadedDashboard)

  const wrapper = mountHomePage()
  await flushPromises()

  expect(wrapper.text()).toContain('今日已签')
  expect(wrapper.text()).toContain('系统消息：今日活动已开启')
  expect(wrapper.text()).toContain('背包')
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/home/__tests__/HomePage.spec.ts
```

Expected: FAIL，如果页面仍然依赖本地默认 mock 初值或没有正确渲染新结构。

**Step 3: Write minimal implementation**

实现要求：
- 页面初始值改成最小安全空态，避免 API 模式下闪现错误 mock 文案
- `loadError` 仍保留
- 页面只消费 `loadHomeDashboard()` 的结果，不直接引用 mock 工厂

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/home/__tests__/HomePage.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/home
git commit -m "refactor: render home page from loaded dashboard state"
```

---

### Task 4: 全量回归验证

**Files:**
- Modify: `docs/plans/2026-03-23-game-web-mock-api-integration-plan.md`

**Step 1: Run focused tests**

Run:
```bash
cd apps/backend && go test ./internal/modules/player -v
pnpm --dir apps/game-web exec vitest run src/services/__tests__/home-dashboard.spec.ts src/pages/home/__tests__/HomePage.spec.ts src/stores/__tests__/session.spec.ts
```

Expected: PASS

**Step 2: Run full verification**

Run:
```bash
pnpm --dir apps/game-web test
pnpm --dir apps/game-web lint
pnpm --dir apps/game-web build
cd apps/backend && go test ./...
git diff --check
```

Expected: PASS

**Step 3: Update execution log**

把本轮首页联调验证记录补进：
- `docs/plans/2026-03-23-game-web-mock-api-integration-plan.md`

**Step 4: Commit**

```bash
git add docs/plans/2026-03-23-game-web-mock-api-integration-plan.md
git commit -m "docs: record home dashboard api integration"
```
