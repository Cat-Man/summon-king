import { readSessionSnapshot } from "@/stores/session"

export interface APIResponse<T> {
  code: number
  message: string
  data: T
  trace_id: string
}

export class APIError extends Error {
  readonly code: number
  readonly status: number
  readonly traceID: string

  constructor(message: string, options: { code?: number; status: number; traceID?: string }) {
    super(message)
    this.name = "APIError"
    this.code = options.code ?? -1
    this.status = options.status
    this.traceID = options.traceID ?? ""
  }
}

const fallbackBaseURL = "http://localhost:8080/api/v1"

function joinURL(baseURL: string, path: string) {
  if (/^https?:\/\//.test(path)) {
    return path
  }
  return `${baseURL.replace(/\/+$/, "")}/${path.replace(/^\/+/, "")}`
}

export const apiBaseURL =
  (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim() || fallbackBaseURL

function buildHeaders(initHeaders: RequestInit["headers"]) {
  const headers = new Headers(initHeaders ?? {})
  if (!headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json")
  }

  const session = readSessionSnapshot()
  if (session.token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${session.token}`)
  }

  return headers
}

export async function apiRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(joinURL(apiBaseURL, path), {
    headers: buildHeaders(init.headers),
    ...init,
  })

  let payload: APIResponse<T> | null = null
  try {
    payload = (await response.json()) as APIResponse<T>
  } catch {
    if (!response.ok) {
      throw new APIError(`request failed with status ${response.status}`, {
        status: response.status,
      })
    }
    throw new APIError("response is not valid JSON", {
      status: response.status,
    })
  }

  if (!response.ok || payload.code !== 0) {
    throw new APIError(payload.message || "request failed", {
      code: payload.code,
      status: response.status,
      traceID: payload.trace_id,
    })
  }

  return payload.data
}
