const MINI_CHANNEL = "wxmini"
const DEFAULT_BRIDGE_API_BASE_URL = "http://localhost:8080/api/v1"

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

export function resolveBridgeAPIBaseURL(options = {}, fallbackBaseURL = DEFAULT_BRIDGE_API_BASE_URL) {
  if (typeof options.apiBaseURL === "string" && options.apiBaseURL.trim()) {
    return options.apiBaseURL
  }

  return fallbackBaseURL
}

export function buildWXMiniLoginExchangeURL(apiBaseURL) {
  return `${String(apiBaseURL || DEFAULT_BRIDGE_API_BASE_URL).replace(/\/+$/, "")}/bridge/wxmini/login/exchange`
}

export function buildWXMiniLoginExchangePayload(code) {
  return {
    code,
  }
}

export function buildWXMiniSessionBootstrapURL(apiBaseURL) {
  return `${String(apiBaseURL || DEFAULT_BRIDGE_API_BASE_URL).replace(/\/+$/, "")}/bridge/wxmini/session/bootstrap`
}

export function buildWXMiniSessionBootstrapPayload(unifiedToken, gameBaseURL) {
  return {
    unified_token: unifiedToken,
    game_base_url: gameBaseURL,
  }
}

export function extractUnifiedToken(responseBody) {
  if (!responseBody || typeof responseBody !== "object") {
    return ""
  }

  const payload = responseBody.code === 0 && responseBody.data && typeof responseBody.data === "object"
    ? responseBody.data
    : responseBody

  if (typeof payload.unified_token === "string" && payload.unified_token.trim()) {
    return payload.unified_token
  }
  if (typeof payload.token === "string" && payload.token.trim()) {
    return payload.token
  }

  return ""
}

export function extractWXMiniBootstrapGameURL(responseBody) {
  if (!responseBody || typeof responseBody !== "object") {
    return ""
  }

  const payload = responseBody.code === 0 && responseBody.data && typeof responseBody.data === "object"
    ? responseBody.data
    : responseBody

  if (typeof payload.game_url === "string" && payload.game_url.trim()) {
    return payload.game_url
  }

  return ""
}
