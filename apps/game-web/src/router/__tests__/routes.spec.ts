import { router } from "../index"

test("contains dungeon route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/dungeon")

  expect(hasRoute).toBe(true)
})

test("contains arena route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/arena")

  expect(hasRoute).toBe(true)
})

test("contains pet route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/pet")

  expect(hasRoute).toBe(true)
})

test("contains alliance route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/alliance")

  expect(hasRoute).toBe(true)
})

test("contains alliance war route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/alliance-war")

  expect(hasRoute).toBe(true)
})
