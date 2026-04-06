import { describe, expect, test } from "vitest"

import { buildGameURL, buildWebviewRoute, resolveLaunchToken } from "../utils/bridge.js"

describe("bridge helpers", () => {
  test("buildGameURL adds channel and token", () => {
    expect(buildGameURL("https://game.xxx.com", "abc")).toBe("https://game.xxx.com?channel=wxmini&token=abc")
  })

  test("buildGameURL preserves existing query params", () => {
    expect(buildGameURL("https://game.xxx.com/play?foo=1", "abc")).toBe(
      "https://game.xxx.com/play?foo=1&channel=wxmini&token=abc",
    )
  })

  test("buildWebviewRoute encodes the target url", () => {
    expect(buildWebviewRoute("https://game.xxx.com?channel=wxmini&token=abc")).toBe(
      "/pages/webview/index?url=https%3A%2F%2Fgame.xxx.com%3Fchannel%3Dwxmini%26token%3Dabc",
    )
  })

  test("resolveLaunchToken prefers the incoming token", () => {
    expect(resolveLaunchToken({ token: "abc" }, "guest-token")).toBe("abc")
  })

  test("resolveLaunchToken falls back when token is missing", () => {
    expect(resolveLaunchToken({}, "guest-token")).toBe("guest-token")
  })
})
