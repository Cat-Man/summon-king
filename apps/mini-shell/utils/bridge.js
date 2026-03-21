function buildGameURL(baseURL, token) {
  const url = new URL(baseURL)
  url.searchParams.set('channel', 'wxmini')
  url.searchParams.set('token', token)
  return url.toString()
}

module.exports = {
  buildGameURL
}
