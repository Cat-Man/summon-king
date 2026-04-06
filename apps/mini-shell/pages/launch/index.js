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
} from "../../utils/bridge.js"

const DEFAULT_GAME_URL = "https://game.xxx.com"
const DEFAULT_FALLBACK_TOKEN = "guest-token"
const DEFAULT_API_BASE_URL = "http://localhost:8080/api/v1"

function redirectToGame(gameBaseURL, token) {
  const gameURL = buildGameURL(gameBaseURL, token)
  wx.redirectTo({
    url: buildWebviewRoute(gameURL),
  })
}

function loginWithWX() {
  return new Promise((resolve, reject) => {
    wx.login({
      success: (res) => {
        if (typeof res.code === "string" && res.code.trim()) {
          resolve(res.code)
          return
        }
        reject(new Error("wx.login code is empty"))
      },
      fail: (error) => reject(error),
    })
  })
}

function exchangeWXCode(apiBaseURL, code) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: buildWXMiniLoginExchangeURL(apiBaseURL),
      method: "POST",
      header: {
        "Content-Type": "application/json",
      },
      data: buildWXMiniLoginExchangePayload(code),
      success: (res) => {
        if (res.statusCode < 200 || res.statusCode >= 300) {
          reject(new Error(`bridge exchange failed: ${res.statusCode}`))
          return
        }
        const token = extractUnifiedToken(res.data)
        if (!token) {
          reject(new Error("bridge exchange token is empty"))
          return
        }
        resolve(token)
      },
      fail: (error) => reject(error),
    })
  })
}

function bootstrapWXMiniSession(apiBaseURL, unifiedToken, gameBaseURL) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: buildWXMiniSessionBootstrapURL(apiBaseURL),
      method: "POST",
      header: {
        "Content-Type": "application/json",
      },
      data: buildWXMiniSessionBootstrapPayload(unifiedToken, gameBaseURL),
      success: (res) => {
        if (res.statusCode < 200 || res.statusCode >= 300) {
          reject(new Error(`bridge bootstrap failed: ${res.statusCode}`))
          return
        }
        const gameURL = extractWXMiniBootstrapGameURL(res.data)
        if (!gameURL) {
          reject(new Error("bridge bootstrap game_url is empty"))
          return
        }
        resolve(gameURL)
      },
      fail: (error) => reject(error),
    })
  })
}

Page({
  data: {
    statusText: "正在进入游戏...",
  },

  async onLoad(options = {}) {
    const app = getApp()
    const gameBaseURL = options.gameURL || app?.globalData?.gameBaseURL || DEFAULT_GAME_URL
    const fallbackToken = app?.globalData?.fallbackToken || DEFAULT_FALLBACK_TOKEN
    const apiBaseURL = resolveBridgeAPIBaseURL(
      options,
      app?.globalData?.apiBaseURL || DEFAULT_API_BASE_URL,
    )
    const directToken = resolveLaunchToken(options, "")

    if (directToken) {
      redirectToGame(gameBaseURL, directToken)
      return
    }

    try {
      this.setData({ statusText: "正在拉起微信登录..." })
      const code = await loginWithWX()
      this.setData({ statusText: "正在验证登录..." })
      const unifiedToken = await exchangeWXCode(apiBaseURL, code)
      this.setData({ statusText: "正在生成启动参数..." })
      const bootstrapGameURL = await bootstrapWXMiniSession(apiBaseURL, unifiedToken, gameBaseURL)
      wx.redirectTo({
        url: buildWebviewRoute(bootstrapGameURL),
      })
    } catch (_error) {
      this.setData({ statusText: "桥接失败，正在使用访客身份进入..." })
      redirectToGame(gameBaseURL, fallbackToken)
    }
  },
})
