function buildGameURL(baseURL, token) {
  const url = new URL(baseURL)
  url.searchParams.set('channel', 'wxmini')
  url.searchParams.set('token', token)
  return url.toString()
}

function buildMiniLoginPayload(code) {
  return {
    code,
    channel: 'wxmini'
  }
}

function buildMiniPayPayload(orderNo) {
  return {
    order_no: orderNo,
    channel: 'wxmini'
  }
}

module.exports = {
  buildGameURL,
  buildMiniLoginPayload,
  buildMiniPayPayload
}
