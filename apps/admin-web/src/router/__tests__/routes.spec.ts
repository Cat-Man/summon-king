import { router } from "../index"

test("contains admin section routes", () => {
  const routeNames = router.getRoutes().map((route) => route.name)

  expect(routeNames).toContain("dashboard")
  expect(routeNames).toContain("config")
  expect(routeNames).toContain("gm")
})

test("redirects root to dashboard", async () => {
  await router.push("/")
  await router.isReady()

  expect(router.currentRoute.value.name).toBe("dashboard")
})
