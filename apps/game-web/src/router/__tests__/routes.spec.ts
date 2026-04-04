import { router } from "../index"

test("contains dungeon route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/dungeon")

  expect(hasRoute).toBe(true)
})

test("contains arena route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/arena")

  expect(hasRoute).toBe(true)
})
