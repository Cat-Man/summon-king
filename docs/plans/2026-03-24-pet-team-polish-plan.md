# Pet Team Polish Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 补齐阵容页的脏状态、满员提示和返回首页查看战力入口，让阵容编辑闭环更完整。

**Architecture:** 继续保持阵容页本地 `team pet ids` 为唯一编辑源，不新增全局事件总线。页面增加 `savedTeamPetIds` 基线用于判断脏状态，保存成功后更新基线并展示返回首页入口；首页仍依赖现有路由切换后重新挂载拉取最新数据。

**Tech Stack:** Vue 3、Vue Router、Pinia、Vitest、Vue Test Utils、TypeScript

### Task 1: 补失败测试

**Files:**
- Modify: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`

**Step 1: Write the failing test for dirty-state save**

```ts
test('disables save button when team has no unsaved changes', async () => {
  // initial load => disabled
  // after add action => enabled
  // after save => disabled again
})
```

**Step 2: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- PetTeamPage.spec.ts`
Expected: FAIL because save button is always enabled.

**Step 3: Write the failing test for full-team guard**

```ts
test('shows guard message when trying to add a bench pet to a full team', async () => {
  // mock 5 equipped pets + 1 bench pet
  // click bench add action
  // expect warning text and unchanged team
})
```

**Step 4: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- PetTeamPage.spec.ts`
Expected: FAIL because page currently accepts add attempt silently.

**Step 5: Write the failing test for return-home CTA**

```ts
test('offers a return-home action after save success', async () => {
  // edit team, save, click go-home CTA
  // expect router.push({ name: 'home' })
})
```

**Step 6: Run test to verify it fails**

Run: `pnpm --dir apps/game-web test -- PetTeamPage.spec.ts`
Expected: FAIL because page has no post-save home navigation CTA.

### Task 2: 实现脏状态、满员提示和首页返回入口

**Files:**
- Modify: `apps/game-web/src/pages/pets/PetTeamPage.vue`
- Modify: `apps/game-web/src/pages/pets/__tests__/PetTeamPage.spec.ts`

**Step 1: Add saved baseline and dirty-state computed**

```ts
const savedTeamPetIds = ref<Array<number | null>>(createEmptyTeamPetIds())
const isDirty = computed(() => JSON.stringify(savedTeamPetIds.value) !== JSON.stringify(editableTeamPetIds.value))
```

**Step 2: Guard bench add when team is full**

```ts
if (activeTeamCount >= TEAM_SLOT_COUNT) {
  loadError.value = '当前战斗队已满，请先下阵后再上阵'
  return
}
```

**Step 3: Update save button disabled state**

```ts
:disabled="pendingSave || !isDirty"
```

保存成功后同步：

```ts
savedTeamPetIds.value = [...editableTeamPetIds.value]
```

**Step 4: Add post-save return-home CTA**

```ts
function navigateHome() {
  router.push({ name: 'home' })
}
```

仅在保存成功提示存在时展示。

**Step 5: Run focused tests**

Run: `pnpm --dir apps/game-web test -- PetTeamPage.spec.ts`
Expected: PASS

### Task 3: 回归验证

**Files:**
- Modify: `docs/plans/2026-03-24-pet-team-polish-plan.md`

**Step 1: Run full frontend verification**

Run: `pnpm --dir apps/game-web test`
Expected: PASS

Run: `pnpm --dir apps/game-web lint`
Expected: PASS

Run: `pnpm --dir apps/game-web build`
Expected: PASS

Run: `git diff --check`
Expected: no output
