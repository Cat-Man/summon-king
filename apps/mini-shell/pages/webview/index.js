const { buildGameURL } = require('../../utils/bridge')

Page({
  data: {
    gameURL: ''
  },
  onLoad(options) {
    const app = getApp()
    const token = options.token || 'guest-token'
    const baseURL = app.globalData.gameBaseURL || 'https://game.xxx.com'
    this.setData({
      gameURL: buildGameURL(baseURL, token)
    })
  }
})
