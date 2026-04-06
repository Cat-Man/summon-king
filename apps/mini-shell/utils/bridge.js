const MINI_CHANNEL = "wxmini"

export function buildGameURL(baseURL, token) {
  const [baseWithQuery, hash = ""] = baseURL.split("#")
  const separator = baseWithQuery.includes("?") ? "&" : "?"
  const nextURL = `${baseWithQuery}${separator}channel=${MINI_CHANNEL}&token=${encodeURIComponent(token)}`

  if (!hash) {
    return nextURL
  }

  return `${nextURL}#${hash}`
}

export function buildWebviewRoute(gameURL) {
  return `/pages/webview/index?url=${encodeURIComponent(gameURL)}`
}

export function resolveLaunchToken(options = {}, fallbackToken = "") {
  if (typeof options.token === "string" && options.token.trim()) {
    return options.token
  }

  return fallbackToken
}
