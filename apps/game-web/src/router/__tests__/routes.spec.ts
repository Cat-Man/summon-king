import { router } from "../index"

test("contains dungeon route", () => {
  const hasRoute = router.getRoutes().some((route) => route.path === "/dungeon")

  expect(hasRoute).toBe(true)
})
