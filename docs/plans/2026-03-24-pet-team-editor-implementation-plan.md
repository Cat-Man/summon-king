# Pet Team Editor Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 为阵容页补齐可编辑上阵、下阵、换位交互，并保证保存后再次进入首页时战力按最新阵容展示。

**Architecture:** 本次仅改 Web 前端页面层与必要的前端数据映射，不改后端接口契约。阵容页维护本地可编辑 `team pet ids`，由它派生当前战斗队、幻兽栏状态、概览统计和保存载荷；首页继续依赖现有 `home/index` 接口按最新队伍重新拉取战力。

**Tech Stack:** Vue 3、Pinia、Vue Test Utils、Vitest、TypeScript

### Task 1: 补齐阵容编辑失败测试

**Files:**
- Modify: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`
- Test: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`

**Step 1: Write the failing test**

补 3 组行为测试：
- 替补幻兽点击上阵后进入当前战斗队，且幻兽栏状态切为“已上阵”
- 已上阵幻兽点击下阵后移出当前战斗队，且幻兽栏状态切为“可替补”
- 当前战斗队点击换位按钮后，保存接口按最新顺序提交

示例断言：

```ts
await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')
expect(wrapper.text()).toContain('3 / 5')
expect(wrapper.get('[data-testid="team-member-3"]').text()).toContain('雷角牛')

await wrapper.get('[data-testid="team-remove-2"]').trigger('click')
expect(wrapper.text()).toContain('1 / 5')
expect(wrapper.text()).toContain('雷角牛')

await wrapper.get('[data-testid="team-move-down-2"]').trigger('click')
await wrapper.get('[data-testid="save-team"]').trigger('click')
expect(savePetTeamSelection).toHaveBeenCalledWith(
  expect.objectContaining({ petIds: [1, 2] })
)
```

**Step 2: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts`

Expected: FAIL，报缺少编辑按钮或保存参数顺序不匹配。

**Step 3: Write minimal implementation**

先不要改服务层，只让测试准确描述目标交互与保存顺序。

**Step 4: Run test to verify it passes**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts`

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts
git commit -m "test: cover pet team editor interactions"
```

### Task 2: 实现阵容页本地编辑状态

**Files:**
- Modify: `apps/game-web/src/pages/pets/PetTeamPage.vue`
- Modify: `apps/game-web/src/services/pet-dashboard.types.ts`
- Modify: `apps/game-web/src/services/pet-dashboard.ts`
- Test: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`
- Test: `apps/game-web/src/services/__tests__/pet-dashboard.spec.ts`

**Step 1: Write the failing test**

在 `apps/game-web/src/services/__tests__/pet-dashboard.spec.ts` 增加 1 条映射测试，要求 API 模式下阵容页 `roster` 数据包含可用于回填当前战斗队的最小信息，例如 `role`。

示例断言：

```ts
expect(dashboard.roster[0]).toMatchObject({
  petId: 3,
  name: '雷角牛',
  role: '前排承伤'
})
```

**Step 2: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- src/services/__tests__/pet-dashboard.spec.ts`

Expected: FAIL，报 `role` 缺失。

**Step 3: Write minimal implementation**

实现最小闭环：
- `PetTeamRosterItem` 增加 `role`
- `loadPetTeamDashboard()` 构造 `roster` 时补齐 `role`
- `PetTeamPage.vue` 增加本地 `editableTeamIds`
- 基于 `editableTeamIds` 计算当前战斗队、幻兽栏状态、概览统计、主养成目标
- 增加 `上阵`、`下阵`、`上移`、`下移` 按钮和对应 `data-testid`
- 队伍满 5 人时阻止继续上阵，并给出页面提示
- 保存成功后清空脏状态提示，保留最新编辑顺序

**Step 4: Run test to verify it passes**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts src/services/__tests__/pet-dashboard.spec.ts`

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/pets/PetTeamPage.vue apps/game-web/src/services/pet-dashboard.ts apps/game-web/src/services/pet-dashboard.types.ts apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts apps/game-web/src/services/__tests__/pet-dashboard.spec.ts
git commit -m "feat: add editable pet team interactions"
```

### Task 3: 验证首页联动边界

**Files:**
- Modify: `apps/game-web/src/pages/pets/PetTeamPage.vue`
- Optional: `apps/game-web/src/stores/session.ts`
- Test: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`

**Step 1: Write the failing test**

补 1 条页面行为测试，要求保存成功后仍显示最新综合战力与最新队伍顺序，不回退到旧快照。

示例断言：

```ts
await wrapper.get('[data-testid="roster-action-雷角牛"]').trigger('click')
await wrapper.get('[data-testid="save-team"]').trigger('click')
expect(wrapper.text()).toContain('450')
expect(wrapper.text()).toContain('阵容已保存')
```

**Step 2: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts`

Expected: FAIL，若页面仍使用旧 `dashboard.team` 快照则战力不会更新。

**Step 3: Write minimal implementation**

只保持当前前端联动边界：
- 阵容页保存后以本地编辑状态作为当前展示真值
- 不新增首页缓存失效或事件总线
- 依赖首页路由重新挂载时调用现有 `loadHomeDashboard()` 拉取最新战力

**Step 4: Run test to verify it passes**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts`

Expected: PASS

**Step 5: Commit**

```bash
git add apps/game-web/src/pages/pets/PetTeamPage.vue apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts
git commit -m "test: verify pet team save keeps latest local state"
```

### Task 4: 全量验证

**Files:**
- Modify: `docs/plans/2026-03-24-pet-team-editor-implementation-plan.md`

**Step 1: Run targeted tests**

Run: `pnpm --dir apps/game-web test -- src/pages/pets/__tests__/PetTeamPage.spec.ts src/services/__tests__/pet-dashboard.spec.ts`

Expected: PASS

**Step 2: Run project verification**

Run: `pnpm --dir apps/game-web test`
Expected: PASS

Run: `pnpm --dir apps/game-web lint`
Expected: PASS

Run: `pnpm --dir apps/game-web build`
Expected: PASS

Run: `git diff --check`
Expected: no output

**Step 3: Commit**

```bash
git add docs/plans/2026-03-24-pet-team-editor-implementation-plan.md
git commit -m "docs: add pet team editor implementation plan"
```
