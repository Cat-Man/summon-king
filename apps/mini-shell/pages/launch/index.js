import { buildGameURL, buildWebviewRoute, resolveLaunchToken } from "../../utils/bridge.js"

const DEFAULT_GAME_URL = "https://game.xxx.com"
const DEFAULT_FALLBACK_TOKEN = "guest-token"

Page({
  data: {
    statusText: "正在进入游戏...",
  },

  onLoad(options = {}) {
    const app = getApp()
    const gameBaseURL = options.gameURL || app?.globalData?.gameBaseURL || DEFAULT_GAME_URL
    const fallbackToken = app?.globalData?.fallbackToken || DEFAULT_FALLBACK_TOKEN
    const token = resolveLaunchToken(options, fallbackToken)
    const gameURL = buildGameURL(gameBaseURL, token)

    wx.redirectTo({
      url: buildWebviewRoute(gameURL),
    })
  },
})
