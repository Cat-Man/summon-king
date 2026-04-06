import { describe, expect, test } from "vitest"

import {
  buildGameURL,
  buildWebviewRoute,
  buildWXMiniLoginExchangePayload,
  buildWXMiniLoginExchangeURL,
  buildWXMiniSessionBootstrapPayload,
  buildWXMiniSessionBootstrapURL,
  extractUnifiedToken,
  extractWXMiniBootstrapGameURL,
  resolveBridgeAPIBaseURL,
  resolveLaunchToken,
} from "../utils/bridge.js"

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

  test("resolveBridgeAPIBaseURL prefers launch option", () => {
    expect(resolveBridgeAPIBaseURL({ apiBaseURL: "https://api.example.com/v1" }, "http://localhost:8080/api/v1")).toBe(
      "https://api.example.com/v1",
    )
  })

  test("buildWXMiniLoginExchangeURL appends bridge exchange path", () => {
    expect(buildWXMiniLoginExchangeURL("http://localhost:8080/api/v1/")).toBe(
      "http://localhost:8080/api/v1/bridge/wxmini/login/exchange",
    )
  })

  test("buildWXMiniLoginExchangePayload serializes wx.login code", () => {
    expect(buildWXMiniLoginExchangePayload("wx-code-1")).toEqual({
      code: "wx-code-1",
    })
  })

  test("buildWXMiniSessionBootstrapURL appends bridge bootstrap path", () => {
    expect(buildWXMiniSessionBootstrapURL("http://localhost:8080/api/v1/")).toBe(
      "http://localhost:8080/api/v1/bridge/wxmini/session/bootstrap",
    )
  })

  test("buildWXMiniSessionBootstrapPayload serializes unified token and game base url", () => {
    expect(buildWXMiniSessionBootstrapPayload("wx-token-1", "https://game.xxx.com")).toEqual({
      unified_token: "wx-token-1",
      game_base_url: "https://game.xxx.com",
    })
  })

  test("extractUnifiedToken supports standard API envelope", () => {
    expect(
      extractUnifiedToken({
        code: 0,
        data: {
          unified_token: "u-token",
        },
      }),
    ).toBe("u-token")
  })

  test("extractUnifiedToken supports direct payload token fallback", () => {
    expect(
      extractUnifiedToken({
        token: "guest-token",
      }),
    ).toBe("guest-token")
  })

  test("extractWXMiniBootstrapGameURL supports standard API envelope", () => {
    expect(
      extractWXMiniBootstrapGameURL({
        code: 0,
        data: {
          game_url: "https://game.xxx.com?channel=wxmini&token=wx-token-1",
        },
      }),
    ).toBe("https://game.xxx.com?channel=wxmini&token=wx-token-1")
  })
})
