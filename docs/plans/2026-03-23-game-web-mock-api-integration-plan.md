# Game Web Mock API Integration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为 `apps/game-web` 建立可切换的 mock / real 数据接入层，并先打通 `游客登录 -> session 建立 -> 首页聚合渲染` 的首条闭环，再为 `背包` 页联调铺好第二条链路。

**Architecture:** 前端继续保持 `Vue 3 + Pinia + fetch` 的轻量结构，不引入额外请求库。数据分为三层：`http` 只负责协议解包，`api services` 负责调用后端或 mock，页面只消费已经整理好的 view model。首条链路采用 `hybrid` 思路：会话与首页主数据优先走真实接口，当前后端尚未提供的“每日必做/消息流/功能矩阵”继续由前端 mock 补齐，避免因为后端聚合未完成而卡住首屏联调。

**Tech Stack:** Vue 3, TypeScript, Pinia, Vue Router, Vite, Vitest, fetch, Go Gin API.

**Prerequisites:** 先阅读并冻结以下文档与代码：
- `docs/plans/2026-03-21-召唤之王详细设计文档.md`
- `docs/plans/2026-03-21-召唤之王实施计划.md`
- `apps/game-web/src/api/http.ts`
- `apps/game-web/src/stores/session.ts`
- `apps/game-web/src/pages/home/HomePage.vue`
- `apps/game-web/src/pages/assets/InventoryPage.vue`
- `apps/backend/internal/bootstrap/router.go`
- `apps/backend/internal/modules/account/http.go`
- `apps/backend/internal/modules/player/http.go`
- `apps/backend/internal/modules/asset/http.go`

---

## 执行原则

- 所有新增功能必须先写失败测试，再写最小实现，再跑验证。
- 前端只新增最小必要目录：`src/api/services`、`src/mocks`。不提前引入 `axios`、`msw`、`vue-query`。
- 页面不直接调用 `fetch`，统一走 `request()` + `service`。
- 首页第一阶段只动态化后端已具备的数据：`nickname`、`level`、`coin`、`diamond`、`last_login_at`。
- 背包第二阶段才接 `wallet + inventory` 双接口；不要和首页首条闭环并行改，避免把 session 问题与 DTO 映射问题混在一起。

---

### Task 1: 冻结前端 mock / real 接入边界

**Files:**
- Modify: `apps/game-web/src/api/http.ts`
- Create: `apps/game-web/src/api/services/runtime.ts`
- Create: `apps/game-web/src/api/services/contracts.ts`
- Test: `apps/game-web/src/api/__tests__/http.spec.ts`

**Step 1: Write the failing test**

```ts
import { describe, expect, test, vi } from 'vitest'

import { request } from '../http'

describe('request', () => {
  test('prefixes api base url when input is relative', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({
          code: 0,
          message: 'ok',
          data: { ok: true },
          trace_id: 'trace-1'
        })
      })
    )

    await request('/player/home/index')

    expect(fetch).toHaveBeenCalledWith(
      'http://127.0.0.1:8080/player/home/index',
      expect.any(Object)
    )
  })
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/api/__tests__/http.spec.ts
```

Expected: FAIL because `request()` 还没有统一的 base url / mode 逻辑。

**Step 3: Write minimal implementation**

```ts
// apps/game-web/src/api/services/runtime.ts
export type APIMode = 'mock' | 'real' | 'hybrid'

export function getAPIMode(): APIMode {
  return (import.meta.env.VITE_API_MODE as APIMode | undefined) ?? 'hybrid'
}

export function getAPIBaseURL() {
  return import.meta.env.VITE_API_BASE_URL ?? '/api/v1'
}
```

```ts
// apps/game-web/src/api/services/contracts.ts
export interface LoginResponse {
  player_id: number
  token: string
  channel: string
}

export interface HomeIndexResponse {
  player_id: number
  nickname: string
  level: number
  coin: number
  diamond: number
  last_login_at: string
}
```

```ts
// apps/game-web/src/api/http.ts
import { getAPIBaseURL } from './services/runtime'

function resolveURL(input: string) {
  if (/^https?:\/\//.test(input)) {
    return input
  }
  return `${getAPIBaseURL()}${input}`
}
```

并保持 `request<T>()` 继续负责解包 `code/message/data/trace_id`。

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/api/__tests__/http.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/api
git commit -m "feat: add game web api runtime boundary"
```

---

### Task 2: 修复 session store 并打通游客登录

**Files:**
- Modify: `apps/game-web/src/stores/session.ts`
- Create: `apps/game-web/src/stores/__tests__/session.spec.ts`
- Create: `apps/game-web/src/api/services/session.ts`

**Step 1: Write the failing test**

```ts
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, expect, test, vi } from 'vitest'

import { useSessionStore } from '../session'

beforeEach(() => {
  setActivePinia(createPinia())
})

test('ensureGuestSession logs in and stores player id token', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          code: 0,
          message: 'ok',
          data: { player_id: 1001, token: 'guest-token-1001', channel: 'web' },
          trace_id: 'trace-login'
        })
      })
  )

  const store = useSessionStore()
  const session = await store.ensureGuestSession()

  expect(session.player_id).toBe(1001)
  expect(store.playerId).toBe(1001)
  expect(store.token).toBe('guest-token-1001')
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/stores/__tests__/session.spec.ts
```

Expected: FAIL because 现有 `session.ts` 错把 `request()` 的返回值再次当成 `APIResponse` 解包。

**Step 3: Write minimal implementation**

```ts
// apps/game-web/src/api/services/session.ts
import { request } from '@/api/http'
import type { LoginResponse } from './contracts'

export function loginGuest(channel = 'web') {
  return request<LoginResponse>('/player/auth/login', {
    method: 'POST',
    body: JSON.stringify({ channel })
  })
}
```

```ts
// apps/game-web/src/stores/session.ts
async ensureGuestSession(channel = 'web') {
  if (this.playerId && this.token) {
    return { player_id: this.playerId, token: this.token, channel }
  }
  const session = await loginGuest(channel)
  this.setSession(session.player_id, session.token)
  return session
}
```

保留 `fetchHomeIndex()`，但改成直接返回 `HomeIndexResponse`，不再访问不存在的 `response.data`。

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/stores/__tests__/session.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/stores apps/game-web/src/api/services
git commit -m "fix: repair guest session store flow"
```

---

### Task 3: 建立首页 service，支持 hybrid 聚合

**Files:**
- Create: `apps/game-web/src/api/services/home.ts`
- Create: `apps/game-web/src/mocks/home.ts`
- Test: `apps/game-web/src/api/services/__tests__/home.spec.ts`

**Step 1: Write the failing test**

```ts
import { expect, test, vi } from 'vitest'

import { fetchHomeDashboard } from '../home'

test('fetchHomeDashboard merges real home index with mock workstation blocks', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        code: 0,
        message: 'ok',
        data: {
          player_id: 1001,
          nickname: '游客1001',
          level: 1,
          coin: 5000,
          diamond: 0,
          last_login_at: '2026-03-23T10:00:00Z'
        },
        trace_id: 'trace-home'
      })
    })
  )

  const dashboard = await fetchHomeDashboard({ playerId: 1001, token: 'guest-token-1001' })

  expect(dashboard.nickname).toBe('游客1001')
  expect(dashboard.resourceIcons[0].value).toBe('5,000')
  expect(dashboard.dailyTodos).toHaveLength(3)
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/api/services/__tests__/home.spec.ts
```

Expected: FAIL because `fetchHomeDashboard()` 尚不存在。

**Step 3: Write minimal implementation**

```ts
// apps/game-web/src/mocks/home.ts
export const homeWorkbenchMock = {
  dailyTodos: [
    { title: '签到状态', value: '今日未签', action: '前往签到' },
    { title: '当前修行状态', value: '还有 18 分钟可领取', action: '查看修行' },
    { title: '当前推荐副本', value: '青木林地 · Boss 可挑战', action: '进入副本' }
  ],
  messages: [
    '世界消息：青木林地今日双倍经验已开启',
    '联盟消息：今晚 20:00 盟战锁定名单',
    '系统消息：VIP 每日宝箱可领取'
  ],
  entries: ['世界地图', '联盟', '幻兽', '背包', '竞技场', '庄园', '修行', '排行']
}
```

```ts
// apps/game-web/src/api/services/home.ts
import { request } from '@/api/http'
import { getAPIMode } from './runtime'
import { homeWorkbenchMock } from '@/mocks/home'

export async function fetchHomeDashboard(ctx: { playerId: number; token: string }) {
  const homeIndex =
    getAPIMode() === 'mock'
      ? {
          player_id: ctx.playerId,
          nickname: '本地游客',
          level: 1,
          coin: 5000,
          diamond: 0,
          last_login_at: new Date().toISOString()
        }
      : await request<HomeIndexResponse>(`/player/home/index?player_id=${ctx.playerId}`, {
          headers: ctx.token ? { Authorization: `Bearer ${ctx.token}` } : undefined
        })

  return {
    nickname: homeIndex.nickname,
    level: homeIndex.level,
    coin: homeIndex.coin,
    diamond: homeIndex.diamond,
    lastLoginAt: homeIndex.last_login_at,
    ...homeWorkbenchMock
  }
}
```

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/api/services/__tests__/home.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/api/services apps/game-web/src/mocks
git commit -m "feat: add home dashboard hybrid service"
```

---

### Task 4: 首页页面接入真实 session + 首页数据

**Files:**
- Modify: `apps/game-web/src/pages/home/HomePage.vue`
- Modify: `apps/game-web/src/pages/home/__tests__/HomePage.spec.ts`
- Optionally Modify: `apps/game-web/src/layouts/GameLayout.vue`

**Step 1: Write the failing test**

```ts
import { flushPromises, mount } from '@vue/test-utils'
import { createPinia } from 'pinia'
import { vi } from 'vitest'

import HomePage from '../HomePage.vue'

test('loads guest session and renders dynamic home stats', async () => {
  vi.stubGlobal(
    'fetch',
    vi.fn()
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          code: 0,
          message: 'ok',
          data: { player_id: 1001, token: 'guest-token-1001', channel: 'web' },
          trace_id: 'trace-login'
        })
      })
      .mockResolvedValueOnce({
        ok: true,
        json: async () => ({
          code: 0,
          message: 'ok',
          data: {
            player_id: 1001,
            nickname: '游客1001',
            level: 1,
            coin: 5000,
            diamond: 0,
            last_login_at: '2026-03-23T10:00:00Z'
          },
          trace_id: 'trace-home'
        })
      })
  )

  const wrapper = mount(HomePage, {
    global: {
      plugins: [createPinia()]
    }
  })

  await flushPromises()

  expect(wrapper.text()).toContain('游客1001')
  expect(wrapper.text()).toContain('Lv.1')
  expect(wrapper.text()).toContain('5,000')
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/home/__tests__/HomePage.spec.ts
```

Expected: FAIL because `HomePage.vue` 仍然只渲染静态常量。

**Step 3: Write minimal implementation**

```vue
<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useSessionStore } from '@/stores/session'
import { fetchHomeDashboard } from '@/api/services/home'

const sessionStore = useSessionStore()
const loading = ref(true)
const error = ref('')
const dashboard = ref<Awaited<ReturnType<typeof fetchHomeDashboard>> | null>(null)

onMounted(async () => {
  try {
    const session = await sessionStore.ensureGuestSession()
    dashboard.value = await fetchHomeDashboard({
      playerId: session.player_id,
      token: session.token
    })
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'home load failed'
  } finally {
    loading.value = false
  }
})

const resources = computed(() => {
  if (!dashboard.value) return []
  return [
    { label: '等级', value: `Lv.${dashboard.value.level}` },
    { label: '铜钱', value: dashboard.value.coinLabel },
    { label: '元宝', value: dashboard.value.diamondLabel }
  ]
})
</script>
```

模板要求：
- 保留 legacy 素材布局
- 新增 `loading` / `error` 占位
- 页面至少显示 `nickname`、`Lv.level`、`coin`、`diamond`
- 原有 `dailyTodos/messages/entries` 改为消费 `dashboard`

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/home/__tests__/HomePage.spec.ts
pnpm --dir apps/game-web exec vitest run src/stores/__tests__/session.spec.ts src/api/services/__tests__/home.spec.ts src/api/__tests__/http.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/home apps/game-web/src/layouts/GameLayout.vue
git commit -m "feat: integrate home page with guest session and api data"
```

---

### Task 5: 第二条链路预备，接背包 wallet + inventory

**Files:**
- Create: `apps/game-web/src/api/services/assets.ts`
- Create: `apps/game-web/src/mocks/assets.ts`
- Modify: `apps/game-web/src/pages/assets/InventoryPage.vue`
- Modify: `apps/game-web/src/pages/assets/__tests__/InventoryPage.spec.ts`

**Step 1: Write the failing test**

```ts
test('loads wallet and inventory from asset service', async () => {
  const wrapper = mount(InventoryPage, {
    global: { plugins: [createPinia()] }
  })

  await flushPromises()

  expect(wrapper.text()).toContain('普通背包')
  expect(wrapper.find('[data-testid="wallet-icon-铜钱"]').exists()).toBe(true)
})
```

**Step 2: Run test to verify it fails**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/assets/__tests__/InventoryPage.spec.ts
```

Expected: FAIL because页面尚未从 `wallet` / `inventory` 服务加载数据。

**Step 3: Write minimal implementation**

```ts
// apps/game-web/src/api/services/assets.ts
export async function fetchInventoryDashboard(ctx: { playerId: number; token: string }) {
  const [wallet, inventory] = await Promise.all([
    request<WalletResponse>(`/player/assets/wallet?player_id=${ctx.playerId}`, { headers }),
    request<InventoryItemResponse[]>(`/player/assets/inventory?player_id=${ctx.playerId}`, { headers })
  ])

  return adaptInventoryDashboard(wallet, inventory)
}
```

其中 `adaptInventoryDashboard()` 必须负责：
- 英文字段转前端展示字段
- `item_id / item_name` 到 legacy icon 的映射
- mock 模式下给出中文体验文案

**Step 4: Run test to verify it passes**

Run:
```bash
pnpm --dir apps/game-web exec vitest run src/pages/assets/__tests__/InventoryPage.spec.ts
```

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/assets apps/game-web/src/api/services/assets.ts apps/game-web/src/mocks/assets.ts
git commit -m "feat: integrate inventory page with asset services"
```

---

## 阶段验收命令

按顺序执行：

```bash
pnpm --dir apps/game-web exec vitest run src/api/__tests__/http.spec.ts
pnpm --dir apps/game-web exec vitest run src/stores/__tests__/session.spec.ts
pnpm --dir apps/game-web exec vitest run src/api/services/__tests__/home.spec.ts
pnpm --dir apps/game-web exec vitest run src/pages/home/__tests__/HomePage.spec.ts
pnpm --dir apps/game-web test
pnpm --dir apps/game-web lint
cd apps/backend && go test ./...
git diff --check
```

## 风险清单

- 后端 `home/index` 当前只返回基础信息，不包含“每日可领收益/消息流/主力幻兽/战力”；首页第一阶段仍是 hybrid 聚合。
- 后端 `asset` 默认库存仍是英文领域数据；新增 DTO 映射当前只覆盖 `potion_small`、`summon_scroll` 两种默认道具，后续扩品类时必须继续补映射表。
- `InventoryPage` 已补“预计售价展示”“二次确认弹层”“高价值道具默认锁定”前端安全层，但锁定状态目前只保存在页面会话中，刷新后会恢复默认锁定，后端尚未持久化玩家自定义锁定偏好。
- `/api/v1/player/assets/logs` 当前返回玩家全量内存流水，尚未补分页、时间筛选与高价值审计标签，前端只展示最近列表视图。

### Execution log (2026-03-23)

- `pnpm --dir apps/game-web exec vitest run src/api/__tests__/http.spec.ts src/stores/__tests__/session.spec.ts src/services/__tests__/home-dashboard.spec.ts src/pages/home/__tests__/HomePage.spec.ts` → PASS
- `pnpm --dir apps/game-web test` → PASS（22 files / 25 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- `go test ./cmd/api -run TestRunUsesConfiguredHTTPServer -v` → PASS
- `go test ./...` in `apps/backend` → PASS
- `curl -s -6 'http://[::1]:8080/healthz'` → 返回 `code=0`
- `curl -s -6 -X POST 'http://[::1]:8080/api/v1/player/auth/login' ...` → 返回 `player_id=1`
- `curl -s -6 'http://[::1]:8080/api/v1/player/home/index?player_id=1' -H 'Authorization: Bearer guest-token-1'` → 返回 `nickname=游客1, level=1, coin=0, diamond=0`
- 浏览器打开 `http://127.0.0.1:4173/`（默认 `mock` 模式）→ 首页正常渲染游客昵称、Lv.1、资源卡片与功能矩阵
- `pnpm --dir apps/game-web exec vitest run src/config/__tests__/vite-dev-proxy.spec.ts src/services/__tests__/inventory-dashboard.spec.ts src/pages/assets/__tests__/InventoryPage.spec.ts` → PASS
- `pnpm --dir apps/game-web test` → PASS（24 files / 30 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- 浏览器打开 `http://127.0.0.1:4174/assets`（`VITE_GAME_DATA_SOURCE=api`）→ 通过 Vite proxy 成功触发 `/api/v1/player/auth/login`、`/api/v1/player/assets/wallet`、`/api/v1/player/assets/inventory`，页面显示 `2 / 30`、`中级经验丹 x5`、`召唤卷轴 x1`
- 浏览器在同一 `api` 模式页面点击 `中级经验丹 -> 使用1个` → 成功提示 `已使用 1 个中级经验丹`，库存刷新为 `x4`，Network 记录 `POST /api/v1/player/assets/inventory/use [200]`，随后 `wallet` 与 `inventory` 回源刷新均为 `[200]`
- 浏览器在同一 `api` 模式页面点击 `召唤卷轴 -> 出售1个` → 成功提示 `已出售 1 个召唤卷轴`，铜钱更新为 `50`，`召唤卷轴` 从列表移除，Network 记录 `POST /api/v1/player/assets/inventory/sell [200]`，随后 `wallet` 与 `inventory` 回源刷新均为 `[200]`
- `pnpm --dir apps/game-web test` → PASS（24 files / 34 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- `go test ./...` in `apps/backend` → PASS
- `git diff --check` → PASS
- `pnpm --dir apps/game-web exec vitest run src/pages/assets/__tests__/InventoryPage.spec.ts` → PASS（新增高价值锁定、出售确认、取消出售覆盖）
- `pnpm --dir apps/game-web test` → PASS（24 files / 36 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- `go test ./internal/modules/asset -run TestHandler_LogsReturnsResourceChangeEntries -v` → 先 FAIL（`404`），补 `/api/v1/player/assets/logs` 后 PASS
- `go test ./internal/modules/asset` → PASS
- `pnpm --dir apps/game-web exec vitest run src/services/__tests__/inventory-dashboard.spec.ts src/pages/assets/__tests__/InventoryPage.spec.ts` → 先 FAIL（缺 `loadInventoryLogs` 与页面入口），补服务与页面后 PASS（12 tests）
- `pnpm --dir apps/game-web test` → PASS（24 files / 38 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- `go test ./...` in `apps/backend` → PASS
- `git diff --check` → PASS
- 真实浏览器 `api` 模式 smoke 本轮未执行：本机 `8080` 已被外部服务占用，当前后端配置固定监听 `8080`
- `pnpm --dir apps/game-web exec vitest run src/services/__tests__/inventory-dashboard.spec.ts src/pages/assets/__tests__/InventoryPage.spec.ts` → 先 FAIL（缺 `loadRecentAssetLogs` 导出与流水字段展示），补最小实现后 PASS（13 tests）
- `pnpm --dir apps/game-web test` → PASS（24 files / 39 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- `go test ./internal/bootstrap -run TestLoadConfig_UsesHTTPPortEnv -v` → 先 FAIL（固定 `8080`），补 `HTTP_PORT` 环境变量读取后 PASS
- `go test ./...` in `apps/backend` → PASS
- 临时启动 `HTTP_PORT=18080 go run ./cmd/api` 与 `VITE_GAME_DATA_SOURCE=api VITE_DEV_API_PROXY_TARGET=http://127.0.0.1:18080 pnpm --dir apps/game-web dev --host 127.0.0.1 --port 4176`
- 浏览器打开 `http://127.0.0.1:4176/assets` → 登录、钱包、背包读取成功；页面显示 `2 / 30`、`中级经验丹 x5`、`召唤卷轴 x1`
- 浏览器点击 `查看最近流水` → 成功触发 `GET /api/v1/player/assets/logs?player_id=2&limit=5 [200]`，初始游客无流水，面板仅展示标题
- 浏览器点击 `中级经验丹 -> 使用1个` → 成功提示 `已使用 1 个中级经验丹`，库存刷新为 `x4`，流水面板自动新增 `使用中级经验丹 / 使用道具 / 原因：inventory_use`
- 浏览器点击 `召唤卷轴 -> 解除锁定 -> 出售1个 -> 确认出售` → 成功提示 `已出售 1 个召唤卷轴`，铜钱更新为 `50`，背包容量变为 `1 / 30`，流水面板自动新增 `出售召唤卷轴 / 出售道具 / 原因：inventory_sell / +50 铜钱`
- `pnpm --dir apps/game-web exec vitest run src/services/__tests__/inventory-dashboard.spec.ts src/pages/assets/__tests__/InventoryPage.spec.ts` → 先 FAIL（新增流水分页/筛选契约测试，服务与页面均未实现）
- 同一命令二次执行 → PASS（16 tests），`loadRecentAssetLogs` 已改为 `{ page, pageSize, changeType } -> { items, total, page, pageSize }`，页面补齐筛选与上一页/下一页
- `pnpm --dir apps/game-web test` → PASS（24 files / 42 tests）
- `pnpm --dir apps/game-web lint` → PASS
- `pnpm --dir apps/game-web build` → PASS
- 真实浏览器 `api` 模式二次 smoke 发现契约偏差：前端将 `changeType=all` 原样发送为 `change_type=all`，后端按精确类型过滤，导致使用/出售成功后流水查询仍返回空列表
- `pnpm --dir apps/game-web exec vitest run src/services/__tests__/inventory-dashboard.spec.ts` → 先 FAIL（新增回归测试，要求 API 模式下 `all` 转为空过滤值），补最小实现后 PASS（8 tests）
- 再次启动 `HTTP_PORT=18080 go run ./cmd/api` 与 `VITE_GAME_DATA_SOURCE=api VITE_DEV_API_PROXY_TARGET=http://127.0.0.1:18080 pnpm --dir apps/game-web dev --host 127.0.0.1 --port 4176`
- 浏览器打开 `http://127.0.0.1:4176/assets` → 成功触发 `GET /api/v1/player/assets/logs?player_id=4&page=1&page_size=2&change_type= [200]`，初始游客无流水
- 浏览器点击 `中级经验丹 -> 使用1个` → 成功提示 `已使用 1 个中级经验丹`，库存刷新为 `x4`，流水面板展示 `使用中级经验丹 / 使用道具 / 原因：inventory_use`
- 浏览器点击 `召唤卷轴 -> 解除锁定 -> 出售1个 -> 确认出售` → 成功提示 `已出售 1 个召唤卷轴`，铜钱更新为 `50`，背包容量变为 `1 / 30`，流水面板展示 `出售召唤卷轴 / 出售道具 / +50 铜钱`
- 浏览器点击 `出售道具` 筛选 → 成功触发 `GET /api/v1/player/assets/logs?player_id=4&page=1&page_size=2&change_type=inventory_sell [200]`，仅展示出售流水 1 条
